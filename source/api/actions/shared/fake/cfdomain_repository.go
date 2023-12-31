// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"context"
	"sync"

	"code.cloudfoundry.org/korifi/api/actions/shared"
	"code.cloudfoundry.org/korifi/api/authorization"
	"code.cloudfoundry.org/korifi/api/repositories"
)

type CFDomainRepository struct {
	GetDomainByNameStub        func(context.Context, authorization.Info, string) (repositories.DomainRecord, error)
	getDomainByNameMutex       sync.RWMutex
	getDomainByNameArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}
	getDomainByNameReturns struct {
		result1 repositories.DomainRecord
		result2 error
	}
	getDomainByNameReturnsOnCall map[int]struct {
		result1 repositories.DomainRecord
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *CFDomainRepository) GetDomainByName(arg1 context.Context, arg2 authorization.Info, arg3 string) (repositories.DomainRecord, error) {
	fake.getDomainByNameMutex.Lock()
	ret, specificReturn := fake.getDomainByNameReturnsOnCall[len(fake.getDomainByNameArgsForCall)]
	fake.getDomainByNameArgsForCall = append(fake.getDomainByNameArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.GetDomainByNameStub
	fakeReturns := fake.getDomainByNameReturns
	fake.recordInvocation("GetDomainByName", []interface{}{arg1, arg2, arg3})
	fake.getDomainByNameMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFDomainRepository) GetDomainByNameCallCount() int {
	fake.getDomainByNameMutex.RLock()
	defer fake.getDomainByNameMutex.RUnlock()
	return len(fake.getDomainByNameArgsForCall)
}

func (fake *CFDomainRepository) GetDomainByNameCalls(stub func(context.Context, authorization.Info, string) (repositories.DomainRecord, error)) {
	fake.getDomainByNameMutex.Lock()
	defer fake.getDomainByNameMutex.Unlock()
	fake.GetDomainByNameStub = stub
}

func (fake *CFDomainRepository) GetDomainByNameArgsForCall(i int) (context.Context, authorization.Info, string) {
	fake.getDomainByNameMutex.RLock()
	defer fake.getDomainByNameMutex.RUnlock()
	argsForCall := fake.getDomainByNameArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFDomainRepository) GetDomainByNameReturns(result1 repositories.DomainRecord, result2 error) {
	fake.getDomainByNameMutex.Lock()
	defer fake.getDomainByNameMutex.Unlock()
	fake.GetDomainByNameStub = nil
	fake.getDomainByNameReturns = struct {
		result1 repositories.DomainRecord
		result2 error
	}{result1, result2}
}

func (fake *CFDomainRepository) GetDomainByNameReturnsOnCall(i int, result1 repositories.DomainRecord, result2 error) {
	fake.getDomainByNameMutex.Lock()
	defer fake.getDomainByNameMutex.Unlock()
	fake.GetDomainByNameStub = nil
	if fake.getDomainByNameReturnsOnCall == nil {
		fake.getDomainByNameReturnsOnCall = make(map[int]struct {
			result1 repositories.DomainRecord
			result2 error
		})
	}
	fake.getDomainByNameReturnsOnCall[i] = struct {
		result1 repositories.DomainRecord
		result2 error
	}{result1, result2}
}

func (fake *CFDomainRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getDomainByNameMutex.RLock()
	defer fake.getDomainByNameMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *CFDomainRepository) recordInvocation(key string, args []interface{}) {
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

var _ shared.CFDomainRepository = new(CFDomainRepository)
