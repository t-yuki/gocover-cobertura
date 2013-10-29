package testdata

import (
	"testing"
)

func TestFunc1(t *testing.T) {
	val := 0
	Func1(&val)
	if val != 0 {
		t.Fail()
	}
}
