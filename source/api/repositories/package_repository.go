package repositories

import (
	"context"
	"fmt"
	"time"

	"code.cloudfoundry.org/korifi/api/authorization"
	apierrors "code.cloudfoundry.org/korifi/api/errors"
	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/controllers/controllers/workloads/packages"
	"code.cloudfoundry.org/korifi/tools/dockercfg"
	"code.cloudfoundry.org/korifi/tools/k8s"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	kind = "CFPackage"

	PackageStateAwaitingUpload = "AWAITING_UPLOAD"
	PackageStateReady          = "READY"

	PackageResourceType = "Package"
)

var packageTypeToLifecycleType = map[korifiv1alpha1.PackageType]korifiv1alpha1.LifecycleType{
	"bits":   "buildpack",
	"docker": "docker",
}

type PackageRepo struct {
	userClientFactory    authorization.UserK8sClientFactory
	namespaceRetriever   NamespaceRetriever
	namespacePermissions *authorization.NamespacePermissions
	repositoryCreator    RepositoryCreator
	repositoryPrefix     string
	awaiter              Awaiter[*korifiv1alpha1.CFPackage]
}

func NewPackageRepo(
	userClientFactory authorization.UserK8sClientFactory,
	namespaceRetriever NamespaceRetriever,
	authPerms *authorization.NamespacePermissions,
	repositoryCreator RepositoryCreator,
	repositoryPrefix string,
	awaiter Awaiter[*korifiv1alpha1.CFPackage],
) *PackageRepo {
	return &PackageRepo{
		userClientFactory:    userClientFactory,
		namespaceRetriever:   namespaceRetriever,
		namespacePermissions: authPerms,
		repositoryCreator:    repositoryCreator,
		repositoryPrefix:     repositoryPrefix,
		awaiter:              awaiter,
	}
}

type PackageRecord struct {
	GUID        string
	UID         types.UID
	Type        string
	AppGUID     string
	SpaceGUID   string
	State       string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	Labels      map[string]string
	Annotations map[string]string
	ImageRef    string
}

type ListPackagesMessage struct {
	GUIDs    []string
	AppGUIDs []string
	States   []string
}

type CreatePackageMessage struct {
	Type      string
	AppGUID   string
	SpaceGUID string
	Metadata  Metadata
	Data      *PackageData
}

type PackageData struct {
	Image    string
	Username *string
	Password *string
}

func (message CreatePackageMessage) toCFPackage() *korifiv1alpha1.CFPackage {
	packageGUID := uuid.NewString()
	pkg := &korifiv1alpha1.CFPackage{
		TypeMeta: metav1.TypeMeta{
			Kind:       kind,
			APIVersion: korifiv1alpha1.GroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        packageGUID,
			Namespace:   message.SpaceGUID,
			Labels:      message.Metadata.Labels,
			Annotations: message.Metadata.Annotations,
		},
		Spec: korifiv1alpha1.CFPackageSpec{
			Type: korifiv1alpha1.PackageType(message.Type),
			AppRef: corev1.LocalObjectReference{
				Name: message.AppGUID,
			},
		},
	}

	if message.Type == "docker" {
		pkg.Spec.Source.Registry.Image = message.Data.Image
	}

	return pkg
}

type UpdatePackageMessage struct {
	GUID          string
	MetadataPatch MetadataPatch
}

type UpdatePackageSourceMessage struct {
	GUID                string
	SpaceGUID           string
	ImageRef            string
	RegistrySecretNames []string
}

func (r *PackageRepo) CreatePackage(ctx context.Context, authInfo authorization.Info, message CreatePackageMessage) (PackageRecord, error) {
	userClient, err := r.userClientFactory.BuildClient(authInfo)
	if err != nil {
		return PackageRecord{}, fmt.Errorf("failed to build user client: %w", err)
	}

	cfApp := &korifiv1alpha1.CFApp{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: message.SpaceGUID,
			Name:      message.AppGUID,
		},
	}

	err = userClient.Get(ctx, client.ObjectKeyFromObject(cfApp), cfApp)
	if err != nil {
		return PackageRecord{},
			apierrors.AsUnprocessableEntity(
				apierrors.FromK8sError(err, ServiceBindingResourceType),
				"Referenced app not found. Ensure that the app exists and you have access to it.",
				apierrors.ForbiddenError{},
				apierrors.NotFoundError{},
			)
	}

	cfPackage := message.toCFPackage()
	err = userClient.Create(ctx, cfPackage)
	if err != nil {
		return PackageRecord{}, apierrors.FromK8sError(err, PackageResourceType)
	}

	if packageTypeToLifecycleType[cfPackage.Spec.Type] != cfApp.Spec.Lifecycle.Type {
		return PackageRecord{}, apierrors.NewUnprocessableEntityError(nil, fmt.Sprintf("cannot create %s package for a %s app", cfPackage.Spec.Type, cfApp.Spec.Lifecycle.Type))
	}

	if cfPackage.Spec.Type == "bits" {
		err = r.repositoryCreator.CreateRepository(ctx, r.repositoryRef(cfPackage))
		if err != nil {
			return PackageRecord{}, fmt.Errorf("failed to create package repository: %w", err)
		}
	}

	if isPrivateDockerImage(message) {
		err = createImagePullSecret(ctx, userClient, cfPackage, message)
		if err != nil {
			return PackageRecord{}, fmt.Errorf("failed to build docker image pull secret: %w", err)
		}
	}

	cfPackage, err = r.awaiter.AwaitCondition(ctx, userClient, cfPackage, packages.InitializedConditionType)
	if err != nil {
		return PackageRecord{}, fmt.Errorf("failed waiting for Initialized condition: %w", err)
	}

	return r.cfPackageToPackageRecord(cfPackage), nil
}

func isPrivateDockerImage(message CreatePackageMessage) bool {
	return message.Type == "docker" &&
		message.Data.Username != nil &&
		message.Data.Password != nil
}

func createImagePullSecret(ctx context.Context, userClient client.Client, cfPackage *korifiv1alpha1.CFPackage, message CreatePackageMessage) error {
	ref, err := name.ParseReference(message.Data.Image)
	if err != nil {
		return fmt.Errorf("failed to parse image ref: %w", err)
	}

	imgPullSecret, err := dockercfg.CreateDockerConfigSecret(
		cfPackage.Namespace,
		cfPackage.Name,
		dockercfg.DockerServerConfig{
			Server:   ref.Context().RegistryStr(),
			Username: *message.Data.Username,
			Password: *message.Data.Password,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to generate image pull secret: %w", err)
	}

	err = controllerutil.SetOwnerReference(cfPackage, imgPullSecret, scheme.Scheme)
	if err != nil {
		return fmt.Errorf("failed to set ownership from the package to the image pull secret: %w", err)
	}

	err = userClient.Create(ctx, imgPullSecret)
	if err != nil {
		return fmt.Errorf("failed create the image pull secret: %w", err)
	}

	err = k8s.PatchResource(ctx, userClient, cfPackage, func() {
		cfPackage.Spec.Source.Registry.ImagePullSecrets = []corev1.LocalObjectReference{{Name: imgPullSecret.Name}}
	})
	if err != nil {
		return fmt.Errorf("failed set the package image pull secret: %w", err)
	}

	return nil
}

func (r *PackageRepo) UpdatePackage(ctx context.Context, authInfo authorization.Info, updateMessage UpdatePackageMessage) (PackageRecord, error) {
	ns, err := r.namespaceRetriever.NamespaceFor(ctx, updateMessage.GUID, PackageResourceType)
	if err != nil {
		return PackageRecord{}, err
	}

	userClient, err := r.userClientFactory.BuildClient(authInfo)
	if err != nil {
		return PackageRecord{}, fmt.Errorf("failed to build user client: %w", err)
	}

	cfPackage := &korifiv1alpha1.CFPackage{}

	err = userClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: updateMessage.GUID}, cfPackage)
	if err != nil {
		return PackageRecord{}, fmt.Errorf("failed to get package: %w", apierrors.ForbiddenAsNotFound(apierrors.FromK8sError(err, PackageResourceType)))
	}

	err = k8s.PatchResource(ctx, userClient, cfPackage, func() {
		updateMessage.MetadataPatch.Apply(cfPackage)
	})
	if err != nil {
		return PackageRecord{}, fmt.Errorf("failed to patch package metadata: %w", apierrors.FromK8sError(err, PackageResourceType))
	}

	return r.cfPackageToPackageRecord(cfPackage), nil
}

func (r *PackageRepo) GetPackage(ctx context.Context, authInfo authorization.Info, guid string) (PackageRecord, error) {
	ns, err := r.namespaceRetriever.NamespaceFor(ctx, guid, PackageResourceType)
	if err != nil {
		return PackageRecord{}, err
	}

	userClient, err := r.userClientFactory.BuildClient(authInfo)
	if err != nil {
		return PackageRecord{}, fmt.Errorf("failed to build user k8s client: %w", err)
	}

	cfPackage := new(korifiv1alpha1.CFPackage)
	if err := userClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: guid}, cfPackage); err != nil {
		return PackageRecord{}, fmt.Errorf("failed to get package %q: %w", guid, apierrors.FromK8sError(err, PackageResourceType))
	}

	return r.cfPackageToPackageRecord(cfPackage), nil
}

func (r *PackageRepo) ListPackages(ctx context.Context, authInfo authorization.Info, message ListPackagesMessage) ([]PackageRecord, error) {
	nsList, err := r.namespacePermissions.GetAuthorizedSpaceNamespaces(ctx, authInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces for spaces with user role bindings: %w", err)
	}
	userClient, err := r.userClientFactory.BuildClient(authInfo)
	if err != nil {
		return []PackageRecord{}, fmt.Errorf("failed to build user client: %w", err)
	}

	preds := []func(korifiv1alpha1.CFPackage) bool{
		SetPredicate(message.GUIDs, func(s korifiv1alpha1.CFPackage) string { return s.Name }),
		SetPredicate(message.AppGUIDs, func(s korifiv1alpha1.CFPackage) string { return s.Spec.AppRef.Name }),
	}
	if len(message.States) > 0 {
		stateSet := NewSet(message.States...)
		preds = append(preds, func(p korifiv1alpha1.CFPackage) bool {
			return (stateSet.Includes(PackageStateReady) && meta.IsStatusConditionTrue(p.Status.Conditions, korifiv1alpha1.StatusConditionReady)) ||
				(stateSet.Includes(PackageStateAwaitingUpload) && !meta.IsStatusConditionTrue(p.Status.Conditions, korifiv1alpha1.StatusConditionReady))
		})
	}

	var filteredPackages []korifiv1alpha1.CFPackage
	for ns := range nsList {
		packageList := &korifiv1alpha1.CFPackageList{}
		err = userClient.List(ctx, packageList, client.InNamespace(ns))
		if k8serrors.IsForbidden(err) {
			continue
		}
		if err != nil {
			return []PackageRecord{}, fmt.Errorf("failed to list packages in namespace %s: %w", ns, apierrors.FromK8sError(err, PackageResourceType))
		}
		filteredPackages = append(filteredPackages, Filter(packageList.Items, preds...)...)
	}
	return r.convertToPackageRecords(filteredPackages), nil
}

func (r *PackageRepo) UpdatePackageSource(ctx context.Context, authInfo authorization.Info, message UpdatePackageSourceMessage) (PackageRecord, error) {
	userClient, err := r.userClientFactory.BuildClient(authInfo)
	if err != nil {
		return PackageRecord{}, fmt.Errorf("failed to build user k8s client: %w", err)
	}

	cfPackage := &korifiv1alpha1.CFPackage{}
	if err = userClient.Get(ctx, client.ObjectKey{Name: message.GUID, Namespace: message.SpaceGUID}, cfPackage); err != nil {
		return PackageRecord{}, fmt.Errorf("failed to get cf package: %w", apierrors.FromK8sError(err, PackageResourceType))
	}

	if err = k8s.PatchResource(ctx, userClient, cfPackage, func() {
		cfPackage.Spec.Source.Registry.Image = message.ImageRef
		imagePullSecrets := []corev1.LocalObjectReference{}
		for _, secret := range message.RegistrySecretNames {
			imagePullSecrets = append(imagePullSecrets, corev1.LocalObjectReference{Name: secret})
		}
		cfPackage.Spec.Source.Registry.ImagePullSecrets = imagePullSecrets
	}); err != nil {
		return PackageRecord{}, fmt.Errorf("failed to update package source: %w", apierrors.FromK8sError(err, PackageResourceType))
	}

	cfPackage, err = r.awaiter.AwaitCondition(ctx, userClient, cfPackage, korifiv1alpha1.StatusConditionReady)
	if err != nil {
		return PackageRecord{}, fmt.Errorf("failed awaiting Ready status condition: %w", err)
	}

	record := r.cfPackageToPackageRecord(cfPackage)
	return record, nil
}

func (r *PackageRepo) cfPackageToPackageRecord(cfPackage *korifiv1alpha1.CFPackage) PackageRecord {
	state := PackageStateAwaitingUpload
	if meta.IsStatusConditionTrue(cfPackage.Status.Conditions, korifiv1alpha1.StatusConditionReady) {
		state = PackageStateReady
	}
	return PackageRecord{
		GUID:        cfPackage.Name,
		UID:         cfPackage.UID,
		SpaceGUID:   cfPackage.Namespace,
		Type:        string(cfPackage.Spec.Type),
		AppGUID:     cfPackage.Spec.AppRef.Name,
		State:       state,
		CreatedAt:   cfPackage.CreationTimestamp.Time,
		UpdatedAt:   getLastUpdatedTime(cfPackage),
		Labels:      cfPackage.Labels,
		Annotations: cfPackage.Annotations,
		ImageRef:    r.repositoryRef(cfPackage),
	}
}

func (r *PackageRepo) convertToPackageRecords(packages []korifiv1alpha1.CFPackage) []PackageRecord {
	packageRecords := make([]PackageRecord, 0, len(packages))

	for i := range packages {
		packageRecords = append(packageRecords, r.cfPackageToPackageRecord(&packages[i]))
	}
	return packageRecords
}

func (r *PackageRepo) repositoryRef(cfPackage *korifiv1alpha1.CFPackage) string {
	if cfPackage.Spec.Type == "docker" {
		return cfPackage.Spec.Source.Registry.Image
	}

	return r.repositoryPrefix + cfPackage.Spec.AppRef.Name + "-packages"
}
