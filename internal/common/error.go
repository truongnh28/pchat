package common

import "errors"

type CustomError struct {
	Code SubReturnCode
	err  error
}

func (c *CustomError) Error() string {
	return c.err.Error()
}

func (c *CustomError) Is(err error) bool {
	if cErr, ok := err.(*CustomError); ok {
		return cErr.Code == c.Code
	}
	return false
}

func IsCustomError(err error, errCode SubReturnCode) bool {
	return errors.Is(err, &CustomError{Code: errCode})
}

func (c *CustomError) Unwrap() error {
	return c.err
}

func WrapError(code SubReturnCode, err error) *CustomError {
	return &CustomError{
		Code: code,
		err:  err,
	}
}

func ConvertErrorToCustomError(err error) (*CustomError, bool) {
	if cErr, ok := err.(*CustomError); ok {
		return cErr, true
	}
	return WrapError(SystemError, err), false
}

type LoginError error

var (
	InvalidAccount   = errors.New("account not valid")
	BlockedAccount   = errors.New("account has been blocked")
	LoginInfoInvalid = errors.New("wrong login information")
	LoginSystemError = errors.New("system error")
)

type FormDataError error

var (
	GetHttpCtxFail = errors.New("get httpCtx fail")
	FiledInvalid   = errors.New("field is valid")
	ParseDataFail  = errors.New("cloud not parse form data")
)
