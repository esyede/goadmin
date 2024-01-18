package util

import (
	"encoding/json"
	"fmt"
)

func Struct2Json(obj interface{}) string {
	str, err := json.Marshal(obj)

	if err != nil {
		panic(fmt.Sprintf("[Struct2Json] Error: %v", err))
	}

	return string(str)
}

func Json2Struct(str string, obj interface{}) {
	err := json.Unmarshal([]byte(str), obj)

	if err != nil {
		panic(fmt.Sprintf("[Json2Struct] Error: %v", err))
	}
}

func JsonI2Struct(str interface{}, obj interface{}) {
	JsonStr := str.(string)
	Json2Struct(JsonStr, obj)
}
