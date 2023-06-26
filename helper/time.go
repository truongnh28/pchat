package helper

import (
	"errors"
	"time"
)

type TimeFormatString string

const APIClientDateTimeFormat TimeFormatString = "02/01/2006 15:04:05"
const ApiClientTimeFormat TimeFormatString = "15:04:05"
const ApiClientDateFormat TimeFormatString = "02/01/2006"

func (t TimeFormatString) IsValid() bool {
	switch t {
	case APIClientDateTimeFormat:
		return true
	}
	return false
}

func (t TimeFormatString) String() string {
	return string(t)
}

var (
	LocLocal *time.Location
)

func init() {
	l, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		panic(err)
	}
	LocLocal = l
}

func ParseLocalTime(layout TimeFormatString, value string) (time.Time, error) {
	return time.ParseInLocation(string(layout), value, LocLocal)
}

func ParseClientTime(value string, format TimeFormatString) (time.Time, error) {
	if !format.IsValid() {
		return time.Time{}, errors.New("time format is not valid")
	}
	return ParseLocalTime(format, value)
}

func FromTimeToString(value time.Time, format TimeFormatString) string {
	if !format.IsValid() {
		return ""
	}
	return value.Format(string(format))
}

func FromUnixSecondTimeString(sec int64, format TimeFormatString) string {
	return time.UnixMilli(sec).In(LocLocal).Format(string(format))
}
