package config_test

import (
	"github.com/qeelyn/go-common/config"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	cnf, err := config.LoadConfig("../_fixtrue/data/config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if !cnf.IsSet("appmode") {
		t.Error("miss appmode")
	}
}
