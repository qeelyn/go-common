package config

import (
	"fmt"
	"github.com/qeelyn/go-common/config/options"
	"github.com/qeelyn/go-common/grpcx/registry"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Path(p string) options.Option {
	return func(o *options.Options) {
		o.Path = p
	}
}

func FileName(f string) options.Option {
	return func(o *options.Options) {
		o.FileName = f
	}
}

func Registry(r registry.Registry) options.Option {
	return func(o *options.Options) {
		o.Registry = r
	}
}

func ParseOptions(opts ...options.Option) *options.Options {
	opt := &options.Options{}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

// if use remote must set NewRemoteFunc
func LoadConfig(opts *options.Options) (*viper.Viper, error) {
	if opts.Registry != nil {
		return LoadRemoteConfig(opts)
	}
	return LoadLocalConfig(opts)
}

// LoadLocalConfig loads configuration from the given list of paths and populates it into the Config variable.
// The configuration file(s) should be named as app.yaml.
// Example (Use App Model):
//	isDebug := app.GetBool("debug")
//  // appmode has three mode: debug,prod,test,use it in your project
//  evn := app.GetString("appmode")
func LoadLocalConfig(opts *options.Options) (*viper.Viper, error) {
	//var filename, ext string = "app", "yaml"
	configFile := path.Join(opts.Path, opts.FileName)
	realPath, _ := filepath.Abs(configFile)
	file, err := os.Stat(realPath)
	if err != nil {
		return nil, err
	}
	configPath := path.Dir(filepath.ToSlash(realPath))
	fn := strings.Split(file.Name(), ".")
	filename := fn[0]
	ext := fn[1]
	cnf := viper.New()
	//cnf.WatchConfig()
	cnf.SetConfigName(filename)
	cnf.SetConfigType(ext)
	cnf.AutomaticEnv()

	cnf.AddConfigPath(configPath)
	cnf.SetDefault("debug", false)

	if err := cnf.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Failed to read the configuration file: %s", err)
	}
	// local
	localConfig := path.Join(configPath, filename+"-local."+ext)
	if _, err := os.Stat(localConfig); err == nil {
		cnf.SetConfigName(filename + "-local")
		if err := cnf.MergeInConfig(); err != nil {
			return nil, err
		}
	}
	defaultSet(cnf)
	return cnf, nil
}

// LoadRemoteConfig loads configuration from the config center and populates it into the Config variable.
// To enable remote support in Viper, do a blank import of the viper/remote package:
// current only support etcd
func LoadRemoteConfig(opts *options.Options) (*viper.Viper, error) {
	cnf := viper.New()
	cnf.AddRemoteProvider("etcd", opts.Addr, path.Clean(opts.Path)+"/"+opts.FileName)
	cnf.SetConfigType("yaml")
	if err := cnf.ReadRemoteConfig(); err != nil {
		return nil, err
	}
	defaultSet(cnf)
	return cnf, nil
}

func defaultSet(cnf *viper.Viper) {
	switch cnf.GetString("appmode") {
	case "debug":
		cnf.Set("debug", true)
	}
}
