package bindings

import (
	"context"
	"fmt"

	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/controllers/webhooks"
	validation "code.cloudfoundry.org/korifi/controllers/webhooks/validation"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	ServiceBindingEntityType = "servicebinding"
	ServiceBindingErrorType  = "ServiceBindingValidationError"
)

// log is for logging in this package.
var cfservicebindinglog = logf.Log.WithName("cfservicebinding-validator")

//+kubebuilder:webhook:path=/validate-korifi-cloudfoundry-org-v1alpha1-cfservicebinding,mutating=false,failurePolicy=fail,sideEffects=NoneOnDryRun,groups=korifi.cloudfoundry.org,resources=cfservicebindings,verbs=create;update;delete,versions=v1alpha1,name=vcfservicebinding.korifi.cloudfoundry.org,admissionReviewVersions={v1,v1beta1}

func (v *CFServiceBindingValidator) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&korifiv1alpha1.CFServiceBinding{}).
		WithValidator(v).
		Complete()
}

type CFServiceBindingValidator struct {
	duplicateValidator webhooks.NameValidator
}

var _ webhook.CustomValidator = &CFServiceBindingValidator{}

func NewCFServiceBindingValidator(duplicateValidator webhooks.NameValidator) *CFServiceBindingValidator {
	return &CFServiceBindingValidator{
		duplicateValidator: duplicateValidator,
	}
}

func (v *CFServiceBindingValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	serviceBinding, ok := obj.(*korifiv1alpha1.CFServiceBinding)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFServiceBinding but got a %T", obj))
	}

	return nil, v.duplicateValidator.ValidateCreate(ctx, cfservicebindinglog, serviceBinding.Namespace, serviceBinding)
}

func (v *CFServiceBindingValidator) ValidateUpdate(ctx context.Context, oldObj, obj runtime.Object) (admission.Warnings, error) {
	serviceBinding, ok := obj.(*korifiv1alpha1.CFServiceBinding)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFServiceBinding but got a %T", obj))
	}

	if !serviceBinding.GetDeletionTimestamp().IsZero() {
		return nil, nil
	}

	oldServiceBinding, ok := oldObj.(*korifiv1alpha1.CFServiceBinding)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFServiceBinding but got a %T", oldObj))
	}

	if oldServiceBinding.Spec.AppRef.Name != serviceBinding.Spec.AppRef.Name {
		return nil, validation.ValidationError{Type: ServiceBindingErrorType, Message: "AppRef.Name is immutable"}
	}

	if oldServiceBinding.Spec.Service.Name != serviceBinding.Spec.Service.Name {
		return nil, validation.ValidationError{Type: ServiceBindingErrorType, Message: "Service.Name is immutable"}
	}

	if oldServiceBinding.Spec.Service.Namespace != serviceBinding.Spec.Service.Namespace {
		return nil, validation.ValidationError{Type: ServiceBindingErrorType, Message: "Service.Namespace is immutable"}
	}

	return nil, nil
}

func (v *CFServiceBindingValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	serviceBinding, ok := obj.(*korifiv1alpha1.CFServiceBinding)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFServiceBinding but got a %T", obj))
	}

	return nil, v.duplicateValidator.ValidateDelete(ctx, cfservicebindinglog, serviceBinding.Namespace, serviceBinding)
}
