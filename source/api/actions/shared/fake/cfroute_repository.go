// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"context"
	"sync"

	"code.cloudfoundry.org/korifi/api/actions/shared"
	"code.cloudfoundry.org/korifi/api/authorization"
	"code.cloudfoundry.org/korifi/api/repositories"
)

type CFRouteRepository struct {
	AddDestinationsToRouteStub        func(context.Context, authorization.Info, repositories.AddDestinationsToRouteMessage) (repositories.RouteRecord, error)
	addDestinationsToRouteMutex       sync.RWMutex
	addDestinationsToRouteArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.AddDestinationsToRouteMessage
	}
	addDestinationsToRouteReturns struct {
		result1 repositories.RouteRecord
		result2 error
	}
	addDestinationsToRouteReturnsOnCall map[int]struct {
		result1 repositories.RouteRecord
		result2 error
	}
	GetOrCreateRouteStub        func(context.Context, authorization.Info, repositories.CreateRouteMessage) (repositories.RouteRecord, error)
	getOrCreateRouteMutex       sync.RWMutex
	getOrCreateRouteArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateRouteMessage
	}
	getOrCreateRouteReturns struct {
		result1 repositories.RouteRecord
		result2 error
	}
	getOrCreateRouteReturnsOnCall map[int]struct {
		result1 repositories.RouteRecord
		result2 error
	}
	ListRoutesForAppStub        func(context.Context, authorization.Info, string, string) ([]repositories.RouteRecord, error)
	listRoutesForAppMutex       sync.RWMutex
	listRoutesForAppArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
		arg4 string
	}
	listRoutesForAppReturns struct {
		result1 []repositories.RouteRecord
		result2 error
	}
	listRoutesForAppReturnsOnCall map[int]struct {
		result1 []repositories.RouteRecord
		result2 error
	}
	RemoveDestinationFromRouteStub        func(context.Context, authorization.Info, repositories.RemoveDestinationFromRouteMessage) (repositories.RouteRecord, error)
	removeDestinationFromRouteMutex       sync.RWMutex
	removeDestinationFromRouteArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.RemoveDestinationFromRouteMessage
	}
	removeDestinationFromRouteReturns struct {
		result1 repositories.RouteRecord
		result2 error
	}
	removeDestinationFromRouteReturnsOnCall map[int]struct {
		result1 repositories.RouteRecord
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *CFRouteRepository) AddDestinationsToRoute(arg1 context.Context, arg2 authorization.Info, arg3 repositories.AddDestinationsToRouteMessage) (repositories.RouteRecord, error) {
	fake.addDestinationsToRouteMutex.Lock()
	ret, specificReturn := fake.addDestinationsToRouteReturnsOnCall[len(fake.addDestinationsToRouteArgsForCall)]
	fake.addDestinationsToRouteArgsForCall = append(fake.addDestinationsToRouteArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.AddDestinationsToRouteMessage
	}{arg1, arg2, arg3})
	stub := fake.AddDestinationsToRouteStub
	fakeReturns := fake.addDestinationsToRouteReturns
	fake.recordInvocation("AddDestinationsToRoute", []interface{}{arg1, arg2, arg3})
	fake.addDestinationsToRouteMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFRouteRepository) AddDestinationsToRouteCallCount() int {
	fake.addDestinationsToRouteMutex.RLock()
	defer fake.addDestinationsToRouteMutex.RUnlock()
	return len(fake.addDestinationsToRouteArgsForCall)
}

func (fake *CFRouteRepository) AddDestinationsToRouteCalls(stub func(context.Context, authorization.Info, repositories.AddDestinationsToRouteMessage) (repositories.RouteRecord, error)) {
	fake.addDestinationsToRouteMutex.Lock()
	defer fake.addDestinationsToRouteMutex.Unlock()
	fake.AddDestinationsToRouteStub = stub
}

func (fake *CFRouteRepository) AddDestinationsToRouteArgsForCall(i int) (context.Context, authorization.Info, repositories.AddDestinationsToRouteMessage) {
	fake.addDestinationsToRouteMutex.RLock()
	defer fake.addDestinationsToRouteMutex.RUnlock()
	argsForCall := fake.addDestinationsToRouteArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFRouteRepository) AddDestinationsToRouteReturns(result1 repositories.RouteRecord, result2 error) {
	fake.addDestinationsToRouteMutex.Lock()
	defer fake.addDestinationsToRouteMutex.Unlock()
	fake.AddDestinationsToRouteStub = nil
	fake.addDestinationsToRouteReturns = struct {
		result1 repositories.RouteRecord
		result2 error
	}{result1, result2}
}

func (fake *CFRouteRepository) AddDestinationsToRouteReturnsOnCall(i int, result1 repositories.RouteRecord, result2 error) {
	fake.addDestinationsToRouteMutex.Lock()
	defer fake.addDestinationsToRouteMutex.Unlock()
	fake.AddDestinationsToRouteStub = nil
	if fake.addDestinationsToRouteReturnsOnCall == nil {
		fake.addDestinationsToRouteReturnsOnCall = make(map[int]struct {
			result1 repositories.RouteRecord
			result2 error
		})
	}
	fake.addDestinationsToRouteReturnsOnCall[i] = struct {
		result1 repositories.RouteRecord
		result2 error
	}{result1, result2}
}

func (fake *CFRouteRepository) GetOrCreateRoute(arg1 context.Context, arg2 authorization.Info, arg3 repositories.CreateRouteMessage) (repositories.RouteRecord, error) {
	fake.getOrCreateRouteMutex.Lock()
	ret, specificReturn := fake.getOrCreateRouteReturnsOnCall[len(fake.getOrCreateRouteArgsForCall)]
	fake.getOrCreateRouteArgsForCall = append(fake.getOrCreateRouteArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateRouteMessage
	}{arg1, arg2, arg3})
	stub := fake.GetOrCreateRouteStub
	fakeReturns := fake.getOrCreateRouteReturns
	fake.recordInvocation("GetOrCreateRoute", []interface{}{arg1, arg2, arg3})
	fake.getOrCreateRouteMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFRouteRepository) GetOrCreateRouteCallCount() int {
	fake.getOrCreateRouteMutex.RLock()
	defer fake.getOrCreateRouteMutex.RUnlock()
	return len(fake.getOrCreateRouteArgsForCall)
}

func (fake *CFRouteRepository) GetOrCreateRouteCalls(stub func(context.Context, authorization.Info, repositories.CreateRouteMessage) (repositories.RouteRecord, error)) {
	fake.getOrCreateRouteMutex.Lock()
	defer fake.getOrCreateRouteMutex.Unlock()
	fake.GetOrCreateRouteStub = stub
}

func (fake *CFRouteRepository) GetOrCreateRouteArgsForCall(i int) (context.Context, authorization.Info, repositories.CreateRouteMessage) {
	fake.getOrCreateRouteMutex.RLock()
	defer fake.getOrCreateRouteMutex.RUnlock()
	argsForCall := fake.getOrCreateRouteArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFRouteRepository) GetOrCreateRouteReturns(result1 repositories.RouteRecord, result2 error) {
	fake.getOrCreateRouteMutex.Lock()
	defer fake.getOrCreateRouteMutex.Unlock()
	fake.GetOrCreateRouteStub = nil
	fake.getOrCreateRouteReturns = struct {
		result1 repositories.RouteRecord
		result2 error
	}{result1, result2}
}

func (fake *CFRouteRepository) GetOrCreateRouteReturnsOnCall(i int, result1 repositories.RouteRecord, result2 error) {
	fake.getOrCreateRouteMutex.Lock()
	defer fake.getOrCreateRouteMutex.Unlock()
	fake.GetOrCreateRouteStub = nil
	if fake.getOrCreateRouteReturnsOnCall == nil {
		fake.getOrCreateRouteReturnsOnCall = make(map[int]struct {
			result1 repositories.RouteRecord
			result2 error
		})
	}
	fake.getOrCreateRouteReturnsOnCall[i] = struct {
		result1 repositories.RouteRecord
		result2 error
	}{result1, result2}
}

func (fake *CFRouteRepository) ListRoutesForApp(arg1 context.Context, arg2 authorization.Info, arg3 string, arg4 string) ([]repositories.RouteRecord, error) {
	fake.listRoutesForAppMutex.Lock()
	ret, specificReturn := fake.listRoutesForAppReturnsOnCall[len(fake.listRoutesForAppArgsForCall)]
	fake.listRoutesForAppArgsForCall = append(fake.listRoutesForAppArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
		arg4 string
	}{arg1, arg2, arg3, arg4})
	stub := fake.ListRoutesForAppStub
	fakeReturns := fake.listRoutesForAppReturns
	fake.recordInvocation("ListRoutesForApp", []interface{}{arg1, arg2, arg3, arg4})
	fake.listRoutesForAppMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFRouteRepository) ListRoutesForAppCallCount() int {
	fake.listRoutesForAppMutex.RLock()
	defer fake.listRoutesForAppMutex.RUnlock()
	return len(fake.listRoutesForAppArgsForCall)
}

func (fake *CFRouteRepository) ListRoutesForAppCalls(stub func(context.Context, authorization.Info, string, string) ([]repositories.RouteRecord, error)) {
	fake.listRoutesForAppMutex.Lock()
	defer fake.listRoutesForAppMutex.Unlock()
	fake.ListRoutesForAppStub = stub
}

func (fake *CFRouteRepository) ListRoutesForAppArgsForCall(i int) (context.Context, authorization.Info, string, string) {
	fake.listRoutesForAppMutex.RLock()
	defer fake.listRoutesForAppMutex.RUnlock()
	argsForCall := fake.listRoutesForAppArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *CFRouteRepository) ListRoutesForAppReturns(result1 []repositories.RouteRecord, result2 error) {
	fake.listRoutesForAppMutex.Lock()
	defer fake.listRoutesForAppMutex.Unlock()
	fake.ListRoutesForAppStub = nil
	fake.listRoutesForAppReturns = struct {
		result1 []repositories.RouteRecord
		result2 error
	}{result1, result2}
}

func (fake *CFRouteRepository) ListRoutesForAppReturnsOnCall(i int, result1 []repositories.RouteRecord, result2 error) {
	fake.listRoutesForAppMutex.Lock()
	defer fake.listRoutesForAppMutex.Unlock()
	fake.ListRoutesForAppStub = nil
	if fake.listRoutesForAppReturnsOnCall == nil {
		fake.listRoutesForAppReturnsOnCall = make(map[int]struct {
			result1 []repositories.RouteRecord
			result2 error
		})
	}
	fake.listRoutesForAppReturnsOnCall[i] = struct {
		result1 []repositories.RouteRecord
		result2 error
	}{result1, result2}
}

func (fake *CFRouteRepository) RemoveDestinationFromRoute(arg1 context.Context, arg2 authorization.Info, arg3 repositories.RemoveDestinationFromRouteMessage) (repositories.RouteRecord, error) {
	fake.removeDestinationFromRouteMutex.Lock()
	ret, specificReturn := fake.removeDestinationFromRouteReturnsOnCall[len(fake.removeDestinationFromRouteArgsForCall)]
	fake.removeDestinationFromRouteArgsForCall = append(fake.removeDestinationFromRouteArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.RemoveDestinationFromRouteMessage
	}{arg1, arg2, arg3})
	stub := fake.RemoveDestinationFromRouteStub
	fakeReturns := fake.removeDestinationFromRouteReturns
	fake.recordInvocation("RemoveDestinationFromRoute", []interface{}{arg1, arg2, arg3})
	fake.removeDestinationFromRouteMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFRouteRepository) RemoveDestinationFromRouteCallCount() int {
	fake.removeDestinationFromRouteMutex.RLock()
	defer fake.removeDestinationFromRouteMutex.RUnlock()
	return len(fake.removeDestinationFromRouteArgsForCall)
}

func (fake *CFRouteRepository) RemoveDestinationFromRouteCalls(stub func(context.Context, authorization.Info, repositories.RemoveDestinationFromRouteMessage) (repositories.RouteRecord, error)) {
	fake.removeDestinationFromRouteMutex.Lock()
	defer fake.removeDestinationFromRouteMutex.Unlock()
	fake.RemoveDestinationFromRouteStub = stub
}

func (fake *CFRouteRepository) RemoveDestinationFromRouteArgsForCall(i int) (context.Context, authorization.Info, repositories.RemoveDestinationFromRouteMessage) {
	fake.removeDestinationFromRouteMutex.RLock()
	defer fake.removeDestinationFromRouteMutex.RUnlock()
	argsForCall := fake.removeDestinationFromRouteArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFRouteRepository) RemoveDestinationFromRouteReturns(result1 repositories.RouteRecord, result2 error) {
	fake.removeDestinationFromRouteMutex.Lock()
	defer fake.removeDestinationFromRouteMutex.Unlock()
	fake.RemoveDestinationFromRouteStub = nil
	fake.removeDestinationFromRouteReturns = struct {
		result1 repositories.RouteRecord
		result2 error
	}{result1, result2}
}

func (fake *CFRouteRepository) RemoveDestinationFromRouteReturnsOnCall(i int, result1 repositories.RouteRecord, result2 error) {
	fake.removeDestinationFromRouteMutex.Lock()
	defer fake.removeDestinationFromRouteMutex.Unlock()
	fake.RemoveDestinationFromRouteStub = nil
	if fake.removeDestinationFromRouteReturnsOnCall == nil {
		fake.removeDestinationFromRouteReturnsOnCall = make(map[int]struct {
			result1 repositories.RouteRecord
			result2 error
		})
	}
	fake.removeDestinationFromRouteReturnsOnCall[i] = struct {
		result1 repositories.RouteRecord
		result2 error
	}{result1, result2}
}

func (fake *CFRouteRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addDestinationsToRouteMutex.RLock()
	defer fake.addDestinationsToRouteMutex.RUnlock()
	fake.getOrCreateRouteMutex.RLock()
	defer fake.getOrCreateRouteMutex.RUnlock()
	fake.listRoutesForAppMutex.RLock()
	defer fake.listRoutesForAppMutex.RUnlock()
	fake.removeDestinationFromRouteMutex.RLock()
	defer fake.removeDestinationFromRouteMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *CFRouteRepository) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ shared.CFRouteRepository = new(CFRouteRepository)
