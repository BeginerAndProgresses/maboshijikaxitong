package utils

import (
	"os"
)

// IsExist
// 如果存在 返回ture，否则返回flase
func IsExist(filepath string) bool {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// CreatFile 创建文件
func CreatFile(filepath string) bool {
	// Create 可读可写不可执行 666
	_, err := os.Create(filepath)
	if err != nil {
		return false
	}
	return true
}
