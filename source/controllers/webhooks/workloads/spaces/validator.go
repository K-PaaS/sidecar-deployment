package spaces

import (
	"context"
	"errors"
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
	CFSpaceEntityType = "cfspace"
)

var spaceLogger = logf.Log.WithName("cfspace-validate")

//+kubebuilder:webhook:path=/validate-korifi-cloudfoundry-org-v1alpha1-cfspace,mutating=false,failurePolicy=fail,sideEffects=NoneOnDryRun,groups=korifi.cloudfoundry.org,resources=cfspaces,verbs=create;update;delete,versions=v1alpha1,name=vcfspace.korifi.cloudfoundry.org,admissionReviewVersions={v1,v1beta1}

type Validator struct {
	duplicateValidator webhooks.NameValidator
	placementValidator webhooks.NamespaceValidator
}

var _ webhook.CustomValidator = &Validator{}

func NewValidator(duplicateSpaceValidator webhooks.NameValidator, placementValidator webhooks.NamespaceValidator) *Validator {
	return &Validator{
		duplicateValidator: duplicateSpaceValidator,
		placementValidator: placementValidator,
	}
}

func (v *Validator) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&korifiv1alpha1.CFSpace{}).
		WithValidator(v).
		Complete()
}

func (v *Validator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	space, ok := obj.(*korifiv1alpha1.CFSpace)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFSpace but got a %T", obj))
	}

	if len(space.Name) > webhooks.MaxLabelLength {
		return nil, errors.New("space name cannot be longer than 63 chars")
	}

	err := v.duplicateValidator.ValidateCreate(ctx, spaceLogger, space.Namespace, space)
	if err != nil {
		return nil, err
	}

	return nil, v.placementValidator.ValidateSpaceCreate(*space)
}

func (v *Validator) ValidateUpdate(ctx context.Context, oldObj, obj runtime.Object) (admission.Warnings, error) {
	space, ok := obj.(*korifiv1alpha1.CFSpace)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFSpace but got a %T", obj))
	}

	if !space.GetDeletionTimestamp().IsZero() {
		return nil, nil
	}

	oldSpace, ok := oldObj.(*korifiv1alpha1.CFSpace)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFSpace but got a %T", obj))
	}

	return nil, v.duplicateValidator.ValidateUpdate(ctx, spaceLogger, oldSpace.Namespace, oldSpace, space)
}

func (v *Validator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	space, ok := obj.(*korifiv1alpha1.CFSpace)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFSpace but got a %T", obj))
	}

	return nil, v.duplicateValidator.ValidateDelete(ctx, spaceLogger, space.Namespace, space)
}
