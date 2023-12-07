package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"code.cloudfoundry.org/korifi/api/authorization"
	apierrors "code.cloudfoundry.org/korifi/api/errors"
	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/tools/k8s"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ProcessResourceType = "Process"
)

func NewProcessRepo(namespaceRetriever NamespaceRetriever, userClientFactory authorization.UserK8sClientFactory, namespacePermissions *authorization.NamespacePermissions) *ProcessRepo {
	return &ProcessRepo{
		namespaceRetriever:   namespaceRetriever,
		clientFactory:        userClientFactory,
		namespacePermissions: namespacePermissions,
	}
}

type ProcessRepo struct {
	namespaceRetriever   NamespaceRetriever
	clientFactory        authorization.UserK8sClientFactory
	namespacePermissions *authorization.NamespacePermissions
}

type ProcessRecord struct {
	GUID             string
	SpaceGUID        string
	AppGUID          string
	Type             string
	Command          string
	DesiredInstances int
	MemoryMB         int64
	DiskQuotaMB      int64
	HealthCheck      HealthCheck
	Labels           map[string]string
	Annotations      map[string]string
	CreatedAt        time.Time
	UpdatedAt        *time.Time
}

type HealthCheck struct {
	Type string
	Data HealthCheckData
}

type HealthCheckData struct {
	HTTPEndpoint             string
	InvocationTimeoutSeconds int64
	TimeoutSeconds           int64
}

type ScaleProcessMessage struct {
	GUID      string
	SpaceGUID string
	ProcessScaleValues
}

type ProcessScaleValues struct {
	Instances *int
	MemoryMB  *int64
	DiskMB    *int64
}

type CreateProcessMessage struct {
	AppGUID          string
	SpaceGUID        string
	Type             string
	Command          string
	DiskQuotaMB      int64
	HealthCheck      HealthCheck
	DesiredInstances *int
	MemoryMB         int64
}

type PatchProcessMessage struct {
	SpaceGUID                           string
	ProcessGUID                         string
	Command                             *string
	DiskQuotaMB                         *int64
	HealthCheckHTTPEndpoint             *string
	HealthCheckInvocationTimeoutSeconds *int64
	HealthCheckTimeoutSeconds           *int64
	HealthCheckType                     *string
	DesiredInstances                    *int
	MemoryMB                            *int64
	MetadataPatch                       *MetadataPatch
}

type ListProcessesMessage struct {
	AppGUIDs  []string
	SpaceGUID string
}

func (r *ProcessRepo) GetProcess(ctx context.Context, authInfo authorization.Info, processGUID string) (ProcessRecord, error) {
	ns, err := r.namespaceRetriever.NamespaceFor(ctx, processGUID, ProcessResourceType)
	if err != nil {
		return ProcessRecord{}, err
	}

	userClient, err := r.clientFactory.BuildClient(authInfo)
	if err != nil {
		return ProcessRecord{}, fmt.Errorf("get-process: failed to build user k8s client: %w", err)
	}

	var process korifiv1alpha1.CFProcess
	err = userClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: processGUID}, &process)
	if err != nil {
		return ProcessRecord{}, fmt.Errorf("failed to get process %q: %w", processGUID, apierrors.FromK8sError(err, ProcessResourceType))
	}

	return cfProcessToProcessRecord(process), nil
}

func (r *ProcessRepo) ListProcesses(ctx context.Context, authInfo authorization.Info, message ListProcessesMessage) ([]ProcessRecord, error) {
	nsList, err := r.namespacePermissions.GetAuthorizedSpaceNamespaces(ctx, authInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces for spaces with user role bindings: %w", err)
	}

	userClient, err := r.clientFactory.BuildClient(authInfo)
	if err != nil {
		return []ProcessRecord{}, fmt.Errorf("get-process: failed to build user k8s client: %w", err)
	}

	preds := []func(korifiv1alpha1.CFProcess) bool{
		SetPredicate(message.AppGUIDs, func(s korifiv1alpha1.CFProcess) string { return s.Spec.AppRef.Name }),
	}

	processList := &korifiv1alpha1.CFProcessList{}
	var matches []korifiv1alpha1.CFProcess
	for ns := range nsList {
		if message.SpaceGUID != "" && message.SpaceGUID != ns {
			continue
		}
		err = userClient.List(ctx, processList, client.InNamespace(ns))
		if k8serrors.IsForbidden(err) {
			continue
		}
		if err != nil {
			return []ProcessRecord{}, apierrors.FromK8sError(err, ProcessResourceType)
		}
		allProcesses := processList.Items
		matches = append(matches, Filter(allProcesses, preds...)...)
	}

	return returnProcesses(matches)
}

func (r *ProcessRepo) ScaleProcess(ctx context.Context, authInfo authorization.Info, scaleProcessMessage ScaleProcessMessage) (ProcessRecord, error) {
	userClient, err := r.clientFactory.BuildClient(authInfo)
	if err != nil {
		return ProcessRecord{}, fmt.Errorf("get-process: failed to build user k8s client: %w", err)
	}

	cfProcess := &korifiv1alpha1.CFProcess{
		ObjectMeta: metav1.ObjectMeta{
			Name:      scaleProcessMessage.GUID,
			Namespace: scaleProcessMessage.SpaceGUID,
		},
	}
	err = k8s.PatchResource(ctx, userClient, cfProcess, func() {
		if scaleProcessMessage.Instances != nil {
			cfProcess.Spec.DesiredInstances = scaleProcessMessage.Instances
		}
		if scaleProcessMessage.MemoryMB != nil {
			cfProcess.Spec.MemoryMB = *scaleProcessMessage.MemoryMB
		}
		if scaleProcessMessage.DiskMB != nil {
			cfProcess.Spec.DiskQuotaMB = *scaleProcessMessage.DiskMB
		}
	})
	if err != nil {
		return ProcessRecord{}, fmt.Errorf("failed to scale process %q: %w", scaleProcessMessage.GUID, apierrors.FromK8sError(err, ProcessResourceType))
	}

	return cfProcessToProcessRecord(*cfProcess), nil
}

func (r *ProcessRepo) CreateProcess(ctx context.Context, authInfo authorization.Info, message CreateProcessMessage) error {
	userClient, err := r.clientFactory.BuildClient(authInfo)
	if err != nil {
		return fmt.Errorf("get-process: failed to build user k8s client: %w", err)
	}

	process := &korifiv1alpha1.CFProcess{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: message.SpaceGUID,
		},
		Spec: korifiv1alpha1.CFProcessSpec{
			AppRef:      corev1.LocalObjectReference{Name: message.AppGUID},
			ProcessType: message.Type,
			Command:     message.Command,
			HealthCheck: korifiv1alpha1.HealthCheck{
				Type: korifiv1alpha1.HealthCheckType(message.HealthCheck.Type),
				Data: korifiv1alpha1.HealthCheckData(message.HealthCheck.Data),
			},
			DesiredInstances: message.DesiredInstances,
			MemoryMB:         message.MemoryMB,
			DiskQuotaMB:      message.DiskQuotaMB,
		},
	}
	process.SetStableName(message.AppGUID)
	err = userClient.Create(ctx, process)
	return apierrors.FromK8sError(err, ProcessResourceType)
}

func (r *ProcessRepo) GetProcessByAppTypeAndSpace(ctx context.Context, authInfo authorization.Info, appGUID, processType, spaceGUID string) (ProcessRecord, error) {
	// Could narrow down process results via AppGUID label, but that is set up by a webhook that isn't configured in our integration tests
	// For now, don't use labels
	userClient, err := r.clientFactory.BuildClient(authInfo)
	if err != nil {
		return ProcessRecord{}, fmt.Errorf("get-process-by-app-type-and-space: failed to build user k8s client: %w", err)
	}

	var processList korifiv1alpha1.CFProcessList
	err = userClient.List(ctx, &processList, client.InNamespace(spaceGUID))
	if err != nil {
		return ProcessRecord{}, apierrors.FromK8sError(err, ProcessResourceType)
	}

	var matches []korifiv1alpha1.CFProcess
	for _, process := range processList.Items {
		if process.Spec.AppRef.Name == appGUID && process.Spec.ProcessType == processType {
			matches = append(matches, process)
		}
	}

	return returnProcess(matches)
}

func (r *ProcessRepo) PatchProcess(ctx context.Context, authInfo authorization.Info, message PatchProcessMessage) (ProcessRecord, error) {
	userClient, err := r.clientFactory.BuildClient(authInfo)
	if err != nil {
		return ProcessRecord{}, fmt.Errorf("failed to build user client: %w", err)
	}

	updatedProcess := &korifiv1alpha1.CFProcess{
		ObjectMeta: metav1.ObjectMeta{
			Name:      message.ProcessGUID,
			Namespace: message.SpaceGUID,
		},
	}
	err = k8s.PatchResource(ctx, userClient, updatedProcess, func() {
		if message.Command != nil {
			updatedProcess.Spec.Command = *message.Command
		}
		if message.DesiredInstances != nil {
			updatedProcess.Spec.DesiredInstances = message.DesiredInstances
		}
		if message.MemoryMB != nil {
			updatedProcess.Spec.MemoryMB = *message.MemoryMB
		}
		if message.DiskQuotaMB != nil {
			updatedProcess.Spec.DiskQuotaMB = *message.DiskQuotaMB
		}
		if message.HealthCheckType != nil {
			// TODO: how do we handle when the type changes? Clear the HTTPEndpoint when type != http? Should we require the endpoint when type == http?
			updatedProcess.Spec.HealthCheck.Type = korifiv1alpha1.HealthCheckType(*message.HealthCheckType)
		}
		if message.HealthCheckHTTPEndpoint != nil {
			updatedProcess.Spec.HealthCheck.Data.HTTPEndpoint = *message.HealthCheckHTTPEndpoint
		}
		if message.HealthCheckInvocationTimeoutSeconds != nil {
			updatedProcess.Spec.HealthCheck.Data.InvocationTimeoutSeconds = *message.HealthCheckInvocationTimeoutSeconds
		}
		if message.HealthCheckTimeoutSeconds != nil {
			updatedProcess.Spec.HealthCheck.Data.TimeoutSeconds = *message.HealthCheckTimeoutSeconds
		}
		if message.MetadataPatch != nil {
			message.MetadataPatch.Apply(updatedProcess)
		}
	})
	if err != nil {
		return ProcessRecord{}, apierrors.FromK8sError(err, ProcessResourceType)
	}

	return cfProcessToProcessRecord(*updatedProcess), nil
}

func returnProcess(processes []korifiv1alpha1.CFProcess) (ProcessRecord, error) {
	if len(processes) == 0 {
		return ProcessRecord{}, apierrors.NewNotFoundError(nil, ProcessResourceType)
	}
	if len(processes) > 1 {
		return ProcessRecord{}, errors.New("duplicate processes exist")
	}

	return cfProcessToProcessRecord(processes[0]), nil
}

func returnProcesses(processes []korifiv1alpha1.CFProcess) ([]ProcessRecord, error) {
	processRecords := make([]ProcessRecord, 0, len(processes))
	for _, process := range processes {
		processRecord := cfProcessToProcessRecord(process)
		processRecords = append(processRecords, processRecord)
	}

	return processRecords, nil
}

func cfProcessToProcessRecord(cfProcess korifiv1alpha1.CFProcess) ProcessRecord {
	cmd := cfProcess.Spec.Command
	if cmd == "" {
		cmd = cfProcess.Spec.DetectedCommand
	}

	return ProcessRecord{
		GUID:             cfProcess.Name,
		SpaceGUID:        cfProcess.Namespace,
		AppGUID:          cfProcess.Spec.AppRef.Name,
		Type:             cfProcess.Spec.ProcessType,
		Command:          cmd,
		DesiredInstances: *cfProcess.Spec.DesiredInstances,
		MemoryMB:         cfProcess.Spec.MemoryMB,
		DiskQuotaMB:      cfProcess.Spec.DiskQuotaMB,
		HealthCheck: HealthCheck{
			Type: string(cfProcess.Spec.HealthCheck.Type),
			Data: HealthCheckData{
				HTTPEndpoint:             cfProcess.Spec.HealthCheck.Data.HTTPEndpoint,
				InvocationTimeoutSeconds: cfProcess.Spec.HealthCheck.Data.InvocationTimeoutSeconds,
				TimeoutSeconds:           cfProcess.Spec.HealthCheck.Data.TimeoutSeconds,
			},
		},
		Labels:      cfProcess.Labels,
		Annotations: cfProcess.Annotations,
		CreatedAt:   cfProcess.CreationTimestamp.Time,
		UpdatedAt:   getLastUpdatedTime(&cfProcess),
	}
}
