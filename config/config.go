package config

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	// "smartcall.im/airgo/modules/logging"
)

var config *AppConfig

type AppConfig struct {
	ListName List `xml:"listName"`
}

type List struct {
	phoneList string `xml:"phone"`
}

func getCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	return filepath.Dir(file)
}

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// Config files by default on : conf/config.xml
func (this *AppConfig) load() error {
	var adapter XmlConfig
	var filename string

	filename = "conf/config.xml"
	if isExist(filename) {
		return adapter.Prase(filename, this)
	}
	return errors.New("Configuration File Not Found.")
}

func Default() *AppConfig {
	if config == nil {
		var cfg AppConfig
		if err := cfg.load(); err != nil {
			// logging.Fatal(err)
			os.Exit(1)
		}
		config = &cfg
	}
	return config
}
