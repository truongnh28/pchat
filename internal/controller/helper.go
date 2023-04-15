package controller

import (
	"chat-app/internal/common"
	"reflect"
)

var (
	fieldReturnCodes = []string{"ReturnCode", "ReturnMessage", "SubReturnCode", "SubReturnMessage"}
)

func buildResponse(subCode common.SubReturnCode, resp any) {
	code := getCode(subCode)

	stype := reflect.ValueOf(resp).Elem()
	for index, fieldName := range fieldReturnCodes {
		field := stype.FieldByName(fieldName)
		if field.IsValid() && field.CanSet() {
			switch index {
			case 0:
				field.SetInt(int64(code)) //set return_code
			case 1:
				field.SetString(code.Message()) //set return_message
			case 2:
				field.SetInt(int64(subCode)) //set sub_return_code
			case 3:
				field.SetString(subCode.Message()) //set sub_return_message
			}
		}
	}
}

func getCode(subCode common.SubReturnCode) common.ReturnCode {
	if subCode == common.OK {
		return common.Success
	}
	return common.Fail
}
