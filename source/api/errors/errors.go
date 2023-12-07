package errors

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"code.cloudfoundry.org/korifi/controllers/webhooks"

	"github.com/go-logr/logr"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ApiError interface {
	Detail() string
	Title() string
	Code() int
	HttpStatus() int
	Unwrap() error
	Error() string
}

// LogAndReturn logs api errors at the info level and other errors at the
// error level since api errors are expected recoverable conditions.
// It returns the error for convenience.
func LogAndReturn(logger logr.Logger, err error, msg string, keysAndValues ...interface{}) error {
	var apiError ApiError
	if errors.As(err, &apiError) {
		keysAndValues = append(keysAndValues, "reason", err)
		logger.Info(msg, keysAndValues...)
	} else {
		logger.Error(err, msg, keysAndValues...)
	}

	return err
}

type apiError struct {
	cause             error
	detail            string
	title             string
	code              int
	httpStatus        int
	additionalDetails map[string]string
}

func (e apiError) Error() string {
	if e.cause == nil {
		return "unknown"
	}

	return e.cause.Error()
}

func (e apiError) Unwrap() error {
	return e.cause
}

func (e apiError) Detail() string {
	detail := e.detail
	for k, v := range e.additionalDetails {
		detail += fmt.Sprintf(" %s=%q", k, v)
	}
	return detail
}

func (e apiError) Title() string {
	return e.title
}

func (e apiError) Code() int {
	return e.code
}

func (e apiError) HttpStatus() int {
	return e.httpStatus
}

func toKeyValues(s ...string) map[string]string {
	result := map[string]string{}

	for i := 0; i < len(s); i += 2 {
		key := s[i]
		val := ""
		if i+1 < len(s) {
			val = s[i+1]
		}
		result[key] = val
	}

	return result
}

type UnprocessableEntityError struct {
	apiError
}

func NewUnprocessableEntityError(cause error, detail string) UnprocessableEntityError {
	return UnprocessableEntityError{
		apiError{
			cause:      cause,
			title:      "CF-UnprocessableEntity",
			detail:     detail,
			code:       10008,
			httpStatus: http.StatusUnprocessableEntity,
		},
	}
}

type MessageParseError struct {
	apiError
}

func NewMessageParseError(cause error) MessageParseError {
	return MessageParseError{
		apiError{
			cause:      cause,
			title:      "CF-MessageParseError",
			detail:     "Request invalid due to parse error: invalid request body",
			code:       1001,
			httpStatus: http.StatusBadRequest,
		},
	}
}

// UnknownError is a generic wrapper over an error Korifi cannot recover from.
// Unknown errors should be only used by the presentation layer to present such
// an error to the user. Other components (handlers, repositories, etc.) should
// simply return the incoming error, it would be mapped to `UnknownError` by
// the presentation layer
type UnknownError struct {
	apiError
}

// NewUnknownError creates an UnknownError. One should generally not create
// unknown errors as generic errors are automatically presented as unknown
// errors to the user
func NewUnknownError(cause error) UnknownError {
	return UnknownError{
		apiError{
			cause:      cause,
			title:      "UnknownError",
			detail:     "An unknown error occurred.",
			code:       10001,
			httpStatus: http.StatusInternalServerError,
		},
	}
}

type NotFoundError struct {
	apiError
}

func NewNotFoundError(cause error, resourceType string, additionalDetails ...string) NotFoundError {
	return NotFoundError{
		apiError{
			cause:             cause,
			title:             "CF-ResourceNotFound",
			detail:            fmt.Sprintf("%s not found. Ensure it exists and you have access to it.", resourceType),
			additionalDetails: toKeyValues(additionalDetails...),
			code:              10010,
			httpStatus:        http.StatusNotFound,
		},
	}
}

type EndpointNotFoundError struct {
	apiError
}

func NewEndpointNotFoundError() EndpointNotFoundError {
	return EndpointNotFoundError{
		apiError{
			title:      "CF-NotFound",
			detail:     "Unknown request",
			code:       10000,
			httpStatus: http.StatusNotFound,
		},
	}
}

type InvalidAuthError struct {
	apiError
}

func NewInvalidAuthError(cause error) InvalidAuthError {
	return InvalidAuthError{
		apiError{
			cause:      cause,
			title:      "CF-InvalidAuthToken",
			detail:     "Invalid Auth Token",
			code:       1000,
			httpStatus: http.StatusUnauthorized,
		},
	}
}

type NotAuthenticatedError struct {
	apiError
}

func NewNotAuthenticatedError(cause error) NotAuthenticatedError {
	return NotAuthenticatedError{
		apiError{
			cause:      cause,
			title:      "CF-NotAuthenticated",
			detail:     "Authentication error",
			code:       10002,
			httpStatus: http.StatusUnauthorized,
		},
	}
}

type ForbiddenError struct {
	apiError
	resourceType string
}

func (e ForbiddenError) ResourceType() string {
	return e.resourceType
}

func NewForbiddenError(cause error, resourceType string) ForbiddenError {
	return ForbiddenError{
		apiError: apiError{
			cause:      cause,
			title:      "CF-NotAuthorized",
			detail:     "You are not authorized to perform the requested action",
			code:       10003,
			httpStatus: http.StatusForbidden,
		},
		resourceType: resourceType,
	}
}

type BadQueryParamValueError struct {
	apiError
}

func NewBadQueryParamValueError(key string, validValues ...string) BadQueryParamValueError {
	return BadQueryParamValueError{
		apiError: apiError{
			title:      "CF-BadQueryParameter",
			detail:     fmt.Sprintf("The query parameter is invalid: %s can only be: %s", key, quotedCommaSeparatedList(validValues)),
			code:       10005,
			httpStatus: http.StatusBadRequest,
		},
	}
}

type UnknownKeyError struct {
	apiError
}

func NewUnknownKeyError(cause error, validKeys []string) UnknownKeyError {
	return UnknownKeyError{
		apiError: apiError{
			cause:      cause,
			title:      "CF-BadQueryParameter",
			detail:     fmt.Sprintf("The query parameter is invalid: Valid parameters are: %s", quotedCommaSeparatedList(validKeys)),
			code:       10005,
			httpStatus: http.StatusBadRequest,
		},
	}
}

func quotedCommaSeparatedList(in []string) string {
	var out []string
	for _, i := range in {
		out = append(out, fmt.Sprintf("'%s'", i))
	}
	return strings.Join(out, ", ")
}

type UniquenessError struct {
	apiError
}

func NewUniquenessError(cause error, detail string) UniquenessError {
	return UniquenessError{
		apiError: apiError{
			cause:      cause,
			title:      "CF-UniquenessError",
			detail:     detail,
			code:       10016,
			httpStatus: http.StatusUnprocessableEntity,
		},
	}
}

type InvalidRequestError struct {
	apiError
}

func NewInvalidRequestError(cause error, detail string) InvalidRequestError {
	return InvalidRequestError{
		apiError: apiError{
			cause:      cause,
			title:      "CF-InvalidRequest",
			detail:     detail,
			code:       10004,
			httpStatus: http.StatusBadRequest,
		},
	}
}

type PackageBitsAlreadyUploadedError struct {
	apiError
}

func NewPackageBitsAlreadyUploadedError(cause error) PackageBitsAlreadyUploadedError {
	return PackageBitsAlreadyUploadedError{
		apiError: apiError{
			cause:      cause,
			title:      "CF-PackageBitsAlreadyUploaded",
			detail:     "Bits may be uploaded only once. Create a new package to upload different bits.",
			code:       150004,
			httpStatus: http.StatusBadRequest,
		},
	}
}

type BlobstoreUnavailableError struct {
	apiError
}

func NewBlobstoreUnavailableError(cause error) BlobstoreUnavailableError {
	return BlobstoreUnavailableError{
		apiError: apiError{
			cause:      cause,
			title:      "CF-BlobstoreUnavailable",
			detail:     "Error uploading source package to the container registry",
			code:       150006,
			httpStatus: http.StatusBadGateway,
		},
	}
}

type ResourceNotReadyError struct {
	apiError
}

func NewResourceNotReadyError(cause error) ResourceNotReadyError {
	return ResourceNotReadyError{
		apiError: apiError{
			cause:      cause,
			title:      "CF-ResourceNotReady",
			detail:     cause.Error(),
			code:       420000,
			httpStatus: http.StatusInternalServerError,
		},
	}
}

type RollingDeployNotSupportedError struct {
	apiError
}

func NewRollingDeployNotSupportedError(runnerName string) RollingDeployNotSupportedError {
	detail := fmt.Sprintf("The configured runner '%s' does not support rolling deploys", runnerName)
	return RollingDeployNotSupportedError{
		apiError: apiError{
			title:      "CF-RollingDeployNotSupported",
			detail:     detail,
			code:       42000,
			httpStatus: http.StatusBadRequest,
		},
	}
}

func FromK8sError(err error, resourceType string) error {
	if webhookValidationError, ok := webhooks.WebhookErrorToValidationError(err); ok {
		return NewUnprocessableEntityError(err, webhookValidationError.GetMessage())
	}

	switch {
	case k8serrors.IsUnauthorized(err):
		return NewInvalidAuthError(err)
	case k8serrors.IsNotFound(err):
		return NewNotFoundError(err, resourceType)
	case k8serrors.IsForbidden(err):
		return NewForbiddenError(err, resourceType)
	case k8serrors.IsInvalid(err):
		cause, ok := k8serrors.StatusCause(err, metav1.CauseTypeFieldValueInvalid)
		if ok {
			return NewUnprocessableEntityError(err, fmt.Sprintf("%s is invalid: %s", cause.Field, cause.Message))
		}
		return NewUnprocessableEntityError(err, resourceType)
	default:
		return err
	}
}

func AsUnprocessableEntity(err error, detail string, errTypes ...ApiError) error {
	if err == nil {
		return nil
	}

	for i := range errTypes {
		// At this point in time the errors in the errType array are downgraded
		// to `ApiError`. This means that pointers to api errors that only
		// embed `apiError` are assignable to each other. Therefore `errors.As`
		// would return `true` and would change the initial value type of the
		// array entry. That is why we need to get the "desiredType" first and
		// compare it to the type that has been set by `errors.As`
		desiredErrType := reflect.ValueOf(errTypes[i]).Type()

		if !errors.As(err, &errTypes[i]) {
			continue
		}

		asErrType := reflect.ValueOf(errTypes[i]).Type()
		if asErrType != desiredErrType {
			continue
		}

		return NewUnprocessableEntityError(errTypes[i].Unwrap(), detail)
	}

	return err
}

func ForbiddenAsNotFound(err error) error {
	var forbiddenErr ForbiddenError
	if errors.As(err, &forbiddenErr) {
		return NewNotFoundError(forbiddenErr.Unwrap(), forbiddenErr.ResourceType())
	}
	return err
}

// DropletForbiddenAsNotFound is a special case due to the CF CLI expecting the error message "Droplet not found" exactly instead of the generic case
// https://github.com/cloudfoundry/korifi/issues/965
func DropletForbiddenAsNotFound(err error) error {
	var forbiddenErr ForbiddenError
	if errors.As(err, &forbiddenErr) {
		return NotFoundError{
			apiError{
				cause:      forbiddenErr.Unwrap(),
				title:      "CF-ResourceNotFound",
				detail:     "Droplet not found",
				code:       10010,
				httpStatus: http.StatusNotFound,
			},
		}
	}
	var notFoundErr NotFoundError
	if errors.As(err, &notFoundErr) {
		return NotFoundError{
			apiError{
				cause:      notFoundErr.Unwrap(),
				title:      "CF-ResourceNotFound",
				detail:     "Droplet not found",
				code:       10010,
				httpStatus: http.StatusNotFound,
			},
		}
	}
	return err
}
