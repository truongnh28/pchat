package helper

import "strings"

type GenderType uint8

const (
	Male GenderType = iota
	Female
)

func (rt GenderType) String() string {
	switch rt {
	case Male:
		return "male"
	case Female:
		return "female"
	}
	panic("invalid gender type")
}

func ConvertToGenderType(s string) GenderType {
	switch strings.ToLower(s) {
	case "male":
		return Male
	case "female":
		return Female
	}
	panic("invalid redis type")
}

func SafeConvertToGenderType(s string) (rType GenderType, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()

	return ConvertToGenderType(s), true
}
