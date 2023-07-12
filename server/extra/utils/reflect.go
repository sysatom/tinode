package utils

import (
	"reflect"
	"runtime"
	"strings"
)

func GetFunctionName(i interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	names := strings.Split(name, ".")
	if len(names) > 0 {
		return names[len(names)-1]
	}
	return ""
}
