package qcloud

import "fmt"

const (
	// internal sdk error code
	internal_sdk_error_code_base = -9000
	PARAM_ERROR = internal_sdk_error_code_base + iota
)

type Error struct {
	Code int
	CodeDesc string
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s(%d): %s", e.CodeDesc, e.Code, e.Message)
}

func NewError(code int, codeDesc , message string) error {
	return &Error{
		Code: code,
		CodeDesc: codeDesc,
		Message: message,
	}
}

func NewParamError(msg string) error {
	return NewError(PARAM_ERROR, "ParamError", msg)
}
