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
