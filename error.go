package mft

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
)

// Error type with internal
type Error struct {
	Code              int      `json:"code,omitempty"`
	Msg               string   `json:"msg,omitempty"`
	InternalErrorText string   `json:"iet,omitempty"`
	InternalError     *Error   `json:"ie,omitempty"`
	InternalErrors    []*Error `json:"ies,omitempty"`
	CallStack         string   `json:"call_stack,omitempty"`
}

// ErrorCommonCode - no code error
const ErrorCommonCode int = 50000

func rowPrefixAdd(s string, prefix string) (out string) {
	return prefix + strings.ReplaceAll(s, "\n", "\n"+prefix)
}

var InnerErrorPrefix string = "  "
var FillCallStack bool = false

// Error implement error interface
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	cd := ""
	if e.Code != 0 && e.Code != ErrorCommonCode {
		cd = "[" + strconv.Itoa(e.Code) + "] "
	}

	listMsgs := ""
	if len(e.InternalErrors) > 0 {
		for i, erI := range e.InternalErrors {
			if i > 0 {
				listMsgs += ";\n"
			}
			listMsgs += erI.Error()
		}
		listMsgs = rowPrefixAdd(listMsgs, InnerErrorPrefix)
		listMsgs = "[\n" + listMsgs + "\n]"
		listMsgs = "\n" + rowPrefixAdd(listMsgs, InnerErrorPrefix)
	}

	callStack := ""
	if e.CallStack != "" {
		callStack = "\n" + rowPrefixAdd(e.CallStack, InnerErrorPrefix)
	}

	if e.InternalErrorText == "" && e.InternalError == nil {
		return cd + e.Msg + listMsgs + callStack
	} else if e.InternalErrorText == "" {
		return cd + e.Msg + "\n" + rowPrefixAdd(e.InternalError.Error(), InnerErrorPrefix) + listMsgs + callStack
	} else if e.InternalError == nil {
		return cd + e.Msg + "\n" + rowPrefixAdd(e.InternalErrorText, InnerErrorPrefix) + listMsgs + callStack
	}

	return cd + e.Msg + "\n" + rowPrefixAdd(e.InternalError.Error(), InnerErrorPrefix) +
		"\n" + rowPrefixAdd(e.InternalErrorText, InnerErrorPrefix) + listMsgs + callStack
}

func GetStack() string {
	if !FillCallStack {
		return ""
	}
	return string(debug.Stack())
}

// ErrorCSE make Error from string with internal error
func ErrorCSE(code int, err string, internalError error) *Error {
	er, ok := internalError.(*Error)
	if !ok {
		return &Error{
			Msg:               err,
			Code:              code,
			InternalErrorText: internalError.Error(),
			CallStack:         GetStack(),
		}
	}
	return &Error{
		Msg:           err,
		Code:          code,
		InternalError: er,
		CallStack:     GetStack(),
	}
}

// ErrorCSEf make Error from string with internal error
func ErrorCSEf(code int, internalError error, format string, a ...interface{}) *Error {
	er, ok := internalError.(*Error)
	if !ok {
		return &Error{
			Msg:               fmt.Sprintf(format, a...),
			Code:              code,
			InternalErrorText: internalError.Error(),
			CallStack:         GetStack(),
		}
	}
	return &Error{
		Msg:           fmt.Sprintf(format, a...),
		Code:          code,
		InternalError: er,
		CallStack:     GetStack(),
	}
}

// ErrorCS make Error from string
func ErrorCS(code int, err string) *Error {
	return &Error{
		Msg:       err,
		Code:      code,
		CallStack: GetStack(),
	}
}

// ErrorCSf make Error from string
func ErrorCSf(code int, format string, a ...interface{}) *Error {
	return &Error{
		Msg:       fmt.Sprintf(format, a...),
		Code:      code,
		CallStack: GetStack(),
	}
}

// ErrorS make Error from string
func ErrorS(err string) *Error {
	return &Error{
		Msg:       err,
		Code:      ErrorCommonCode,
		CallStack: GetStack(),
	}
}

// ErrorSf make Error from string
func ErrorSf(format string, a ...interface{}) *Error {
	return &Error{
		Msg:       fmt.Sprintf(format, a...),
		Code:      ErrorCommonCode,
		CallStack: GetStack(),
	}
}

// ErrorCE make Error from any error
func ErrorCE(code int, err error) *Error {
	return &Error{
		Msg:       err.Error(),
		Code:      code,
		CallStack: GetStack(),
	}
}

// AppendS append next error level saving code
func (e *Error) AppendS(errs string) *Error {
	return &Error{
		Msg:           errs,
		InternalError: e,
		Code:          e.Code,
		CallStack:     GetStack(),
	}
}

// AppendSf append next error level saving code
func (e *Error) AppendSf(format string, a ...interface{}) *Error {
	return &Error{
		Msg:           fmt.Sprintf(format, a...),
		InternalError: e,
		Code:          e.Code,
		CallStack:     GetStack(),
	}
}

// AppendE append next error level saving code
func (e *Error) AppendE(errs error) *Error {
	return &Error{
		Msg:           errs.Error(),
		InternalError: e,
		Code:          e.Code,
		CallStack:     GetStack(),
	}
}

// ErrorE make Error from any error
func ErrorE(err error) *Error {
	return &Error{
		Msg:       err.Error(),
		Code:      ErrorCommonCode,
		CallStack: GetStack(),
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
			CallStack:         GetStack(),
		}
	}
	return &Error{
		Msg:           msg,
		InternalError: er,
		Code:          ErrorCommonCode,
		CallStack:     GetStack(),
	}
}

// ErrorNewf - Create new Error from msg and another error
func ErrorNewf(internalError error, format string, a ...interface{}) *Error {
	er, ok := internalError.(*Error)
	if !ok {
		return &Error{
			Msg:               fmt.Sprintf(format, a...),
			InternalErrorText: internalError.Error(),
			Code:              ErrorCommonCode,
			CallStack:         GetStack(),
		}
	}
	return &Error{
		Msg:           fmt.Sprintf(format, a...),
		InternalError: er,
		Code:          ErrorCommonCode,
		CallStack:     GetStack(),
	}
}

// ErrorNew2 - Create new Error
func ErrorNew2(msg string, internalError error, internal2Error error) *Error {
	return ErrorNew(msg, ErrorNew(internalError.Error(), internal2Error))
}

// ErrorNew2f - Create new Error
func ErrorNew2f(internalError error, internal2Error error, format string, a ...interface{}) *Error {
	return ErrorNewf(ErrorNew(internalError.Error(), internal2Error), format, a...)
}

func (e *Error) AppendList(sub ...*Error) (eOut *Error) {
	eOut = e
	for _, ei := range sub {
		if ei == nil {
			continue
		}
		if eOut == nil {
			eOut = &Error{
				CallStack: GetStack(),
			}
		}

		eOut.InternalErrors = append(eOut.InternalErrors, ei)
	}
	return eOut
}
