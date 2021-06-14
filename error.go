package mft

import (
	"strconv"
	"strings"
)

// Error type with internal
type Error struct {
	Code              int    `json:"code,omitempty"`
	Msg               string `json:"msg,omitempty"`
	InternalErrorText string `json:"iet,omitempty"`
	InternalError     *Error `json:"ie,omitempty"`
}

// ErrorCommonCode - no code error
const ErrorCommonCode int = 50000

func rowPrefixAdd(s string, prefix string) (out string) {
	return prefix + strings.ReplaceAll(s, "\n", "\n"+prefix)
}

var InnerErrorPrefix string = "  "

// Error implement error interface
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	cd := ""
	if e.Code != 0 && e.Code != ErrorCommonCode {
		cd = "[" + strconv.Itoa(e.Code) + "] "
	}

	if e.InternalErrorText == "" && e.InternalError == nil {
		return cd + e.Msg
	} else if e.InternalErrorText == "" {
		return cd + e.Msg + "\n" + rowPrefixAdd(e.InternalError.Error(), InnerErrorPrefix)
	} else if e.InternalError == nil {
		return cd + e.Msg + "\n" + rowPrefixAdd(e.InternalErrorText, InnerErrorPrefix)
	}

	return cd + e.Msg + "\n" + rowPrefixAdd(e.InternalError.Error(), InnerErrorPrefix) +
		"\n" + rowPrefixAdd(e.InternalErrorText, InnerErrorPrefix)
}

// ErrorCSE make Error from string with internal error
func ErrorCSE(code int, err string, internalError error) *Error {
	er, ok := internalError.(*Error)
	if !ok {
		return &Error{
			Msg:               err,
			Code:              code,
			InternalErrorText: internalError.Error(),
		}
	}
	return &Error{
		Msg:           err,
		Code:          code,
		InternalError: er,
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
		InternalError: e,
		Code:          e.Code,
	}
}

// AppendE append next error level saving code
func (e *Error) AppendE(errs error) *Error {
	return &Error{
		Msg:           errs.Error(),
		InternalError: e,
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
	er, ok := internalError.(*Error)
	if !ok {
		return &Error{
			Msg:               msg,
			InternalErrorText: internalError.Error(),
			Code:              ErrorCommonCode,
		}
	}
	return &Error{
		Msg:           msg,
		InternalError: er,
		Code:          ErrorCommonCode,
	}
}

// ErrorNew2 - Create new Error
func ErrorNew2(msg string, internalError error, internal2Error error) *Error {
	return ErrorNew(msg, ErrorNew(internalError.Error(), internal2Error))
}
