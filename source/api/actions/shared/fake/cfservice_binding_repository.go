// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"context"
	"sync"

	"code.cloudfoundry.org/korifi/api/actions/shared"
	"code.cloudfoundry.org/korifi/api/authorization"
	"code.cloudfoundry.org/korifi/api/repositories"
)

type CFServiceBindingRepository struct {
	CreateServiceBindingStub        func(context.Context, authorization.Info, repositories.CreateServiceBindingMessage) (repositories.ServiceBindingRecord, error)
	createServiceBindingMutex       sync.RWMutex
	createServiceBindingArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateServiceBindingMessage
	}
	createServiceBindingReturns struct {
		result1 repositories.ServiceBindingRecord
		result2 error
	}
	createServiceBindingReturnsOnCall map[int]struct {
		result1 repositories.ServiceBindingRecord
		result2 error
	}
	DeleteServiceBindingStub        func(context.Context, authorization.Info, string) error
	deleteServiceBindingMutex       sync.RWMutex
	deleteServiceBindingArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}
	deleteServiceBindingReturns struct {
		result1 error
	}
	deleteServiceBindingReturnsOnCall map[int]struct {
		result1 error
	}
	ListServiceBindingsStub        func(context.Context, authorization.Info, repositories.ListServiceBindingsMessage) ([]repositories.ServiceBindingRecord, error)
	listServiceBindingsMutex       sync.RWMutex
	listServiceBindingsArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.ListServiceBindingsMessage
	}
	listServiceBindingsReturns struct {
		result1 []repositories.ServiceBindingRecord
		result2 error
	}
	listServiceBindingsReturnsOnCall map[int]struct {
		result1 []repositories.ServiceBindingRecord
		result2 error
	}
	UpdateServiceBindingStub        func(context.Context, authorization.Info, repositories.UpdateServiceBindingMessage) (repositories.ServiceBindingRecord, error)
	updateServiceBindingMutex       sync.RWMutex
	updateServiceBindingArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.UpdateServiceBindingMessage
	}
	updateServiceBindingReturns struct {
		result1 repositories.ServiceBindingRecord
		result2 error
	}
	updateServiceBindingReturnsOnCall map[int]struct {
		result1 repositories.ServiceBindingRecord
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *CFServiceBindingRepository) CreateServiceBinding(arg1 context.Context, arg2 authorization.Info, arg3 repositories.CreateServiceBindingMessage) (repositories.ServiceBindingRecord, error) {
	fake.createServiceBindingMutex.Lock()
	ret, specificReturn := fake.createServiceBindingReturnsOnCall[len(fake.createServiceBindingArgsForCall)]
	fake.createServiceBindingArgsForCall = append(fake.createServiceBindingArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateServiceBindingMessage
	}{arg1, arg2, arg3})
	stub := fake.CreateServiceBindingStub
	fakeReturns := fake.createServiceBindingReturns
	fake.recordInvocation("CreateServiceBinding", []interface{}{arg1, arg2, arg3})
	fake.createServiceBindingMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFServiceBindingRepository) CreateServiceBindingCallCount() int {
	fake.createServiceBindingMutex.RLock()
	defer fake.createServiceBindingMutex.RUnlock()
	return len(fake.createServiceBindingArgsForCall)
}

func (fake *CFServiceBindingRepository) CreateServiceBindingCalls(stub func(context.Context, authorization.Info, repositories.CreateServiceBindingMessage) (repositories.ServiceBindingRecord, error)) {
	fake.createServiceBindingMutex.Lock()
	defer fake.createServiceBindingMutex.Unlock()
	fake.CreateServiceBindingStub = stub
}

func (fake *CFServiceBindingRepository) CreateServiceBindingArgsForCall(i int) (context.Context, authorization.Info, repositories.CreateServiceBindingMessage) {
	fake.createServiceBindingMutex.RLock()
	defer fake.createServiceBindingMutex.RUnlock()
	argsForCall := fake.createServiceBindingArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFServiceBindingRepository) CreateServiceBindingReturns(result1 repositories.ServiceBindingRecord, result2 error) {
	fake.createServiceBindingMutex.Lock()
	defer fake.createServiceBindingMutex.Unlock()
	fake.CreateServiceBindingStub = nil
	fake.createServiceBindingReturns = struct {
		result1 repositories.ServiceBindingRecord
		result2 error
	}{result1, result2}
}

func (fake *CFServiceBindingRepository) CreateServiceBindingReturnsOnCall(i int, result1 repositories.ServiceBindingRecord, result2 error) {
	fake.createServiceBindingMutex.Lock()
	defer fake.createServiceBindingMutex.Unlock()
	fake.CreateServiceBindingStub = nil
	if fake.createServiceBindingReturnsOnCall == nil {
		fake.createServiceBindingReturnsOnCall = make(map[int]struct {
			result1 repositories.ServiceBindingRecord
			result2 error
		})
	}
	fake.createServiceBindingReturnsOnCall[i] = struct {
		result1 repositories.ServiceBindingRecord
		result2 error
	}{result1, result2}
}

func (fake *CFServiceBindingRepository) DeleteServiceBinding(arg1 context.Context, arg2 authorization.Info, arg3 string) error {
	fake.deleteServiceBindingMutex.Lock()
	ret, specificReturn := fake.deleteServiceBindingReturnsOnCall[len(fake.deleteServiceBindingArgsForCall)]
	fake.deleteServiceBindingArgsForCall = append(fake.deleteServiceBindingArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.DeleteServiceBindingStub
	fakeReturns := fake.deleteServiceBindingReturns
	fake.recordInvocation("DeleteServiceBinding", []interface{}{arg1, arg2, arg3})
	fake.deleteServiceBindingMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CFServiceBindingRepository) DeleteServiceBindingCallCount() int {
	fake.deleteServiceBindingMutex.RLock()
	defer fake.deleteServiceBindingMutex.RUnlock()
	return len(fake.deleteServiceBindingArgsForCall)
}

func (fake *CFServiceBindingRepository) DeleteServiceBindingCalls(stub func(context.Context, authorization.Info, string) error) {
	fake.deleteServiceBindingMutex.Lock()
	defer fake.deleteServiceBindingMutex.Unlock()
	fake.DeleteServiceBindingStub = stub
}

func (fake *CFServiceBindingRepository) DeleteServiceBindingArgsForCall(i int) (context.Context, authorization.Info, string) {
	fake.deleteServiceBindingMutex.RLock()
	defer fake.deleteServiceBindingMutex.RUnlock()
	argsForCall := fake.deleteServiceBindingArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFServiceBindingRepository) DeleteServiceBindingReturns(result1 error) {
	fake.deleteServiceBindingMutex.Lock()
	defer fake.deleteServiceBindingMutex.Unlock()
	fake.DeleteServiceBindingStub = nil
	fake.deleteServiceBindingReturns = struct {
		result1 error
	}{result1}
}

func (fake *CFServiceBindingRepository) DeleteServiceBindingReturnsOnCall(i int, result1 error) {
	fake.deleteServiceBindingMutex.Lock()
	defer fake.deleteServiceBindingMutex.Unlock()
	fake.DeleteServiceBindingStub = nil
	if fake.deleteServiceBindingReturnsOnCall == nil {
		fake.deleteServiceBindingReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteServiceBindingReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *CFServiceBindingRepository) ListServiceBindings(arg1 context.Context, arg2 authorization.Info, arg3 repositories.ListServiceBindingsMessage) ([]repositories.ServiceBindingRecord, error) {
	fake.listServiceBindingsMutex.Lock()
	ret, specificReturn := fake.listServiceBindingsReturnsOnCall[len(fake.listServiceBindingsArgsForCall)]
	fake.listServiceBindingsArgsForCall = append(fake.listServiceBindingsArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.ListServiceBindingsMessage
	}{arg1, arg2, arg3})
	stub := fake.ListServiceBindingsStub
	fakeReturns := fake.listServiceBindingsReturns
	fake.recordInvocation("ListServiceBindings", []interface{}{arg1, arg2, arg3})
	fake.listServiceBindingsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFServiceBindingRepository) ListServiceBindingsCallCount() int {
	fake.listServiceBindingsMutex.RLock()
	defer fake.listServiceBindingsMutex.RUnlock()
	return len(fake.listServiceBindingsArgsForCall)
}

func (fake *CFServiceBindingRepository) ListServiceBindingsCalls(stub func(context.Context, authorization.Info, repositories.ListServiceBindingsMessage) ([]repositories.ServiceBindingRecord, error)) {
	fake.listServiceBindingsMutex.Lock()
	defer fake.listServiceBindingsMutex.Unlock()
	fake.ListServiceBindingsStub = stub
}

func (fake *CFServiceBindingRepository) ListServiceBindingsArgsForCall(i int) (context.Context, authorization.Info, repositories.ListServiceBindingsMessage) {
	fake.listServiceBindingsMutex.RLock()
	defer fake.listServiceBindingsMutex.RUnlock()
	argsForCall := fake.listServiceBindingsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFServiceBindingRepository) ListServiceBindingsReturns(result1 []repositories.ServiceBindingRecord, result2 error) {
	fake.listServiceBindingsMutex.Lock()
	defer fake.listServiceBindingsMutex.Unlock()
	fake.ListServiceBindingsStub = nil
	fake.listServiceBindingsReturns = struct {
		result1 []repositories.ServiceBindingRecord
		result2 error
	}{result1, result2}
}

func (fake *CFServiceBindingRepository) ListServiceBindingsReturnsOnCall(i int, result1 []repositories.ServiceBindingRecord, result2 error) {
	fake.listServiceBindingsMutex.Lock()
	defer fake.listServiceBindingsMutex.Unlock()
	fake.ListServiceBindingsStub = nil
	if fake.listServiceBindingsReturnsOnCall == nil {
		fake.listServiceBindingsReturnsOnCall = make(map[int]struct {
			result1 []repositories.ServiceBindingRecord
			result2 error
		})
	}
	fake.listServiceBindingsReturnsOnCall[i] = struct {
		result1 []repositories.ServiceBindingRecord
		result2 error
	}{result1, result2}
}

func (fake *CFServiceBindingRepository) UpdateServiceBinding(arg1 context.Context, arg2 authorization.Info, arg3 repositories.UpdateServiceBindingMessage) (repositories.ServiceBindingRecord, error) {
	fake.updateServiceBindingMutex.Lock()
	ret, specificReturn := fake.updateServiceBindingReturnsOnCall[len(fake.updateServiceBindingArgsForCall)]
	fake.updateServiceBindingArgsForCall = append(fake.updateServiceBindingArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.UpdateServiceBindingMessage
	}{arg1, arg2, arg3})
	stub := fake.UpdateServiceBindingStub
	fakeReturns := fake.updateServiceBindingReturns
	fake.recordInvocation("UpdateServiceBinding", []interface{}{arg1, arg2, arg3})
	fake.updateServiceBindingMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFServiceBindingRepository) UpdateServiceBindingCallCount() int {
	fake.updateServiceBindingMutex.RLock()
	defer fake.updateServiceBindingMutex.RUnlock()
	return len(fake.updateServiceBindingArgsForCall)
}

func (fake *CFServiceBindingRepository) UpdateServiceBindingCalls(stub func(context.Context, authorization.Info, repositories.UpdateServiceBindingMessage) (repositories.ServiceBindingRecord, error)) {
	fake.updateServiceBindingMutex.Lock()
	defer fake.updateServiceBindingMutex.Unlock()
	fake.UpdateServiceBindingStub = stub
}

func (fake *CFServiceBindingRepository) UpdateServiceBindingArgsForCall(i int) (context.Context, authorization.Info, repositories.UpdateServiceBindingMessage) {
	fake.updateServiceBindingMutex.RLock()
	defer fake.updateServiceBindingMutex.RUnlock()
	argsForCall := fake.updateServiceBindingArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFServiceBindingRepository) UpdateServiceBindingReturns(result1 repositories.ServiceBindingRecord, result2 error) {
	fake.updateServiceBindingMutex.Lock()
	defer fake.updateServiceBindingMutex.Unlock()
	fake.UpdateServiceBindingStub = nil
	fake.updateServiceBindingReturns = struct {
		result1 repositories.ServiceBindingRecord
		result2 error
	}{result1, result2}
}

func (fake *CFServiceBindingRepository) UpdateServiceBindingReturnsOnCall(i int, result1 repositories.ServiceBindingRecord, result2 error) {
	fake.updateServiceBindingMutex.Lock()
	defer fake.updateServiceBindingMutex.Unlock()
	fake.UpdateServiceBindingStub = nil
	if fake.updateServiceBindingReturnsOnCall == nil {
		fake.updateServiceBindingReturnsOnCall = make(map[int]struct {
			result1 repositories.ServiceBindingRecord
			result2 error
		})
	}
	fake.updateServiceBindingReturnsOnCall[i] = struct {
		result1 repositories.ServiceBindingRecord
		result2 error
	}{result1, result2}
}

func (fake *CFServiceBindingRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createServiceBindingMutex.RLock()
	defer fake.createServiceBindingMutex.RUnlock()
	fake.deleteServiceBindingMutex.RLock()
	defer fake.deleteServiceBindingMutex.RUnlock()
	fake.listServiceBindingsMutex.RLock()
	defer fake.listServiceBindingsMutex.RUnlock()
	fake.updateServiceBindingMutex.RLock()
	defer fake.updateServiceBindingMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *CFServiceBindingRepository) recordInvocation(key string, args []interface{}) {
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

var _ shared.CFServiceBindingRepository = new(CFServiceBindingRepository)
