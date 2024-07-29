package orgs

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
	CFOrgEntityType      = "cforg"
	OrgDecodingErrorType = "OrgDecodingError"
)

var cfOrgLog = logf.Log.WithName("cforg-validate")

//+kubebuilder:webhook:path=/validate-korifi-cloudfoundry-org-v1alpha1-cforg,mutating=false,failurePolicy=fail,sideEffects=NoneOnDryRun,groups=korifi.cloudfoundry.org,resources=cforgs,verbs=create;update;delete,versions=v1alpha1,name=vcforg.korifi.cloudfoundry.org,admissionReviewVersions={v1,v1beta1}

type Validator struct {
	duplicateValidator webhooks.NameValidator
	placementValidator webhooks.NamespaceValidator
}

var _ webhook.CustomValidator = &Validator{}

func NewValidator(duplicateValidator webhooks.NameValidator, placementValidator webhooks.NamespaceValidator) *Validator {
	return &Validator{
		duplicateValidator: duplicateValidator,
		placementValidator: placementValidator,
	}
}

func (v *Validator) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&korifiv1alpha1.CFOrg{}).
		WithValidator(v).
		Complete()
}

func (v *Validator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	org, ok := obj.(*korifiv1alpha1.CFOrg)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFOrg but got a %T", obj))
	}

	if len(org.Name) > webhooks.MaxLabelLength {
		return nil, errors.New("org name cannot be longer than 63 chars")
	}

	err := v.placementValidator.ValidateOrgCreate(*org)
	if err != nil {
		cfOrgLog.Info(err.Error())
		return nil, err
	}

	return nil, v.duplicateValidator.ValidateCreate(ctx, cfOrgLog, org.Namespace, org)
}

func (v *Validator) ValidateUpdate(ctx context.Context, oldObj, obj runtime.Object) (admission.Warnings, error) {
	org, ok := obj.(*korifiv1alpha1.CFOrg)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFOrg but got a %T", obj))
	}

	if !org.GetDeletionTimestamp().IsZero() {
		return nil, nil
	}

	oldOrg, ok := oldObj.(*korifiv1alpha1.CFOrg)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFOrg but got a %T", obj))
	}

	return nil, v.duplicateValidator.ValidateUpdate(ctx, cfOrgLog, org.Namespace, oldOrg, org)
}

func (v *Validator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	org, ok := obj.(*korifiv1alpha1.CFOrg)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFOrg but got a %T", obj))
	}

	return nil, v.duplicateValidator.ValidateDelete(ctx, cfOrgLog, org.Namespace, org)
}
