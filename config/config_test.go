// need install etcdctl component to execute remote setting
package config_test

import (
	"github.com/qeelyn/go-common/config"
	"github.com/qeelyn/go-common/config/etcdv3"
	"github.com/qeelyn/go-common/config/options"
	"github.com/qeelyn/go-common/grpcx/registry"
	etcdv32 "github.com/qeelyn/go-common/grpcx/registry/etcdv3"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"os/exec"
	"testing"
)

func init() {
	cmd := exec.Command("/bin/sh", "-c", "cat ../_fixtrue/data/config.yaml | etcdctl put go-common/config.yaml")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	// 读取输出结果
	_, err = ioutil.ReadAll(stdout)
	if err != nil {
		panic(err)
	}
}

func TestParseOptions(t *testing.T) {
	opt := config.ParseOptions(config.Path("config"))
	if opt == nil {
		t.Fatal()
	}
}

func TestLoadConfig(t *testing.T) {
	//local
	opts := &options.Options{
		Path:     "../_fixtrue/data",
		FileName: "config.yaml",
	}
	cnf, err := config.LoadConfig(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !cnf.IsSet("appmode") {
		t.Error("miss appmode")
	}
	//remote
	register, err := etcdv32.NewRegistry(registry.Dsn("127.0.0.1:2379"))
	opts = &options.Options{
		Path:     "go-common",
		FileName: "config.yaml",
		Registry: register,
	}
	if opts.Registry != nil {
		etcdv3.Build(opts)
	}
	cnf, err = config.LoadConfig(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !cnf.IsSet("appmode") {
		t.Error("miss appmode")
	}
}

func TestLoadLocalConfig(t *testing.T) {
	cnf, err := config.LoadLocalConfig(&options.Options{Path: "../_fixtrue/data/config.yaml"})
	if err != nil {
		t.Fatal(err)
	}
	if !cnf.IsSet("appmode") {
		t.Error("miss appmode")
	}
}

func TestLoadRemoteConfig(t *testing.T) {

	register, err := etcdv32.NewRegistry(registry.Dsn("127.0.0.1:2379"))
	if err != nil {
		t.Fatal(err)
	}
	cnf, err := config.LoadRemoteConfig(&options.Options{
		Registry: register,
		Path:     "go-common",
	})
	if err != nil {
		t.Fatal(err)
	}
	if cnf.GetString("appname") == "" {
		t.Fatal()
	}
	if cnf.GetStringMap("log.file") == nil {
		t.Fatal()
	}
	err = cnf.WatchRemoteConfigOnChannel()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetRemotePath(t *testing.T) {
	cmd := exec.Command("/bin/sh", "-c", "etcdctl put go-common/key abcd")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	// 读取输出结果
	_, err = ioutil.ReadAll(stdout)
	if err != nil {
		panic(err)
	}
	register, err := etcdv32.NewRegistry(registry.Dsn("127.0.0.1:2379"))
	if err != nil {
		t.Fatal(err)
	}
	//remote
	opts := &options.Options{
		Path:     "go-common",
		FileName: "config.yaml",
		Registry: register,
	}
	etcdv3.Build(opts)
	cnf, err := config.LoadConfig(opts)
	if err != nil {
		t.Fatal(err)
	}
	io, err := config.GetRemotePath(cnf, "config-local.yaml")
	if err != nil {
		t.Fatal(err)
	}
	val, err := ioutil.ReadAll(io)
	if err != nil {
		t.Error(err)
	}
	if len(val) == 0 {
		t.Error()
	}
}

func TestResetFromSource(t *testing.T) {
	var (
		err      error
		register registry.Registry
		stdout   io.ReadCloser
		cmd      *exec.Cmd
	)
	opts := &options.Options{
		Path:     "../_fixtrue/data",
		FileName: "config.yaml",
	}
	if viper.RemoteConfig == nil {
		goto check1
	}
remote:
	cmd = exec.Command("/bin/sh", "-c", "etcdctl put config-local.yaml abcd")
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	// 读取输出结果
	_, err = ioutil.ReadAll(stdout)
	if err != nil {
		panic(err)
	}
	register, err = etcdv32.NewRegistry(registry.Dsn("127.0.0.1:2379"))
	if err != nil {
		t.Fatal(err)
	}
	//remote
	opts = &options.Options{
		Path:     "go-common",
		FileName: "config.yaml",
		Registry: register,
	}
	etcdv3.Build(opts)
check1:
	cnf, err := config.LoadConfig(opts)
	if err != nil {
		t.Fatal(err)
	}
	if viper.RemoteConfig == nil {
		authConfig := cnf.GetStringMap("auth")
		authConfig["private-key"] = "./" + authConfig["private-key"].(string)

	}
	count := len(cnf.GetStringMap("auth"))
	err = config.ResetFromSource(cnf, "auth.private-key")
	if err != nil {
		t.Error(err)
	} else {
		if len(cnf.GetStringMap("auth")) != count {
			t.Error("sub key changes")
		}
		if !cnf.GetBool("auth.keep") {
			t.Error("key miss")
		}
		_, ok := cnf.Get("auth.private-key").([]byte)
		if !ok {
			t.Error()
		}
	}
	if viper.RemoteConfig == nil {
		goto remote
	}
}
