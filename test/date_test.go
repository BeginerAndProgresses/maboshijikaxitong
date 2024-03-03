package test

import (
	"fmt"
	"maboshijikaxitong/utils"
	"testing"
	"time"
)

func TestDateTime(t *testing.T) {
	fmt.Println(utils.DifferenceOfDay(time.Now().Format(time.DateOnly), "2023-07-21"))
}
