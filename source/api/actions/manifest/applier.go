package manifest

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"code.cloudfoundry.org/korifi/api/actions/shared"
	"code.cloudfoundry.org/korifi/api/authorization"
	apierrors "code.cloudfoundry.org/korifi/api/errors"
	"code.cloudfoundry.org/korifi/api/payloads"
	"code.cloudfoundry.org/korifi/api/repositories"
	"code.cloudfoundry.org/korifi/api/tools/singleton"
	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
)

type Applier struct {
	appRepo             shared.CFAppRepository
	domainRepo          shared.CFDomainRepository
	processRepo         shared.CFProcessRepository
	routeRepo           shared.CFRouteRepository
	serviceInstanceRepo shared.CFServiceInstanceRepository
	serviceBindingRepo  shared.CFServiceBindingRepository
}

func NewApplier(
	appRepo shared.CFAppRepository,
	domainRepo shared.CFDomainRepository,
	processRepo shared.CFProcessRepository,
	routeRepo shared.CFRouteRepository,
	serviceInstanceRepo shared.CFServiceInstanceRepository,
	serviceBindingRepo shared.CFServiceBindingRepository,
) *Applier {
	return &Applier{
		appRepo:             appRepo,
		domainRepo:          domainRepo,
		processRepo:         processRepo,
		routeRepo:           routeRepo,
		serviceInstanceRepo: serviceInstanceRepo,
		serviceBindingRepo:  serviceBindingRepo,
	}
}

func (a *Applier) Apply(ctx context.Context, authInfo authorization.Info, spaceGUID string, appInfo payloads.ManifestApplication, appState AppState) error {
	appState, err := a.applyApp(ctx, authInfo, spaceGUID, appInfo, appState)
	if err != nil {
		return err
	}

	if err := a.applyProcesses(ctx, authInfo, appInfo, appState); err != nil {
		return err
	}

	if err := a.applyRoutes(ctx, authInfo, appInfo, appState); err != nil {
		return err
	}

	return a.applyServices(ctx, authInfo, appInfo, appState)
}

func (a *Applier) applyApp(
	ctx context.Context,
	authInfo authorization.Info,
	spaceGUID string,
	appInfo payloads.ManifestApplication,
	appState AppState,
) (AppState, error) {
	if appState.App.GUID == "" {
		appRecord, err := a.appRepo.CreateApp(ctx, authInfo, appInfo.ToAppCreateMessage(spaceGUID))
		return AppState{App: appRecord}, err
	} else {
		_, err := a.appRepo.PatchApp(ctx, authInfo, appInfo.ToAppPatchMessage(appState.App.GUID, spaceGUID))
		return appState, err
	}
}

func (a *Applier) applyProcesses(
	ctx context.Context,
	authInfo authorization.Info,
	appInfo payloads.ManifestApplication,
	appState AppState,
) error {
	for _, processInfo := range appInfo.Processes {
		if process, ok := appState.Processes[processInfo.Type]; ok {
			if _, err := a.processRepo.PatchProcess(ctx, authInfo, processInfo.ToProcessPatchMessage(process.GUID, appState.App.SpaceGUID)); err != nil {
				return err
			}
			continue
		}

		if err := a.processRepo.CreateProcess(ctx, authInfo, processInfo.ToProcessCreateMessage(appState.App.GUID, appState.App.SpaceGUID)); err != nil {
			return err
		}

	}

	return nil
}

func (a *Applier) applyRoutes(ctx context.Context, authInfo authorization.Info, appInfo payloads.ManifestApplication, appState AppState) error {
	if appInfo.NoRoute {
		return a.deleteAppDestinations(ctx, authInfo, appState.App.GUID, appState.Routes)
	}

	return a.createOrUpdateRoutes(ctx, authInfo, appInfo, appState)
}

func (a *Applier) createOrUpdateRoutes(ctx context.Context, authInfo authorization.Info, appInfo payloads.ManifestApplication, appState AppState) error {
	for _, route := range appInfo.Routes {
		err := a.createOrUpdateRoute(ctx, authInfo, *route.Route, appState)
		if err != nil {
			return fmt.Errorf("createOrUpdateRoutes: %w", err)
		}
	}

	return nil
}

func (a *Applier) createOrUpdateRoute(ctx context.Context, authInfo authorization.Info, routeString string, appState AppState) error {
	if _, routeExists := appState.Routes[routeString]; routeExists {
		return nil
	}

	hostName, domainName, path := splitRoute(routeString)

	domains, err := a.domainRepo.ListDomains(ctx, authInfo, repositories.ListDomainsMessage{
		Names: []string{domainName},
	})
	if err != nil {
		return fmt.Errorf("failed to list domains: %w", err)
	}

	domain, err := singleton.Get(domains)
	if err != nil {
		return err
	}

	routeRecord, err := a.routeRepo.GetOrCreateRoute(
		ctx,
		authInfo,
		repositories.CreateRouteMessage{
			Host:            hostName,
			Path:            path,
			SpaceGUID:       appState.App.SpaceGUID,
			DomainGUID:      domain.GUID,
			DomainNamespace: domain.Namespace,
			DomainName:      domain.Name,
		})
	if err != nil {
		return fmt.Errorf("getOrCreateRoute: %w", err)
	}

	_, err = a.routeRepo.AddDestinationsToRoute(ctx, authInfo, repositories.AddDestinationsMessage{
		RouteGUID:            routeRecord.GUID,
		SpaceGUID:            routeRecord.SpaceGUID,
		ExistingDestinations: routeRecord.Destinations,
		NewDestinations: []repositories.DesiredDestination{{
			AppGUID:     appState.App.GUID,
			ProcessType: korifiv1alpha1.ProcessTypeWeb,
		}},
	})
	if err != nil {
		return fmt.Errorf("addDestinationsToRoute: %w", err)
	}

	return nil
}

func (a *Applier) deleteAppDestinations(
	ctx context.Context,
	authInfo authorization.Info,
	appGUID string,
	existingAppRoutes map[string]repositories.RouteRecord,
) error {
	for _, route := range existingAppRoutes {
		existingDestinations := route.Destinations

		for _, destination := range route.Destinations {
			if destination.AppGUID != appGUID {
				continue
			}

			var err error
			existingDestinations, err = a.deleteAppDestination(ctx, authInfo, route, destination.GUID, existingDestinations)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *Applier) deleteAppDestination(ctx context.Context, authInfo authorization.Info, route repositories.RouteRecord, destinationGUID string, existingDestinations []repositories.DestinationRecord) ([]repositories.DestinationRecord, error) {
	route, err := a.routeRepo.RemoveDestinationFromRoute(ctx, authInfo, repositories.RemoveDestinationMessage{
		RouteGUID: route.GUID,
		SpaceGUID: route.SpaceGUID,
		GUID:      destinationGUID,
	})
	if err != nil {
		return nil, err
	}

	return route.Destinations, nil
}

func (a *Applier) applyServices(ctx context.Context, authInfo authorization.Info, appInfo payloads.ManifestApplication, appState AppState) error {
	desiredServiceNames := map[string]bool{}
	for _, s := range appInfo.Services {
		desiredServiceNames[s.Name] = true
	}
	for serviceName := range appState.ServiceBindings {
		delete(desiredServiceNames, serviceName)
	}

	if len(desiredServiceNames) == 0 {
		return nil
	}

	serviceInstances, err := a.serviceInstanceRepo.ListServiceInstances(ctx, authInfo, repositories.ListServiceInstanceMessage{
		Names: slices.Collect(maps.Keys(desiredServiceNames)),
	})
	if err != nil {
		return err
	}

	serviceNameToServiceInstance := map[string]repositories.ServiceInstanceRecord{}
	for _, serviceInstance := range serviceInstances {
		serviceNameToServiceInstance[serviceInstance.Name] = serviceInstance
	}

	serviceNameToServiceBinding := map[string]*string{}
	for _, manifestService := range appInfo.Services {
		serviceNameToServiceBinding[manifestService.Name] = manifestService.BindingName
	}

	for serviceName := range desiredServiceNames {
		serviceInstance, ok := serviceNameToServiceInstance[serviceName]
		if !ok {
			return apierrors.NewNotFoundError(
				nil,
				repositories.ServiceInstanceResourceType,
				"application", appInfo.Name,
				"service", serviceName,
			)
		}

		_, err := a.serviceBindingRepo.CreateServiceBinding(ctx, authInfo, repositories.CreateServiceBindingMessage{
			Name:                serviceNameToServiceBinding[serviceName],
			ServiceInstanceGUID: serviceInstance.GUID,
			AppGUID:             appState.App.GUID,
			SpaceGUID:           appState.App.SpaceGUID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func splitRoute(route string) (string, string, string) {
	parts := strings.SplitN(route, ".", 2)
	hostName := parts[0]
	domainAndPath := parts[1]

	parts = strings.SplitN(domainAndPath, "/", 2)
	domain := parts[0]
	var path string
	if len(parts) > 1 {
		path = "/" + parts[1]
	}
	return hostName, domain, path
}
