// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"context"
	"sync"

	"code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/controllers/controllers/workloads/build"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type BuildReconciler struct {
	ReconcileBuildStub        func(context.Context, *v1alpha1.CFBuild, *v1alpha1.CFApp, *v1alpha1.CFPackage) (reconcile.Result, error)
	reconcileBuildMutex       sync.RWMutex
	reconcileBuildArgsForCall []struct {
		arg1 context.Context
		arg2 *v1alpha1.CFBuild
		arg3 *v1alpha1.CFApp
		arg4 *v1alpha1.CFPackage
	}
	reconcileBuildReturns struct {
		result1 reconcile.Result
		result2 error
	}
	reconcileBuildReturnsOnCall map[int]struct {
		result1 reconcile.Result
		result2 error
	}
	SetupWithManagerStub        func(manager.Manager) *builder.Builder
	setupWithManagerMutex       sync.RWMutex
	setupWithManagerArgsForCall []struct {
		arg1 manager.Manager
	}
	setupWithManagerReturns struct {
		result1 *builder.Builder
	}
	setupWithManagerReturnsOnCall map[int]struct {
		result1 *builder.Builder
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *BuildReconciler) ReconcileBuild(arg1 context.Context, arg2 *v1alpha1.CFBuild, arg3 *v1alpha1.CFApp, arg4 *v1alpha1.CFPackage) (reconcile.Result, error) {
	fake.reconcileBuildMutex.Lock()
	ret, specificReturn := fake.reconcileBuildReturnsOnCall[len(fake.reconcileBuildArgsForCall)]
	fake.reconcileBuildArgsForCall = append(fake.reconcileBuildArgsForCall, struct {
		arg1 context.Context
		arg2 *v1alpha1.CFBuild
		arg3 *v1alpha1.CFApp
		arg4 *v1alpha1.CFPackage
	}{arg1, arg2, arg3, arg4})
	stub := fake.ReconcileBuildStub
	fakeReturns := fake.reconcileBuildReturns
	fake.recordInvocation("ReconcileBuild", []interface{}{arg1, arg2, arg3, arg4})
	fake.reconcileBuildMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *BuildReconciler) ReconcileBuildCallCount() int {
	fake.reconcileBuildMutex.RLock()
	defer fake.reconcileBuildMutex.RUnlock()
	return len(fake.reconcileBuildArgsForCall)
}

func (fake *BuildReconciler) ReconcileBuildCalls(stub func(context.Context, *v1alpha1.CFBuild, *v1alpha1.CFApp, *v1alpha1.CFPackage) (reconcile.Result, error)) {
	fake.reconcileBuildMutex.Lock()
	defer fake.reconcileBuildMutex.Unlock()
	fake.ReconcileBuildStub = stub
}

func (fake *BuildReconciler) ReconcileBuildArgsForCall(i int) (context.Context, *v1alpha1.CFBuild, *v1alpha1.CFApp, *v1alpha1.CFPackage) {
	fake.reconcileBuildMutex.RLock()
	defer fake.reconcileBuildMutex.RUnlock()
	argsForCall := fake.reconcileBuildArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *BuildReconciler) ReconcileBuildReturns(result1 reconcile.Result, result2 error) {
	fake.reconcileBuildMutex.Lock()
	defer fake.reconcileBuildMutex.Unlock()
	fake.ReconcileBuildStub = nil
	fake.reconcileBuildReturns = struct {
		result1 reconcile.Result
		result2 error
	}{result1, result2}
}

func (fake *BuildReconciler) ReconcileBuildReturnsOnCall(i int, result1 reconcile.Result, result2 error) {
	fake.reconcileBuildMutex.Lock()
	defer fake.reconcileBuildMutex.Unlock()
	fake.ReconcileBuildStub = nil
	if fake.reconcileBuildReturnsOnCall == nil {
		fake.reconcileBuildReturnsOnCall = make(map[int]struct {
			result1 reconcile.Result
			result2 error
		})
	}
	fake.reconcileBuildReturnsOnCall[i] = struct {
		result1 reconcile.Result
		result2 error
	}{result1, result2}
}

func (fake *BuildReconciler) SetupWithManager(arg1 manager.Manager) *builder.Builder {
	fake.setupWithManagerMutex.Lock()
	ret, specificReturn := fake.setupWithManagerReturnsOnCall[len(fake.setupWithManagerArgsForCall)]
	fake.setupWithManagerArgsForCall = append(fake.setupWithManagerArgsForCall, struct {
		arg1 manager.Manager
	}{arg1})
	stub := fake.SetupWithManagerStub
	fakeReturns := fake.setupWithManagerReturns
	fake.recordInvocation("SetupWithManager", []interface{}{arg1})
	fake.setupWithManagerMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *BuildReconciler) SetupWithManagerCallCount() int {
	fake.setupWithManagerMutex.RLock()
	defer fake.setupWithManagerMutex.RUnlock()
	return len(fake.setupWithManagerArgsForCall)
}

func (fake *BuildReconciler) SetupWithManagerCalls(stub func(manager.Manager) *builder.Builder) {
	fake.setupWithManagerMutex.Lock()
	defer fake.setupWithManagerMutex.Unlock()
	fake.SetupWithManagerStub = stub
}

func (fake *BuildReconciler) SetupWithManagerArgsForCall(i int) manager.Manager {
	fake.setupWithManagerMutex.RLock()
	defer fake.setupWithManagerMutex.RUnlock()
	argsForCall := fake.setupWithManagerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *BuildReconciler) SetupWithManagerReturns(result1 *builder.Builder) {
	fake.setupWithManagerMutex.Lock()
	defer fake.setupWithManagerMutex.Unlock()
	fake.SetupWithManagerStub = nil
	fake.setupWithManagerReturns = struct {
		result1 *builder.Builder
	}{result1}
}

func (fake *BuildReconciler) SetupWithManagerReturnsOnCall(i int, result1 *builder.Builder) {
	fake.setupWithManagerMutex.Lock()
	defer fake.setupWithManagerMutex.Unlock()
	fake.SetupWithManagerStub = nil
	if fake.setupWithManagerReturnsOnCall == nil {
		fake.setupWithManagerReturnsOnCall = make(map[int]struct {
			result1 *builder.Builder
		})
	}
	fake.setupWithManagerReturnsOnCall[i] = struct {
		result1 *builder.Builder
	}{result1}
}

func (fake *BuildReconciler) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.reconcileBuildMutex.RLock()
	defer fake.reconcileBuildMutex.RUnlock()
	fake.setupWithManagerMutex.RLock()
	defer fake.setupWithManagerMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *BuildReconciler) recordInvocation(key string, args []interface{}) {
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

var _ build.BuildReconciler = new(BuildReconciler)
