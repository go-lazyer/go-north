package generator

import (
	"bytes"
	"os"
	"strings"
)

func IsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err == nil {
		return true
	}
	if os.IsExist(err) {
		return true
	}
	return false
}

//只能创建目录，不能创建文件
func CreateDir(path string) error {
	if IsExist(path) {
		return nil
	}
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

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
func IsNumberType(inter interface{}) bool {
	if inter == nil {
		return false
	}
	switch inter.(type) {
	case int:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true
	case float32:
		return true
	case float64:
		return true
	default:
		return false
	}
}
