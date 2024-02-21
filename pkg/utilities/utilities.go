package utilities

import (
	"TaskManager/pkg/logger"
	"bytes"
	"encoding/json"
)

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "   ")
	if err != nil {
		return in
	}
	return out.String()
}

func ToJSON(object any) string {
	jsonByte, err := json.Marshal(object)
	if err != nil {
		logger.Error("Ошибка при получении JSON: ", err.Error())
	}
	n := len(jsonByte)
	result := string(jsonByte[:n])

	return jsonPrettyPrint(result)
}
