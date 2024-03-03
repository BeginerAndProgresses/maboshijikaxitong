package utils

import (
	"strings"
)

func PathTransform(path string) string {
	return strings.Replace(path, `\`, "/", -1)
}

// Last4Rune 获取字符串后四位
func Last4Rune(str string) string {
	s := []rune(str)
	if len(s) <= 4 {
		return str
	}
	last4 := string(s[len(s)-4:])
	return last4
}

// OneRuneIsNumber 判断第一个字符是否是数字
func OneRuneIsNumber(str string) bool {
	anyReplaceList := []int32{'.', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	for strIndex, strInt := range str[:1] {
		if strIndex == 0 {
			if strInt == '-' {
				continue
			}
		}
		if IsHaveReplaceList(strInt, anyReplaceList) {
			continue
		}
		return false
	}
	return true
}

// IsHaveReplaceList 是否在过滤列表中
func IsHaveReplaceList(strInt int32, replaceList []int32) bool {
	for _, chStrInt := range replaceList {
		if chStrInt == strInt {
			return true
		}
	}
	return false
}
