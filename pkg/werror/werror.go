package werror

import (
	"encoding/json"
	"errors"

	"google.golang.org/grpc/codes"
)

type AppError struct {
	Code    codes.Code     `json:"code"`
	Message string         `json:"message"`
	Details AppErrorDetail `json:"details"`
}

type AppErrorDetail struct {
	Err       error  `json:"error"`
	ErrDomain string `json:"domain"`
	ErrReason string `json:"reason"`
	RequestId string `json:"request_id"`
}

// Error Returns Message if Details.Err is nil.
func (e *AppError) Error() string {
	if err := e.Details.Err; err != nil {
		return err.Error()
	}

	return e.Message
}

func (e *AppError) Msg() string {
	return e.Message
}

func (e *AppError) ToJson() string {
	b, _ := json.Marshal(e)

	return string(b)
}

func (e *AppError) SetCode(code codes.Code) *AppError {
	e.Code = code

	return e
}

func NewInternalError(appError AppError) *AppError {
	return appError.SetCode(codes.Internal)
}

func NewUnauthenticatedError(appError AppError) *AppError {
	return appError.SetCode(codes.Unauthenticated)
}

func NewForbiddenError(appError AppError) *AppError {
	return appError.SetCode(codes.PermissionDenied)
}

func NewNotFoundError(appError AppError) *AppError {
	return appError.SetCode(codes.NotFound)
}

func NewBadRequestError(appError AppError) *AppError {
	return appError.SetCode(codes.Internal)
}

var (
	ErrRequestMissingMetadata = errors.New("request: missing metadata")
)
