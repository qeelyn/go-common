package config_test

import (
	"github.com/qeelyn/go-common/config"
	"github.com/qeelyn/go-common/config/etcdv3"
	"github.com/qeelyn/go-common/config/options"
	"io/ioutil"
	"log"
	"os/exec"
	"testing"
)

func TestParseOptions(t *testing.T) {
	opt := config.ParseOptions(config.Addrs("127.0.0.1:2379"), config.Secure(true))
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
	opts = &options.Options{
		Addrs:    []string{"127.0.0.1:2379"},
		Path:     "go-common",
		FileName: "config.yaml",
	}
	etcdv3.Build(opts)
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
	cmd := exec.Command("/bin/sh", "-c", "cat ../_fixtrue/data/config.yaml | etcdctl put go-common")
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
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(string(opBytes))

	cnf, err := config.LoadRemoteConfig(&options.Options{
		Addrs: []string{"127.0.0.1:2379"},
		Path:  "go-common",
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
