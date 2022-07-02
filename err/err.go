package err

import (
	"net/http"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
)

type ErrType string

const (
	ErrBadRequest   ErrType = "BadRequest"
	ErrUnauthorized ErrType = "Unauthorized"
	ErrInternal     ErrType = "Internal"
)

const (
	defaultErrType     ErrType    = ErrInternal
	grpcDefaultErrCode codes.Code = codes.Internal
	httpDefaultErrCode int        = http.StatusInternalServerError
)

var grpcErrCode = map[ErrType]codes.Code{
	ErrBadRequest:   codes.InvalidArgument,
	ErrUnauthorized: codes.Unauthenticated,
	ErrInternal:     codes.Internal,
}

var httpErrCode = map[ErrType]int{
	ErrBadRequest:   http.StatusBadRequest,
	ErrUnauthorized: http.StatusUnauthorized,
	ErrInternal:     http.StatusInternalServerError,
}

type (
	// error - abstraction for error object
	Error interface {
		Error() string
		GRPCStatusCode() codes.Code
		HttpStatusCode() int
	}
)

type (
	err struct {
		Err      error
		GRPCCode codes.Code
		HTTPCode int
	}
)

// NewError - function for initializing error
func NewError(e error, errType ErrType) Error {
	return &err{
		Err:      e,
		GRPCCode: setGRPCStatusCode(errType),
		HTTPCode: setHttpStatusCode(errType),
	}
}

// GetError - get error with Error type
func GetError(e error) Error {
	errc := errors.Cause(e)
	switch errc.(type) {
	case *err:
		return e.(*err)
	default:
		return &err{
			Err:      e,
			GRPCCode: setGRPCStatusCode(defaultErrType),
			HTTPCode: setHttpStatusCode(defaultErrType),
		}
	}
}

// Error - function for return error message
func (e *err) Error() string {
	if e == nil {
		return ""
	}
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

// Error - function for return error message
func (e *err) GRPCStatusCode() codes.Code {
	if e == nil {
		return 0
	}
	return e.GRPCCode
}

// Error - function for return error message
func (e *err) HttpStatusCode() int {
	if e == nil {
		return 0
	}
	return e.HTTPCode
}

func setGRPCStatusCode(errType ErrType) codes.Code {
	if errType != "" {
		return grpcErrCode[errType]
	} else {
		return grpcDefaultErrCode
	}
}

func setHttpStatusCode(errType ErrType) int {
	if errType != "" {
		return httpErrCode[errType]
	} else {
		return httpDefaultErrCode
	}
}
