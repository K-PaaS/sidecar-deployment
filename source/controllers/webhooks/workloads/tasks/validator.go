/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tasks

import (
	"context"
	"fmt"

	"code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/controllers/webhooks"
	"code.cloudfoundry.org/korifi/controllers/webhooks/validation"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	runtime "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	CancelationNotPossibleErrorType = "CancelationNotPossibleError"
)

// log is for logging in this package.
var cftasklog = logf.Log.WithName("cftask-resource")

//+kubebuilder:webhook:path=/validate-korifi-cloudfoundry-org-v1alpha1-cftask,mutating=false,failurePolicy=fail,sideEffects=None,groups=korifi.cloudfoundry.org,resources=cftasks;cftasks/status,verbs=create;update,versions=v1alpha1,name=vcftask.korifi.cloudfoundry.org,admissionReviewVersions={v1,v1beta1}

type Validator struct{}

var _ webhook.CustomValidator = &Validator{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&v1alpha1.CFTask{}).
		WithValidator(v).
		Complete()
}

var _ webhook.CustomValidator = &Validator{}

func (v *Validator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	task, ok := obj.(*v1alpha1.CFTask)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFTask but got a %T", obj))
	}

	cftasklog.V(1).Info("validate task creation", "namespace", task.Namespace, "name", task.Name)

	if len(task.Spec.Command) == 0 {
		return nil, validation.ValidationError{
			Type:    webhooks.MissingRequredFieldErrorType,
			Message: fmt.Sprintf("task %s:%s is missing required field 'Spec.Command'", task.Namespace, task.Name),
		}.ExportJSONError()
	}

	if task.Spec.AppRef.Name == "" {
		return nil, validation.ValidationError{
			Type:    webhooks.MissingRequredFieldErrorType,
			Message: fmt.Sprintf("task %s:%s is missing required field 'Spec.AppRef.Name'", task.Namespace, task.Name),
		}.ExportJSONError()
	}

	return nil, nil
}

func (v *Validator) ValidateUpdate(ctx context.Context, oldObj runtime.Object, obj runtime.Object) (admission.Warnings, error) {
	newTask, ok := obj.(*v1alpha1.CFTask)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFTask but got a %T", obj))
	}

	if !newTask.GetDeletionTimestamp().IsZero() {
		return nil, nil
	}

	cftasklog.V(1).Info("validate task update", "namespace", newTask.Namespace, "name", newTask.Name)

	oldTask, ok := oldObj.(*v1alpha1.CFTask)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFTask but got a %T", oldObj))
	}

	if newTask.Status.SequenceID < 0 {
		return nil, validation.ValidationError{
			Type:    webhooks.InvalidFieldValueErrorType,
			Message: fmt.Sprintf("task %s:%s Status.SequenceID cannot be negative", newTask.Namespace, newTask.Name),
		}.ExportJSONError()
	}

	if oldTask.Status.SequenceID != 0 && newTask.Status.SequenceID != oldTask.Status.SequenceID {
		return nil, validation.ValidationError{
			Type:    webhooks.ImmutableFieldModificationErrorType,
			Message: fmt.Sprintf("task %s:%s Status.SequenceID is immutable", newTask.Namespace, newTask.Name),
		}.ExportJSONError()
	}

	if oldTask.Spec.Canceled || !newTask.Spec.Canceled {
		return nil, nil
	}

	state := ""
	if meta.IsStatusConditionTrue(newTask.Status.Conditions, v1alpha1.TaskSucceededConditionType) {
		state = "SUCCEEDED"
	} else if meta.IsStatusConditionTrue(newTask.Status.Conditions, v1alpha1.TaskFailedConditionType) {
		state = "FAILED"
	}

	if state != "" {
		return nil, validation.ValidationError{
			Type:    CancelationNotPossibleErrorType,
			Message: fmt.Sprintf("Task state is %s and therefore cannot be canceled", state),
		}.ExportJSONError()
	}

	return nil, nil
}

func (v *Validator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
