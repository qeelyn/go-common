package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// LoadConfig loads configuration from the given list of paths and populates it into the Config variable.
// The configuration file(s) should be named as app.yaml.
// Example (Use App Model):
//	isDebug := app.GetBool("debug")
//  // appmode has three mode: debug,prod,test,use it in your project
//  evn := app.GetString("appmode")
func LoadConfig(configFile string) (*viper.Viper, error) {
	//var filename, ext string = "app", "yaml"
	realPath, _ := filepath.Abs(configFile)
	file, err := os.Stat(realPath)
	if err != nil {
		return nil, err
	}
	configPath := path.Dir(realPath)
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
	if _, err := os.Stat(localConfig); err != nil {
		return nil, err
	}

	cnf.SetConfigName(filename + "-local")
	if err := cnf.MergeInConfig(); err != nil {
		return nil, err
	}

	switch cnf.GetString("appmode") {
	case "debug":
		cnf.Set("debug", true)
	}

	return cnf, nil
}
