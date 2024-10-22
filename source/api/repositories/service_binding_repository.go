package repositories

import (
	"context"
	"fmt"
	"slices"
	"time"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"

	"code.cloudfoundry.org/korifi/api/authorization"
	apierrors "code.cloudfoundry.org/korifi/api/errors"
	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/controllers/webhooks/services/bindings"
	"code.cloudfoundry.org/korifi/controllers/webhooks/validation"
	"code.cloudfoundry.org/korifi/tools"
	"code.cloudfoundry.org/korifi/tools/k8s"

	"github.com/BooleanCat/go-functional/v2/it"
	"github.com/BooleanCat/go-functional/v2/it/itx"
	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	LabelServiceBindingProvisionedService = "servicebinding.io/provisioned-service"
	ServiceBindingResourceType            = "Service Binding"
	ServiceBindingTypeApp                 = "app"
)

type ServiceBindingRepo struct {
	userClientFactory       authorization.UserK8sClientFactory
	namespacePermissions    *authorization.NamespacePermissions
	namespaceRetriever      NamespaceRetriever
	bindingConditionAwaiter Awaiter[*korifiv1alpha1.CFServiceBinding]
}

func NewServiceBindingRepo(
	namespaceRetriever NamespaceRetriever,
	userClientFactory authorization.UserK8sClientFactory,
	namespacePermissions *authorization.NamespacePermissions,
	bindingConditionAwaiter Awaiter[*korifiv1alpha1.CFServiceBinding],
) *ServiceBindingRepo {
	return &ServiceBindingRepo{
		userClientFactory:       userClientFactory,
		namespacePermissions:    namespacePermissions,
		namespaceRetriever:      namespaceRetriever,
		bindingConditionAwaiter: bindingConditionAwaiter,
	}
}

type ServiceBindingRecord struct {
	GUID                string
	Type                string
	Name                *string
	AppGUID             string
	ServiceInstanceGUID string
	SpaceGUID           string
	Labels              map[string]string
	Annotations         map[string]string
	CreatedAt           time.Time
	UpdatedAt           *time.Time
	LastOperation       ServiceBindingLastOperation
}

func (r ServiceBindingRecord) Relationships() map[string]string {
	return map[string]string{
		"app":              r.AppGUID,
		"service_instance": r.ServiceInstanceGUID,
	}
}

type ServiceBindingLastOperation struct {
	Type        string
	State       string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type CreateServiceBindingMessage struct {
	Name                *string
	ServiceInstanceGUID string
	AppGUID             string
	SpaceGUID           string
}

type DeleteServiceBindingMessage struct {
	GUID string
}

type ListServiceBindingsMessage struct {
	AppGUIDs             []string
	ServiceInstanceGUIDs []string
	LabelSelector        string
}

func (m *ListServiceBindingsMessage) matches(serviceBinding korifiv1alpha1.CFServiceBinding) bool {
	return tools.EmptyOrContains(m.ServiceInstanceGUIDs, serviceBinding.Spec.Service.Name) &&
		tools.EmptyOrContains(m.AppGUIDs, serviceBinding.Spec.AppRef.Name)
}

func (m CreateServiceBindingMessage) toCFServiceBinding() *korifiv1alpha1.CFServiceBinding {
	guid := uuid.NewString()
	return &korifiv1alpha1.CFServiceBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      guid,
			Namespace: m.SpaceGUID,
			Labels:    map[string]string{LabelServiceBindingProvisionedService: "true"},
		},
		Spec: korifiv1alpha1.CFServiceBindingSpec{
			DisplayName: m.Name,
			Service: corev1.ObjectReference{
				Kind:       "CFServiceInstance",
				APIVersion: korifiv1alpha1.GroupVersion.Identifier(),
				Name:       m.ServiceInstanceGUID,
			},
			AppRef: corev1.LocalObjectReference{Name: m.AppGUID},
		},
	}
}

type UpdateServiceBindingMessage struct {
	GUID          string
	MetadataPatch MetadataPatch
}

func (r *ServiceBindingRepo) CreateServiceBinding(ctx context.Context, authInfo authorization.Info, message CreateServiceBindingMessage) (ServiceBindingRecord, error) {
	userClient, err := r.userClientFactory.BuildClient(authInfo)
	if err != nil {
		return ServiceBindingRecord{}, fmt.Errorf("failed to build user client: %w", err)
	}

	cfServiceBinding := message.toCFServiceBinding()

	cfApp := new(korifiv1alpha1.CFApp)
	err = userClient.Get(ctx, types.NamespacedName{Name: cfServiceBinding.Spec.AppRef.Name, Namespace: cfServiceBinding.Namespace}, cfApp)
	if err != nil {
		return ServiceBindingRecord{},
			apierrors.AsUnprocessableEntity(
				apierrors.FromK8sError(err, ServiceBindingResourceType),
				"Unable to use app. Ensure that the app exists and you have access to it.",
				apierrors.ForbiddenError{},
				apierrors.NotFoundError{},
			)
	}

	err = userClient.Create(ctx, cfServiceBinding)
	if err != nil {
		if validationError, ok := validation.WebhookErrorToValidationError(err); ok {
			if validationError.Type == bindings.ServiceBindingErrorType {
				return ServiceBindingRecord{}, apierrors.NewUniquenessError(err, validationError.GetMessage())
			}
		}

		return ServiceBindingRecord{}, apierrors.FromK8sError(err, ServiceBindingResourceType)
	}

	cfServiceBinding, err = r.bindingConditionAwaiter.AwaitCondition(ctx, userClient, cfServiceBinding, korifiv1alpha1.StatusConditionReady)
	if err != nil {
		return ServiceBindingRecord{}, err
	}

	return cfServiceBindingToRecord(*cfServiceBinding), err
}

func (r *ServiceBindingRepo) DeleteServiceBinding(ctx context.Context, authInfo authorization.Info, guid string) error {
	userClient, err := r.userClientFactory.BuildClient(authInfo)
	if err != nil {
		return fmt.Errorf("failed to build user client: %w", err)
	}

	namespace, err := r.namespaceRetriever.NamespaceFor(ctx, guid, ServiceBindingResourceType)
	if err != nil {
		return err
	}

	binding := &korifiv1alpha1.CFServiceBinding{}

	err = userClient.Get(ctx, client.ObjectKey{Namespace: namespace, Name: guid}, binding)
	if err != nil {
		return apierrors.ForbiddenAsNotFound(apierrors.FromK8sError(err, ServiceBindingResourceType))
	}

	err = userClient.Delete(ctx, binding)
	if err != nil {
		return apierrors.FromK8sError(err, ServiceBindingResourceType)
	}
	return nil
}

func (r *ServiceBindingRepo) GetServiceBinding(ctx context.Context, authInfo authorization.Info, guid string) (ServiceBindingRecord, error) {
	ns, err := r.namespaceRetriever.NamespaceFor(ctx, guid, ServiceBindingResourceType)
	if err != nil {
		return ServiceBindingRecord{}, err
	}

	userClient, err := r.userClientFactory.BuildClient(authInfo)
	if err != nil {
		return ServiceBindingRecord{}, fmt.Errorf("get-service-binding failed to create user client: %w", err)
	}

	serviceBinding := &korifiv1alpha1.CFServiceBinding{}
	err = userClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: guid}, serviceBinding)
	if err != nil {
		return ServiceBindingRecord{}, apierrors.FromK8sError(err, ServiceBindingResourceType)
	}

	return cfServiceBindingToRecord(*serviceBinding), nil
}

func (r *ServiceBindingRepo) UpdateServiceBinding(ctx context.Context, authInfo authorization.Info, updateMsg UpdateServiceBindingMessage) (ServiceBindingRecord, error) {
	userClient, err := r.userClientFactory.BuildClient(authInfo)
	if err != nil {
		return ServiceBindingRecord{}, fmt.Errorf("failed to create user client: %w", err)
	}

	ns, err := r.namespaceRetriever.NamespaceFor(ctx, updateMsg.GUID, ServiceBindingResourceType)
	if err != nil {
		return ServiceBindingRecord{}, err
	}

	serviceBinding := &korifiv1alpha1.CFServiceBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      updateMsg.GUID,
			Namespace: ns,
		},
	}

	err = userClient.Get(ctx, client.ObjectKeyFromObject(serviceBinding), serviceBinding)
	if err != nil {
		return ServiceBindingRecord{}, fmt.Errorf("failed to get service binding: %w", apierrors.FromK8sError(err, ServiceBindingResourceType))
	}

	err = k8s.PatchResource(ctx, userClient, serviceBinding, func() {
		updateMsg.MetadataPatch.Apply(serviceBinding)
	})
	if err != nil {
		return ServiceBindingRecord{}, fmt.Errorf("failed to patch service binding metadata: %w", apierrors.FromK8sError(err, ServiceBindingResourceType))
	}

	return cfServiceBindingToRecord(*serviceBinding), nil
}

func cfServiceBindingToRecord(binding korifiv1alpha1.CFServiceBinding) ServiceBindingRecord {
	return ServiceBindingRecord{
		GUID:                binding.Name,
		Type:                ServiceBindingTypeApp,
		Name:                binding.Spec.DisplayName,
		AppGUID:             binding.Spec.AppRef.Name,
		ServiceInstanceGUID: binding.Spec.Service.Name,
		SpaceGUID:           binding.Namespace,
		Labels:              binding.Labels,
		Annotations:         binding.Annotations,
		CreatedAt:           binding.CreationTimestamp.Time,
		UpdatedAt:           getLastUpdatedTime(&binding),
		LastOperation: ServiceBindingLastOperation{
			Type:        "create",
			State:       "succeeded",
			Description: nil,
			CreatedAt:   binding.CreationTimestamp.Time,
			UpdatedAt:   getLastUpdatedTime(&binding),
		},
	}
}

// nolint:dupl
func (r *ServiceBindingRepo) ListServiceBindings(ctx context.Context, authInfo authorization.Info, message ListServiceBindingsMessage) ([]ServiceBindingRecord, error) {
	userClient, err := r.userClientFactory.BuildClient(authInfo)
	if err != nil {
		return []ServiceBindingRecord{}, fmt.Errorf("failed to build user client: %w", err)
	}

	authorizedSpaceNamespaces, err := authorizedSpaceNamespaces(ctx, authInfo, r.namespacePermissions)
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces for spaces with user role bindings: %w", err)
	}

	labelSelector, err := labels.Parse(message.LabelSelector)
	if err != nil {
		return []ServiceBindingRecord{}, apierrors.NewUnprocessableEntityError(err, "invalid label selector")
	}

	var serviceBindings []korifiv1alpha1.CFServiceBinding
	for _, ns := range authorizedSpaceNamespaces.Collect() {
		serviceBindingList := new(korifiv1alpha1.CFServiceBindingList)
		err = userClient.List(ctx, serviceBindingList, client.InNamespace(ns), &client.ListOptions{LabelSelector: labelSelector})
		if k8serrors.IsForbidden(err) {
			continue
		}
		if err != nil {
			return []ServiceBindingRecord{}, fmt.Errorf("failed to list service instances in namespace %s: %w",
				ns,
				apierrors.FromK8sError(err, ServiceBindingResourceType),
			)
		}
		serviceBindings = append(serviceBindings, serviceBindingList.Items...)
	}

	filteredServiceBindings := itx.FromSlice(serviceBindings).Filter(message.matches)
	return slices.Collect(it.Map(filteredServiceBindings, cfServiceBindingToRecord)), nil
}
