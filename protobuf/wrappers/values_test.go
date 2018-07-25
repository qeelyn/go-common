package wrappers_test

import (
	"github.com/qeelyn/go-common/protobuf/wrappers"
	"testing"
)

func TestWrapFloat64(t *testing.T) {
	var a *float64
	v1 := wrappers.WrapFloat64(a)
	if v1 != nil {
		t.Error("nil float64 wrap error")
	}
	a = new(float64)
	*a = 9.1234
	v1 = wrappers.WrapFloat64(a)
	if float64(*v1) != *a {
		t.Error("float64 wrap error")
	}
}
