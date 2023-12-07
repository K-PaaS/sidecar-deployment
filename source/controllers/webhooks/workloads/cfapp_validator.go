package workloads

import (
	"context"
	"fmt"

	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/controllers/webhooks"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	AppEntityType        = "app"
	AppDecodingErrorType = "AppDecodingError"
)

var cfapplog = logf.Log.WithName("cfapp-validate")

//+kubebuilder:webhook:path=/validate-korifi-cloudfoundry-org-v1alpha1-cfapp,mutating=false,failurePolicy=fail,sideEffects=NoneOnDryRun,groups=korifi.cloudfoundry.org,resources=cfapps,verbs=create;update;delete,versions=v1alpha1,name=vcfapp.korifi.cloudfoundry.org,admissionReviewVersions={v1,v1beta1}

type CFAppValidator struct {
	duplicateValidator webhooks.NameValidator
}

var _ webhook.CustomValidator = &CFAppValidator{}

func NewCFAppValidator(duplicateValidator webhooks.NameValidator) *CFAppValidator {
	return &CFAppValidator{
		duplicateValidator: duplicateValidator,
	}
}

func (v *CFAppValidator) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&korifiv1alpha1.CFApp{}).
		WithValidator(v).
		Complete()
}

func (v *CFAppValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	app, ok := obj.(*korifiv1alpha1.CFApp)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFApp but got a %T", obj))
	}

	return nil, v.duplicateValidator.ValidateCreate(ctx, cfapplog, app.Namespace, app)
}

func (v *CFAppValidator) ValidateUpdate(ctx context.Context, oldObj, obj runtime.Object) (admission.Warnings, error) {
	app, ok := obj.(*korifiv1alpha1.CFApp)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFApp but got a %T", obj))
	}

	if !app.GetDeletionTimestamp().IsZero() {
		return nil, nil
	}

	oldApp, ok := oldObj.(*korifiv1alpha1.CFApp)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFApp but got a %T", oldObj))
	}

	if app.Spec.Lifecycle.Type != oldApp.Spec.Lifecycle.Type {
		return nil, webhooks.ValidationError{
			Type:    "ImmutableFieldError",
			Message: fmt.Sprintf("Lifecycle type cannot be changed from %s to %s", oldApp.Spec.Lifecycle.Type, app.Spec.Lifecycle.Type),
		}.ExportJSONError()
	}

	return nil, v.duplicateValidator.ValidateUpdate(ctx, cfapplog, app.Namespace, oldApp, app)
}

func (v *CFAppValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	app, ok := obj.(*korifiv1alpha1.CFApp)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFApp but got a %T", obj))
	}

	return nil, v.duplicateValidator.ValidateDelete(ctx, cfapplog, app.Namespace, app)
}
