package testdata

import (
	"testing"
)

func TestFunc2a(t *testing.T) {
	var a Type1
	val := 2
	a.Func2a(&val)
	if val != 1 {
		t.Fail()
	}
}
