package test

import (
	"fmt"
	"testing"
)

func TestAny(t *testing.T) {
	a := any("nihao")
	switch a.(type) {
	case string:
		fmt.Println(a.(string))
	}
}
