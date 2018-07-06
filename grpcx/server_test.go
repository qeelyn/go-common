package grpcx_test

import (
	"github.com/qeelyn/go-common/grpcx"
	"testing"
)

func TestMicro(t *testing.T) {
	_, err := grpcx.Micro("test", grpcx.WithPrometheus(":9100"))
	if err != nil {
		t.Error(err)
	}
}

func TestMutilServer(t *testing.T) {
	a, err := grpcx.Micro("test", grpcx.WithPrometheus(":9100"))
	if err != nil {
		t.Error(err)
	}
	b, err := grpcx.Micro("test", grpcx.WithPrometheus(":9101"))
	if err != nil {
		t.Error(err)
	}
	arpc := a.BuildGrpcServer()
	//arpc.RegisterService(nil,nil)
	a.StartPrometheus(nil)
	go func() {
		a.Run(arpc, "9009")
	}()

	brpc := a.BuildGrpcServer()
	//brpc.RegisterService(nil,nil)
	b.StartPrometheus(nil)
	b.Run(brpc, "9010")
}
