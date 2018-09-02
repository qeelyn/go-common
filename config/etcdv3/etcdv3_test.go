package etcdv3_test

import (
	"bytes"
	"github.com/qeelyn/go-common/config/etcdv3"
	"github.com/qeelyn/go-common/config/options"
	"github.com/qeelyn/go-common/grpcx/registry"
	retcd3 "github.com/qeelyn/go-common/grpcx/registry/etcdv3"
	"github.com/spf13/viper"
	"os/exec"
	"testing"
)

type defaultRemoteProvider struct {
	provider      string
	endpoint      string
	path          string
	secretKeyring string
}

func (rp defaultRemoteProvider) Provider() string {
	return rp.provider
}

func (rp defaultRemoteProvider) Endpoint() string {
	return rp.endpoint
}

func (rp defaultRemoteProvider) Path() string {
	return rp.path
}

func (rp defaultRemoteProvider) SecretKeyring() string {
	return rp.secretKeyring
}

func ectdPut(t *testing.T) {
	cmd := exec.Command("/bin/sh", "-c", "cat ../../_fixtrue/data/config.yaml | etcdctl put go-common")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
}

func TestEtcdConfigProvider_WatchChannel(t *testing.T) {
	rg, err := retcd3.NewRegistry(registry.Dsn("127.0.0.1:2379"))
	if err != nil {
		t.Fatal(err)
	}
	etcdp, err := etcdv3.NewEtcdConfigProvider(&options.Options{Registry: rg})
	if err != nil {
		t.Fatal(err)
	}
	rp := defaultRemoteProvider{
		path: "go-common",
	}
	respc, _ := etcdp.WatchChannel(rp)
	i := 0
	func(rc <-chan *viper.RemoteResponse) {
		for {
			b := <-rc
			reader := bytes.NewReader(b.Value)
			v := viper.New()
			v.Unmarshal(reader)
			if v.GetString("appname") == "" {
				t.Error()
				return
			}
			i++
			if i == 2 {
				return
			}
		}
	}(respc)
}
