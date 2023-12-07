// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"context"
	"sync"

	"code.cloudfoundry.org/korifi/api/authorization"
	"code.cloudfoundry.org/korifi/api/handlers"
	"code.cloudfoundry.org/korifi/api/repositories"
)

type RunnerInfoRepository struct {
	GetRunnerInfoStub        func(context.Context, authorization.Info, string) (repositories.RunnerInfoRecord, error)
	getRunnerInfoMutex       sync.RWMutex
	getRunnerInfoArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}
	getRunnerInfoReturns struct {
		result1 repositories.RunnerInfoRecord
		result2 error
	}
	getRunnerInfoReturnsOnCall map[int]struct {
		result1 repositories.RunnerInfoRecord
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *RunnerInfoRepository) GetRunnerInfo(arg1 context.Context, arg2 authorization.Info, arg3 string) (repositories.RunnerInfoRecord, error) {
	fake.getRunnerInfoMutex.Lock()
	ret, specificReturn := fake.getRunnerInfoReturnsOnCall[len(fake.getRunnerInfoArgsForCall)]
	fake.getRunnerInfoArgsForCall = append(fake.getRunnerInfoArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.GetRunnerInfoStub
	fakeReturns := fake.getRunnerInfoReturns
	fake.recordInvocation("GetRunnerInfo", []interface{}{arg1, arg2, arg3})
	fake.getRunnerInfoMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *RunnerInfoRepository) GetRunnerInfoCallCount() int {
	fake.getRunnerInfoMutex.RLock()
	defer fake.getRunnerInfoMutex.RUnlock()
	return len(fake.getRunnerInfoArgsForCall)
}

func (fake *RunnerInfoRepository) GetRunnerInfoCalls(stub func(context.Context, authorization.Info, string) (repositories.RunnerInfoRecord, error)) {
	fake.getRunnerInfoMutex.Lock()
	defer fake.getRunnerInfoMutex.Unlock()
	fake.GetRunnerInfoStub = stub
}

func (fake *RunnerInfoRepository) GetRunnerInfoArgsForCall(i int) (context.Context, authorization.Info, string) {
	fake.getRunnerInfoMutex.RLock()
	defer fake.getRunnerInfoMutex.RUnlock()
	argsForCall := fake.getRunnerInfoArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *RunnerInfoRepository) GetRunnerInfoReturns(result1 repositories.RunnerInfoRecord, result2 error) {
	fake.getRunnerInfoMutex.Lock()
	defer fake.getRunnerInfoMutex.Unlock()
	fake.GetRunnerInfoStub = nil
	fake.getRunnerInfoReturns = struct {
		result1 repositories.RunnerInfoRecord
		result2 error
	}{result1, result2}
}

func (fake *RunnerInfoRepository) GetRunnerInfoReturnsOnCall(i int, result1 repositories.RunnerInfoRecord, result2 error) {
	fake.getRunnerInfoMutex.Lock()
	defer fake.getRunnerInfoMutex.Unlock()
	fake.GetRunnerInfoStub = nil
	if fake.getRunnerInfoReturnsOnCall == nil {
		fake.getRunnerInfoReturnsOnCall = make(map[int]struct {
			result1 repositories.RunnerInfoRecord
			result2 error
		})
	}
	fake.getRunnerInfoReturnsOnCall[i] = struct {
		result1 repositories.RunnerInfoRecord
		result2 error
	}{result1, result2}
}

func (fake *RunnerInfoRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getRunnerInfoMutex.RLock()
	defer fake.getRunnerInfoMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *RunnerInfoRepository) recordInvocation(key string, args []interface{}) {
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

var _ handlers.RunnerInfoRepository = new(RunnerInfoRepository)
