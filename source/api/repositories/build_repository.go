package repositories

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"code.cloudfoundry.org/korifi/api/authorization"
	apierrors "code.cloudfoundry.org/korifi/api/errors"
	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"

	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	k8sclient "k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	BuildStateStaging = "STAGING"
	BuildStateStaged  = "STAGED"
	BuildStateFailed  = "FAILED"

	StagingConditionType   = "Staging"
	SucceededConditionType = "Succeeded"

	BuildResourceType    = "Build"
	stagingLogSourceType = "STG"
)

type BuildRecord struct {
	GUID            string
	State           string
	CreatedAt       time.Time
	UpdatedAt       *time.Time
	StagingErrorMsg string
	StagingMemoryMB int
	StagingDiskMB   int
	Lifecycle       Lifecycle
	PackageGUID     string
	DropletGUID     string
	AppGUID         string
	Labels          map[string]string
	Annotations     map[string]string
	ImageRef        string
}

type LogRecord struct {
	Message   string
	Timestamp int64
	Header    string
	Tags      map[string]string
}

type BuildRepo struct {
	namespaceRetriever NamespaceRetriever
	userClientFactory  authorization.UserK8sClientFactory
}

func NewBuildRepo(
	namespaceRetriever NamespaceRetriever,
	userClientFactory authorization.UserK8sClientFactory,
) *BuildRepo {
	return &BuildRepo{
		namespaceRetriever: namespaceRetriever,
		userClientFactory:  userClientFactory,
	}
}

func (b *BuildRepo) GetBuild(ctx context.Context, authInfo authorization.Info, buildGUID string) (BuildRecord, error) {
	ns, err := b.namespaceRetriever.NamespaceFor(ctx, buildGUID, BuildResourceType)
	if err != nil {
		return BuildRecord{}, err
	}

	userClient, err := b.userClientFactory.BuildClient(authInfo)
	if err != nil {
		return BuildRecord{}, fmt.Errorf("get-build failed to build user client: %w", err)
	}

	build := korifiv1alpha1.CFBuild{}
	if err := userClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: buildGUID}, &build); err != nil {
		return BuildRecord{}, fmt.Errorf("failed to get build: %w", apierrors.FromK8sError(err, BuildResourceType))
	}

	return b.cfBuildToBuildRecord(build), nil
}

func (b *BuildRepo) GetLatestBuildByAppGUID(ctx context.Context, authInfo authorization.Info, spaceGUID string, appGUID string) (BuildRecord, error) {
	userClient, err := b.userClientFactory.BuildClient(authInfo)
	if err != nil { // Untested
		return BuildRecord{}, err
	}
	labelSelector, err := labels.ValidatedSelectorFromSet(map[string]string{
		korifiv1alpha1.CFAppGUIDLabelKey: appGUID,
	})
	if err != nil { // Untested
		return BuildRecord{}, err
	}

	listOpts := &client.ListOptions{Namespace: spaceGUID, LabelSelector: labelSelector}
	buildList := &korifiv1alpha1.CFBuildList{}

	err = userClient.List(ctx, buildList, listOpts)
	if err != nil {
		return BuildRecord{}, apierrors.FromK8sError(err, BuildResourceType)
	}

	if len(buildList.Items) == 0 {
		return BuildRecord{}, apierrors.NewNotFoundError(fmt.Errorf("builds for app %q in space %q not found", appGUID, spaceGUID), BuildResourceType)
	}

	sortedBuilds := orderBuilds(buildList.Items)

	return b.cfBuildToBuildRecord(sortedBuilds[0]), nil
}

func orderBuilds(builds []korifiv1alpha1.CFBuild) []korifiv1alpha1.CFBuild {
	sort.Slice(builds, func(i, j int) bool {
		return !builds[i].CreationTimestamp.Before(&builds[j].CreationTimestamp)
	})
	return builds
}

func (b *BuildRepo) GetBuildLogs(ctx context.Context, authInfo authorization.Info, spaceGUID string, buildGUID string) ([]LogRecord, error) {
	userClient, err := b.userClientFactory.BuildK8sClient(authInfo)
	if err != nil { // Untested
		return []LogRecord{}, err
	}
	logWriter := new(strings.Builder)
	err = NewBuildLogsClient(userClient).GetImageLogs(ctx, logWriter, buildGUID, spaceGUID)
	if err != nil {
		return nil, err
	}
	buildLogs := strings.Split(strings.Trim(logWriter.String(), "\n"), "\n")
	toReturn := make([]LogRecord, 0, len(buildLogs))

	for _, log := range buildLogs {
		logLine, logTime, _ := parseRFC3339NanoTime(log)

		toReturn = append(toReturn, LogRecord{
			Message:   logLine,
			Timestamp: logTime,
			Tags: map[string]string{
				"source_type": stagingLogSourceType,
			},
		})
	}

	return toReturn, nil
}

func (b *BuildRepo) cfBuildToBuildRecord(cfBuild korifiv1alpha1.CFBuild) BuildRecord {
	toReturn := BuildRecord{
		GUID:            cfBuild.Name,
		State:           BuildStateStaging,
		CreatedAt:       cfBuild.CreationTimestamp.Time,
		UpdatedAt:       getLastUpdatedTime(&cfBuild),
		StagingErrorMsg: "",
		StagingMemoryMB: cfBuild.Spec.StagingMemoryMB,
		StagingDiskMB:   cfBuild.Spec.StagingDiskMB,
		Lifecycle: Lifecycle{
			Type: string(cfBuild.Spec.Lifecycle.Type),
			Data: LifecycleData{
				Buildpacks: []string{},
				Stack:      cfBuild.Spec.Lifecycle.Data.Stack,
			},
		},
		PackageGUID: cfBuild.Spec.PackageRef.Name,
		DropletGUID: "",
		AppGUID:     cfBuild.Spec.AppRef.Name,
		Labels:      cfBuild.Labels,
		Annotations: cfBuild.Annotations,
	}

	if cfBuild.Spec.Lifecycle.Type == "docker" {
		toReturn.Lifecycle.Data = LifecycleData{}
	}

	if cfBuild.Spec.Lifecycle.Data.Buildpacks != nil {
		toReturn.Lifecycle.Data.Buildpacks = cfBuild.Spec.Lifecycle.Data.Buildpacks
	}

	stagingStatus := getConditionValue(&cfBuild.Status.Conditions, StagingConditionType)
	succeededStatus := getConditionValue(&cfBuild.Status.Conditions, SucceededConditionType)
	// TODO: Consider moving this logic to CRDs repo in case Status Conditions change later?
	if stagingStatus == metav1.ConditionFalse {
		switch succeededStatus {
		case metav1.ConditionTrue:
			toReturn.State = BuildStateStaged
			toReturn.DropletGUID = cfBuild.Name
		case metav1.ConditionFalse:
			toReturn.State = BuildStateFailed
			conditionStatus := meta.FindStatusCondition(cfBuild.Status.Conditions, SucceededConditionType)
			toReturn.StagingErrorMsg = conditionStatus.Message
		}
	}

	return toReturn
}

func (b *BuildRepo) CreateBuild(ctx context.Context, authInfo authorization.Info, message CreateBuildMessage) (BuildRecord, error) {
	cfBuild := message.toCFBuild()
	userClient, err := b.userClientFactory.BuildClient(authInfo)
	if err != nil {
		return BuildRecord{}, fmt.Errorf("failed to build user k8s client: %w", err)
	}

	if err := userClient.Create(ctx, &cfBuild); err != nil {
		return BuildRecord{}, apierrors.FromK8sError(err, BuildResourceType)
	}

	return b.cfBuildToBuildRecord(cfBuild), nil
}

type CreateBuildMessage struct {
	AppGUID         string
	PackageGUID     string
	SpaceGUID       string
	StagingMemoryMB int
	StagingDiskMB   int
	Lifecycle       Lifecycle
	Labels          map[string]string
	Annotations     map[string]string
}

func (m CreateBuildMessage) toCFBuild() korifiv1alpha1.CFBuild {
	guid := uuid.NewString()
	return korifiv1alpha1.CFBuild{
		ObjectMeta: metav1.ObjectMeta{
			Name:        guid,
			Namespace:   m.SpaceGUID,
			Labels:      m.Labels,
			Annotations: m.Annotations,
		},
		Spec: korifiv1alpha1.CFBuildSpec{
			PackageRef: corev1.LocalObjectReference{
				Name: m.PackageGUID,
			},
			AppRef: corev1.LocalObjectReference{
				Name: m.AppGUID,
			},
			StagingMemoryMB: m.StagingMemoryMB,
			StagingDiskMB:   m.StagingDiskMB,
			Lifecycle: korifiv1alpha1.Lifecycle{
				Type: korifiv1alpha1.LifecycleType(m.Lifecycle.Type),
				Data: korifiv1alpha1.LifecycleData{
					Buildpacks: m.Lifecycle.Data.Buildpacks,
					Stack:      m.Lifecycle.Data.Stack,
				},
			},
		},
	}
}

// This code was copied from https://github.com/pivotal/kpack/blob/37764d9940b2f45b1c7acb6c2cd33cf4f489c03b/pkg/logs/build_logs.go
// with minor modifications.
type readyContainer struct {
	podName       string
	containerName string
	namespace     string
}

type BuildLogsClient struct {
	k8sClient k8sclient.Interface
	processed map[readyContainer]interface{}
}

func NewBuildLogsClient(k8sClient k8sclient.Interface) *BuildLogsClient {
	return &BuildLogsClient{
		k8sClient: k8sClient,
		processed: make(map[readyContainer]interface{}),
	}
}

const BuildWorkloadLabelKey = "korifi.cloudfoundry.org/build-workload-name"

func (c *BuildLogsClient) GetImageLogs(ctx context.Context, writer io.Writer, buildGUID, namespace string) error {
	return c.getPodLogs(ctx, writer, namespace, metav1.ListOptions{
		Watch:         false,
		LabelSelector: fmt.Sprintf("%s=%s", BuildWorkloadLabelKey, buildGUID),
	}, false)
}

func (c *BuildLogsClient) getPodLogs(ctx context.Context, writer io.Writer, namespace string, listOptions metav1.ListOptions, follow bool) error {
	readyContainers, err := c.getContainers(ctx, namespace, listOptions)
	if err != nil {
		return err
	}

	for _, container := range readyContainers {
		err := c.streamLogsForContainer(ctx, writer, container, follow)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *BuildLogsClient) getContainers(ctx context.Context, namespace string, listOptions metav1.ListOptions) ([]readyContainer, error) {
	readyContainers := make([]readyContainer, 0)
	pods, err := c.k8sClient.CoreV1().Pods(namespace).List(ctx, listOptions)
	if err != nil {
		return nil, err
	}

	for _, pod := range pods.Items {
		for _, c := range pod.Status.InitContainerStatuses {
			if c.State.Waiting == nil {
				readyContainers = append(readyContainers, readyContainer{
					podName:       pod.Name,
					containerName: c.Name,
					namespace:     pod.Namespace,
				})
			}
		}

		for _, c := range pod.Status.ContainerStatuses {
			if c.State.Waiting == nil {
				readyContainers = append(readyContainers, readyContainer{
					podName:       pod.Name,
					containerName: c.Name,
					namespace:     pod.Namespace,
				})
			}
		}
	}

	return readyContainers, nil
}

func (c *BuildLogsClient) streamLogsForContainer(ctx context.Context, writer io.Writer, readyContainer readyContainer, follow bool) error {
	if _, alreadyProcessed := c.processed[readyContainer]; alreadyProcessed {
		return nil
	}
	c.processed[readyContainer] = nil

	logReadCloser, err := c.k8sClient.CoreV1().Pods(readyContainer.namespace).GetLogs(readyContainer.podName, &corev1.PodLogOptions{
		Container:  readyContainer.containerName,
		Follow:     follow,
		Timestamps: true,
	}).Stream(ctx)
	if err != nil {
		return err
	}
	defer logReadCloser.Close()

	r := bufio.NewReader(logReadCloser)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			line, err := r.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					return nil
				}

				return err
			}

			_, err = writer.Write(line)
			if err != nil {
				return err
			}
		}
	}
}
