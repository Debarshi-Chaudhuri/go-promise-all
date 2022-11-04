package util

import (
	"fmt"
	"reflect"
)

func Assert(condition bool, errorCode string, statusCode int, errorMsg string) {
	if !condition {
		panic(fmt.Sprint(statusCode) + ":" + string(errorCode) + ":" + errorMsg)
	}
}
func Ternary[T any](condition bool, valueIfTrue T, valueIfFalse T) T {
	if condition {
		return valueIfTrue
	} else {
		return valueIfFalse
	}
}

func GetMsgFromError(err interface{}) string {
	var msg string

	if reflect.TypeOf(err).Kind().String() == "error" {
		msg = err.(error).Error()
	} else if reflect.TypeOf(err).Kind() == reflect.String {
		msg = reflect.ValueOf(err).String()
	} else {
		msg = reflect.TypeOf(err).String()
	}

	return msg
}