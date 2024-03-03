package test

import (
	"fmt"
	"testing"
	"time"
)

func TestString2Time(t *testing.T) {
	s := "2023-08-10"
	a, _ := time.Parse(time.DateOnly, s)
	//b, _ := time.Parse("2006-01-02 15:04:05", time.Now().String())
	b, _ := time.Parse(time.DateOnly, s)
	fmt.Println("a", a)
	fmt.Println("b", b)
	d := a.Sub(b)

	fmt.Println("相差时间：", d.Hours()/24, "天")
}
