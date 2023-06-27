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

func GetUniqueElements(a, b []string) []string {
	unique := make([]string, 0)
	bMap := make(map[string]bool)

	for _, element := range b {
		bMap[element] = true
	}

	for _, element := range a {
		if !bMap[element] {
			unique = append(unique, element)
		}
	}
	return unique
}
