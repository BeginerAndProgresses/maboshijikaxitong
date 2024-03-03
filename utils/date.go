package utils

import (
	"fmt"
	"strings"
	"time"
)

// NowDate 获取当前日子
func NowDate() (year, month, day uint) {
	// 设置时区为东八区
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone
	now := time.Now()
	year = uint(now.Year())
	month = uint(now.Month())
	day = uint(now.Day())
	return
}

// NowWeek 获取当前周几
func NowWeek() int {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone
	now := time.Now()
	return int(now.Weekday())
}

// DifferenceOfDay 返回 a - b 的天数
// a与b的格式例如 2023-01-13，
// 可以直接传入当前日期 time.Now().Format(time.DateOnly)
func DifferenceOfDay(a, b string) int {
	t1, _ := time.Parse(time.DateOnly, a)
	t2, _ := time.Parse(time.DateOnly, b)
	days := t1.Sub(t2).Hours() / 24
	return int(days)
}

// DateFormat 判断获取到的日期格式，如果日期格式为08-17-23则转为2023-08-17这样的格式
func DateFormat(date string) string {
	if date == "" {
		return date
	}
	if len(date) <= 8 {
		split := strings.Split(date, "-")
		s := time.Now().String()[0:2]
		return fmt.Sprintf("%s%s-%s-%s", s, split[2], split[0], split[1])
	}
	return date
}

// DateOfNDays 获取后n天的日期，如果n为负数则输出前-n天的日期
func DateOfNDays(n int) string {
	return time.Now().AddDate(0, 0, n).Format(time.DateOnly)
}
