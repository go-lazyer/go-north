package util

import (
	"bytes"
	"strings"
)

//首字母大写驼峰
func ToUpperCamelCase(str string) string {
	if len(str) == 0 {
		return ""
	}
	strArr := strings.Split(str, "_")
	var sb bytes.Buffer
	for _, s := range strArr {
		if len(s) == 0 {
			continue
		}
		sb.WriteString(strings.ToUpper(s[0:1]) + s[1:])
	}
	return sb.String()
}

//首字母小写驼峰
func ToLowerCamelCase(str string) string {
	if len(str) == 0 {
		return ""
	}
	strArr := strings.Split(str, "_")
	var sb bytes.Buffer
	for n, s := range strArr {
		if len(s) == 0 {
			continue
		}
		if n == 0 {
			sb.WriteString(s)
		} else {
			sb.WriteString(strings.ToUpper(s[0:1]) + s[1:])
		}
	}
	return sb.String()
}
