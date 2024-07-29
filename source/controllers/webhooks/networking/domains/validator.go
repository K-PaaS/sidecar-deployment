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

package domains

import (
	"context"
	"fmt"
	"strings"

	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	validationwebhook "code.cloudfoundry.org/korifi/controllers/webhooks/validation"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	DomainDecodingErrorType  = "DomainDecodingError"
	DuplicateDomainErrorType = "DuplicateDomainError"
	InvalidDomainErrorType   = "InvalidDomainError"
)

// log is for logging in this package.
var log = logf.Log.WithName("domain-validation")

//+kubebuilder:webhook:path=/validate-korifi-cloudfoundry-org-v1alpha1-cfdomain,mutating=false,failurePolicy=fail,sideEffects=None,groups=korifi.cloudfoundry.org,resources=cfdomains,verbs=create;update,versions=v1alpha1,name=vcfdomain.korifi.cloudfoundry.org,admissionReviewVersions=v1

type Validator struct {
	client client.Client
}

var _ webhook.CustomValidator = &Validator{}

func NewValidator(client client.Client) *Validator {
	return &Validator{
		client: client,
	}
}

func (v *Validator) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&korifiv1alpha1.CFDomain{}).
		WithValidator(v).
		Complete()
}

func (v *Validator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	domain, ok := obj.(*korifiv1alpha1.CFDomain)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFDomain but got a %T", obj))
	}

	err := validateDomainName(domain.Spec.Name)
	if err != nil {
		return nil, validationwebhook.ValidationError{
			Type:    InvalidDomainErrorType,
			Message: fmt.Sprintf("%q is not a valid domain: %s", domain.Spec.Name, err.Error()),
		}.ExportJSONError()
	}

	isOverlapping, err := v.domainIsOverlapping(ctx, domain.Spec.Name)
	if err != nil {
		log.Info("error checking for overlapping domain", "reason", err)
		return nil, validationwebhook.ValidationError{
			Type:    validationwebhook.UnknownErrorType,
			Message: validationwebhook.UnknownErrorMessage,
		}.ExportJSONError()
	}

	if isOverlapping {
		return nil, validationwebhook.ValidationError{
			Type:    DuplicateDomainErrorType,
			Message: "Overlapping domain exists",
		}.ExportJSONError()
	}

	return nil, nil
}

func validateDomainName(domainName string) error {
	return validation.IsFullyQualifiedDomainName(field.NewPath("CFDomain", "Spec", "Name"), domainName).ToAggregate()
}

func (v *Validator) ValidateUpdate(ctx context.Context, oldObj runtime.Object, obj runtime.Object) (admission.Warnings, error) {
	domain, ok := obj.(*korifiv1alpha1.CFDomain)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFDomain but got a %T", obj))
	}

	if !domain.GetDeletionTimestamp().IsZero() {
		return nil, nil
	}

	oldDomain, ok := oldObj.(*korifiv1alpha1.CFDomain)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a CFDomain but got a %T", obj))
	}

	if oldDomain.Spec.Name != domain.Spec.Name {
		return nil, validationwebhook.ValidationError{
			Type:    validationwebhook.ImmutableFieldErrorType,
			Message: fmt.Sprintf(validationwebhook.ImmutableFieldErrorMessageTemplate, "CFDomain.Spec.Name"),
		}.ExportJSONError()
	}

	return nil, nil
}

func (v *Validator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func (v *Validator) domainIsOverlapping(ctx context.Context, domainName string) (bool, error) {
	var existingDomainList korifiv1alpha1.CFDomainList
	err := v.client.List(ctx, &existingDomainList)
	if err != nil {
		return true, err
	}

	domainElements := strings.Split(domainName, ".")

	for _, existingDomain := range existingDomainList.Items {
		existingDomainElements := strings.Split(existingDomain.Spec.Name, ".")
		if isSubDomain(domainElements, existingDomainElements) {
			return true, nil
		}
	}

	return false, nil
}

func isSubDomain(domainElements, existingDomainElements []string) bool {
	var shorterSlice, longerSlice *[]string

	if len(domainElements) < len(existingDomainElements) {
		shorterSlice = &domainElements
		longerSlice = &existingDomainElements
	} else {
		shorterSlice = &existingDomainElements
		longerSlice = &domainElements
	}

	offset := len(*longerSlice) - len(*shorterSlice)

	for i := len(*shorterSlice) - 1; i >= 0; i-- {
		if (*shorterSlice)[i] != (*longerSlice)[i+offset] {
			return false
		}
	}

	return true
}
