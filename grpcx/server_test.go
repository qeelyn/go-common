package grpcx_test

import (
	"github.com/qeelyn/go-common/grpcx"
	"testing"
)

func TestMicro(t *testing.T) {
	_, err := grpcx.Micro()
	if err != nil {
		t.Error(err)
	}
}
