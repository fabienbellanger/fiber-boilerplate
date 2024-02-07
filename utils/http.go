package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Error status codes
const (
	StatusBadRequest                   = 400
	StatusUnauthorized                 = 401
	StatusPaymentRequired              = 402
	StatusForbidden                    = 403
	StatusNotFound                     = 404
	StatusMethodNotAllowed             = 405
	StatusNotAcceptable                = 406
	StatusProxyAuthRequired            = 407
	StatusRequestTimeout               = 408
	StatusConflict                     = 409
	StatusGone                         = 410
	StatusLengthRequired               = 411
	StatusPreconditionFailed           = 412
	StatusRequestEntityTooLarge        = 413
	StatusRequestURITooLong            = 414
	StatusUnsupportedMediaType         = 415
	StatusRequestedRangeNotSatisfiable = 416
	StatusExpectationFailed            = 417
	StatusTeapot                       = 418
	StatusMisdirectedRequest           = 421
	StatusUnprocessableEntity          = 422
	StatusLocked                       = 423
	StatusFailedDependency             = 424
	StatusTooEarly                     = 425
	StatusUpgradeRequired              = 426
	StatusPreconditionRequired         = 428
	StatusTooManyRequests              = 429
	StatusRequestHeaderFieldsTooLarge  = 431
	StatusUnavailableForLegalReasons   = 451

	StatusInternalServerError           = 500
	StatusNotImplemented                = 501
	StatusBadGateway                    = 502
	StatusServiceUnavailable            = 503
	StatusGatewayTimeout                = 504
	StatusHTTPVersionNotSupported       = 505
	StatusVariantAlsoNegotiates         = 506
	StatusInsufficientStorage           = 507
	StatusLoopDetected                  = 508
	StatusNotExtended                   = 510
	StatusNetworkAuthenticationRequired = 511
)

// HTTPError represents an HTTP error.
type HTTPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Err     error       `json:"-"`
}

// NewHTTPError returns a new HTTPError.
func NewHTTPError(code int, message string, details interface{}, err error) *HTTPError {

	return &HTTPError{
		Code:    code,
		Message: message,
		Details: details,
		Err:     err,
	}
}

func (e HTTPError) Error() string {
	return e.Message
}

// NewError returns a fiber error and log the error.
func NewError(c *fiber.Ctx, logger *zap.Logger, msg, details string, err error) *fiber.Error {
	logger.Error(details, zap.String("description", msg), zap.Error(err), zap.String("requestId", fmt.Sprintf("%v", c.Locals("requestid"))))

	return fiber.NewError(fiber.StatusInternalServerError, msg)
}

// PaginateResponse represents a response with pagination.
// TODO: Move to domain?
type PaginateResponse struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
}
