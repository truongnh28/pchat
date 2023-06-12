package common

import (
	"fmt"
)

type ReturnCode int

const (
	Success ReturnCode = iota + 1
	Fail
)

var (
	returnCodeText = map[ReturnCode]string{
		Success: "Success",
		Fail:    "Failure",
	}
	subReturnCodeText = map[SubReturnCode]string{
		OK:              "Success",
		SystemError:     "System has been error, please try again later!",
		NotSupport:      "Currently, system does not support your request, please check and try again!",
		NotPermission:   "You have not permission to perform this request, please check and try again!",
		InvalidRequest:  "Invalid request, please check and try again!",
		NotFound:        "The resource you requested cannot be found, please check the request and try again!",
		ValidationError: "Your request does not pass resource validation tests, please check and try again",
	}
)

func (r ReturnCode) Message() string {
	return returnCodeText[r]
}

func (r ReturnCode) Int32() int32 {
	return int32(r)
}

type SubReturnCode int

const (
	OK SubReturnCode = iota + 1000
	SystemError
	NotSupport
	NotPermission
	InvalidRequest
	NotFound
	ValidationError
)

func (r SubReturnCode) Message() string {
	msg, ok := subReturnCodeText[r]
	if ok {
		return msg
	}
	return subReturnCodeText[SystemError]
}

func (r SubReturnCode) Int32() int32 {
	return int32(r)
}

var (
	NotFoundErr     error  = fmt.Errorf("not found")
	CookieName      string = "pchat"
	PrefixLoginCode string = "pchat"
)
