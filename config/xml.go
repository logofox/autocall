package config

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type XmlConfig struct {
}

func (this *XmlConfig) Prase(filename string, v interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}
