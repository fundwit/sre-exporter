package fail

import (
	"errors"
	"net/http"
)

var ErrUnexpected = errors.New("common.internal_server_error")
var ErrInvalidArguments = errors.New("invalid arguments")

var ErrUnauthenticated = errors.New("unauthenticated")
var ErrForbidden = errors.New("forbidden")

var ErrNotFound = errors.New("not found")
var ErrNoContent = errors.New("no content")
var ErrInvalidPassword = errors.New("invalid password")

type BizError interface {
	Respond() *BizErrorDetail
}

type BizErrorDetail struct {
	Status  int
	Code    string
	Message string

	Data interface{}
}

type ErrBadParam struct {
	Param        string
	InvalidValue string

	Cause error
}

func (e *ErrBadParam) Unwrap() error {
	return e.Cause
}

func (e *ErrBadParam) Error() string {
	message := "bad param"
	if e.Param != "" {
		message = "invalid " + e.Param + " '" + e.InvalidValue + "'"
	} else if e.Cause != nil {
		message = e.Cause.Error()
	}
	return message
}

func (e *ErrBadParam) Respond() *BizErrorDetail {
	return &BizErrorDetail{
		Status:  http.StatusBadRequest,
		Code:    "common.bad_param",
		Message: e.Error(),
		Data:    nil,
	}
}
