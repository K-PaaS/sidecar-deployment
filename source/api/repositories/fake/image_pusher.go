// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"context"
	"io"
	"sync"

	"code.cloudfoundry.org/korifi/api/repositories"
	"code.cloudfoundry.org/korifi/tools/image"
)

type ImagePusher struct {
	PushStub        func(context.Context, image.Creds, string, io.Reader, ...string) (string, error)
	pushMutex       sync.RWMutex
	pushArgsForCall []struct {
		arg1 context.Context
		arg2 image.Creds
		arg3 string
		arg4 io.Reader
		arg5 []string
	}
	pushReturns struct {
		result1 string
		result2 error
	}
	pushReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ImagePusher) Push(arg1 context.Context, arg2 image.Creds, arg3 string, arg4 io.Reader, arg5 ...string) (string, error) {
	fake.pushMutex.Lock()
	ret, specificReturn := fake.pushReturnsOnCall[len(fake.pushArgsForCall)]
	fake.pushArgsForCall = append(fake.pushArgsForCall, struct {
		arg1 context.Context
		arg2 image.Creds
		arg3 string
		arg4 io.Reader
		arg5 []string
	}{arg1, arg2, arg3, arg4, arg5})
	stub := fake.PushStub
	fakeReturns := fake.pushReturns
	fake.recordInvocation("Push", []interface{}{arg1, arg2, arg3, arg4, arg5})
	fake.pushMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ImagePusher) PushCallCount() int {
	fake.pushMutex.RLock()
	defer fake.pushMutex.RUnlock()
	return len(fake.pushArgsForCall)
}

func (fake *ImagePusher) PushCalls(stub func(context.Context, image.Creds, string, io.Reader, ...string) (string, error)) {
	fake.pushMutex.Lock()
	defer fake.pushMutex.Unlock()
	fake.PushStub = stub
}

func (fake *ImagePusher) PushArgsForCall(i int) (context.Context, image.Creds, string, io.Reader, []string) {
	fake.pushMutex.RLock()
	defer fake.pushMutex.RUnlock()
	argsForCall := fake.pushArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *ImagePusher) PushReturns(result1 string, result2 error) {
	fake.pushMutex.Lock()
	defer fake.pushMutex.Unlock()
	fake.PushStub = nil
	fake.pushReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *ImagePusher) PushReturnsOnCall(i int, result1 string, result2 error) {
	fake.pushMutex.Lock()
	defer fake.pushMutex.Unlock()
	fake.PushStub = nil
	if fake.pushReturnsOnCall == nil {
		fake.pushReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.pushReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *ImagePusher) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.pushMutex.RLock()
	defer fake.pushMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ImagePusher) recordInvocation(key string, args []interface{}) {
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

var _ repositories.ImagePusher = new(ImagePusher)