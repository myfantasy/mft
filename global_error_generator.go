package mft

import (
	"fmt"
	"sync"
)

var ErrorsCodes map[int]string = map[int]string{
	50000: "Undescribe error",
}

var mx sync.RWMutex

func AddErrorsCodes(m map[int]string) {
	mx.Lock()
	defer mx.Unlock()
	for k, v := range m {
		ve, ok := ErrorsCodes[k]
		if ok && v != ve {
			panic(fmt.Sprintf("mft.AddErrorsCodes, append Error `%v` fail. This code already exists."+
				"\n\tNew value:`%v`"+
				"\n\tExists value:`%v`", k, v, ve))
		}
		ErrorsCodes[k] = v
	}
}

// GenerateError - generate error with code
func GenerateError(key int, a ...interface{}) *Error {
	mx.RLock()
	defer mx.RUnlock()
	if text, ok := ErrorsCodes[key]; ok {
		return ErrorCSf(key, text, a...)
	}
	panic(fmt.Sprintf("mft.GenerateError, error not found code:%v", key))
}

// GenerateErrorE - generate error with code and exists error
func GenerateErrorE(key int, err error, a ...interface{}) *Error {
	mx.RLock()
	defer mx.RUnlock()
	if text, ok := ErrorsCodes[key]; ok {
		return ErrorCSEf(key, err, text, a...)
	}
	panic(fmt.Sprintf("mft.GenerateErrorE, error not found code:%v error:%v", key, err))
}

// GenerateErrorSubList -
func GenerateErrorSubList(key int, sub []*Error, a ...interface{}) *Error {
	if text, ok := ErrorsCodes[key]; ok {
		err := ErrorCSf(key, text, a...)
		err.InternalErrors = sub
		return err
	}
	panic(fmt.Sprintf("mft.GenerateErrorSubList, error not found code:%v", key))
}
