package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"code.cloudfoundry.org/korifi/api/actions"
	"code.cloudfoundry.org/korifi/api/authorization"
	apierrors "code.cloudfoundry.org/korifi/api/errors"
	"code.cloudfoundry.org/korifi/api/payloads"
	"code.cloudfoundry.org/korifi/api/presenter"
	"code.cloudfoundry.org/korifi/api/repositories"
	"code.cloudfoundry.org/korifi/api/routing"

	"github.com/go-logr/logr"
)

const (
	ProcessPath         = "/v3/processes/{guid}"
	ProcessSidecarsPath = "/v3/processes/{guid}/sidecars"
	ProcessScalePath    = "/v3/processes/{guid}/actions/scale"
	ProcessStatsPath    = "/v3/processes/{guid}/stats"
	ProcessesPath       = "/v3/processes"
)

//counterfeiter:generate -o fake -fake-name CFProcessRepository . CFProcessRepository
type CFProcessRepository interface {
	GetProcess(context.Context, authorization.Info, string) (repositories.ProcessRecord, error)
	ListProcesses(context.Context, authorization.Info, repositories.ListProcessesMessage) ([]repositories.ProcessRecord, error)
	GetProcessByAppTypeAndSpace(context.Context, authorization.Info, string, string, string) (repositories.ProcessRecord, error)
	PatchProcess(context.Context, authorization.Info, repositories.PatchProcessMessage) (repositories.ProcessRecord, error)
	CreateProcess(context.Context, authorization.Info, repositories.CreateProcessMessage) error
	ScaleProcess(ctx context.Context, authInfo authorization.Info, scaleProcessMessage repositories.ScaleProcessMessage) (repositories.ProcessRecord, error)
}

//counterfeiter:generate -o fake -fake-name ProcessStats . ProcessStats
type ProcessStats interface {
	FetchStats(context.Context, authorization.Info, string) ([]actions.PodStatsRecord, error)
}

type Process struct {
	serverURL        url.URL
	processRepo      CFProcessRepository
	processStats     ProcessStats
	requestValidator RequestValidator
}

func NewProcess(
	serverURL url.URL,
	processRepo CFProcessRepository,
	processStatsFetcher ProcessStats,
	requestValidator RequestValidator,
) *Process {
	return &Process{
		serverURL:        serverURL,
		processRepo:      processRepo,
		processStats:     processStatsFetcher,
		requestValidator: requestValidator,
	}
}

func (h *Process) get(r *http.Request) (*routing.Response, error) {
	authInfo, _ := authorization.InfoFromContext(r.Context())
	logger := logr.FromContextOrDiscard(r.Context()).WithName("handlers.process.get")

	processGUID := routing.URLParam(r, "guid")

	process, err := h.processRepo.GetProcess(r.Context(), authInfo, processGUID)
	if err != nil {
		return nil, apierrors.LogAndReturn(logger, apierrors.ForbiddenAsNotFound(err), "Failed to fetch process from Kubernetes", "ProcessGUID", processGUID)
	}

	return routing.NewResponse(http.StatusOK).WithBody(presenter.ForProcess(process, h.serverURL)), nil
}

func (h *Process) getSidecars(r *http.Request) (*routing.Response, error) {
	authInfo, _ := authorization.InfoFromContext(r.Context())
	logger := logr.FromContextOrDiscard(r.Context()).WithName("handlers.process.get-sidecars")

	processGUID := routing.URLParam(r, "guid")

	_, err := h.processRepo.GetProcess(r.Context(), authInfo, processGUID)
	if err != nil {
		return nil, apierrors.LogAndReturn(logger, apierrors.ForbiddenAsNotFound(err), "Failed to fetch process from Kubernetes", "ProcessGUID", processGUID)
	}

	return routing.NewResponse(http.StatusOK).WithBody(map[string]interface{}{
		"pagination": map[string]interface{}{
			"total_results": 0,
			"total_pages":   1,
			"first": map[string]interface{}{
				"href": fmt.Sprintf("%s/v3/processes/%s/sidecars", h.serverURL.String(), processGUID),
			},
			"last": map[string]interface{}{
				"href": fmt.Sprintf("%s/v3/processes/%s/sidecars", h.serverURL.String(), processGUID),
			},
			"next":     nil,
			"previous": nil,
		},
		"resources": []string{},
	}), nil
}

func (h *Process) scale(r *http.Request) (*routing.Response, error) {
	authInfo, _ := authorization.InfoFromContext(r.Context())
	logger := logr.FromContextOrDiscard(r.Context()).WithName("handlers.process.scale")

	processGUID := routing.URLParam(r, "guid")

	var payload payloads.ProcessScale
	if err := h.requestValidator.DecodeAndValidateJSONPayload(r, &payload); err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "failed to decode payload")
	}

	process, err := h.processRepo.GetProcess(r.Context(), authInfo, processGUID)
	if err != nil {
		return nil, apierrors.ForbiddenAsNotFound(err)
	}

	processRecord, err := h.processRepo.ScaleProcess(r.Context(), authInfo, repositories.ScaleProcessMessage{
		GUID:               process.GUID,
		SpaceGUID:          process.SpaceGUID,
		ProcessScaleValues: payload.ToRecord(),
	})
	if err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "failed to scale process", "processGUID", processGUID)
	}

	return routing.NewResponse(http.StatusOK).WithBody(presenter.ForProcess(processRecord, h.serverURL)), nil
}

func (h *Process) getStats(r *http.Request) (*routing.Response, error) {
	authInfo, _ := authorization.InfoFromContext(r.Context())
	logger := logr.FromContextOrDiscard(r.Context()).WithName("handlers.process.get-stats")

	processGUID := routing.URLParam(r, "guid")

	records, err := h.processStats.FetchStats(r.Context(), authInfo, processGUID)
	if err != nil {
		return nil, apierrors.LogAndReturn(logger, apierrors.ForbiddenAsNotFound(err), "Failed to get process stats from Kubernetes", "ProcessGUID", processGUID)
	}

	return routing.NewResponse(http.StatusOK).WithBody(presenter.ForProcessStats(records)), nil
}

func (h *Process) list(r *http.Request) (*routing.Response, error) { //nolint:dupl
	authInfo, _ := authorization.InfoFromContext(r.Context())
	logger := logr.FromContextOrDiscard(r.Context()).WithName("handlers.process.list")

	processListFilter := new(payloads.ProcessList)
	err := h.requestValidator.DecodeAndValidateURLValues(r, processListFilter)
	if err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "Unable to decode request query parameters")
	}

	processList, err := h.processRepo.ListProcesses(r.Context(), authInfo, processListFilter.ToMessage())
	if err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "Failed to fetch processes(s) from Kubernetes")
	}

	return routing.NewResponse(http.StatusOK).WithBody(presenter.ForProcessList(processList, h.serverURL, *r.URL)), nil
}

func (h *Process) update(r *http.Request) (*routing.Response, error) {
	authInfo, _ := authorization.InfoFromContext(r.Context())
	logger := logr.FromContextOrDiscard(r.Context()).WithName("handlers.process.update")

	processGUID := routing.URLParam(r, "guid")

	var payload payloads.ProcessPatch
	if err := h.requestValidator.DecodeAndValidateJSONPayload(r, &payload); err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "failed to decode json payload")
	}

	process, err := h.processRepo.GetProcess(r.Context(), authInfo, processGUID)
	if err != nil {
		return nil, apierrors.LogAndReturn(logger, apierrors.ForbiddenAsNotFound(err), "Failed to get process from Kubernetes", "ProcessGUID", processGUID)
	}

	updatedProcess, err := h.processRepo.PatchProcess(r.Context(), authInfo, payload.ToProcessPatchMessage(processGUID, process.SpaceGUID))
	if err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "Failed to patch process from Kubernetes", "ProcessGUID", processGUID)
	}

	return routing.NewResponse(http.StatusOK).WithBody(presenter.ForProcess(updatedProcess, h.serverURL)), nil
}

func (h *Process) UnauthenticatedRoutes() []routing.Route {
	return nil
}

func (h *Process) AuthenticatedRoutes() []routing.Route {
	return []routing.Route{
		{Method: "GET", Pattern: ProcessPath, Handler: h.get},
		{Method: "GET", Pattern: ProcessSidecarsPath, Handler: h.getSidecars},
		{Method: "POST", Pattern: ProcessScalePath, Handler: h.scale},
		{Method: "GET", Pattern: ProcessStatsPath, Handler: h.getStats},
		{Method: "GET", Pattern: ProcessesPath, Handler: h.list},
		{Method: "PATCH", Pattern: ProcessPath, Handler: h.update},
	}
}
