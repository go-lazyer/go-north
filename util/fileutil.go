package util

import "os"

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
