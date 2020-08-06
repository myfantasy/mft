package mft

import (
	"strconv"
)

// Error type with internal
type Error struct {
	Code          int    `json:"code,omitempty"`
	Msg           string `json:"msg,omitempty"`
	InternalError string `json:"ie,omitempty"`
}

// ErrorCommonCode - no code error
const ErrorCommonCode int = 50000

// Error implement error interface
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	cd := ""
	if e.Code != 0 {
		cd = "[" + strconv.Itoa(e.Code) + "] "
	}

	if e.InternalError == "" {
		return cd + e.Msg
	}
	return cd + e.Msg + "\n" + e.InternalError
}

// ErrorCSE make Error from string with internal error
func ErrorCSE(code int, err string, internalError error) *Error {
	return &Error{
		Msg:           err,
		Code:          code,
		InternalError: internalError.Error(),
	}
}

// ErrorCS make Error from string
func ErrorCS(code int, err string) *Error {
	return &Error{
		Msg:  err,
		Code: code,
	}
}

// ErrorS make Error from string
func ErrorS(err string) *Error {
	return &Error{
		Msg:  err,
		Code: ErrorCommonCode,
	}
}

// ErrorCE make Error from any error
func ErrorCE(code int, err error) *Error {
	return &Error{
		Msg:  err.Error(),
		Code: code,
	}
}

// AppendS append next error level saving code
func (e *Error) AppendS(errs string) *Error {
	return &Error{
		Msg:           errs,
		InternalError: e.Error(),
		Code:          e.Code,
	}
}

// AppendE append next error level saving code
func (e *Error) AppendE(errs error) *Error {
	return &Error{
		Msg:           errs.Error(),
		InternalError: e.Error(),
		Code:          e.Code,
	}
}

// ErrorE make Error from any error
func ErrorE(err error) *Error {
	return &Error{
		Msg:  err.Error(),
		Code: ErrorCommonCode,
	}
}

// ErrorNew - Create new Error from msg and another error
func ErrorNew(msg string, internalError error) *Error {
	return &Error{
		Msg:           msg,
		InternalError: internalError.Error(),
		Code:          ErrorCommonCode,
	}
}

// ErrorNew2 - Create new Error
func ErrorNew2(msg string, internalError error, internal2Error error) *Error {
	return ErrorNew(msg, ErrorNew(internalError.Error(), internal2Error))
}
