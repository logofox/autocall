package config

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	var c AppConfig
	c.load()
	fmt.Println("%v", c)

}
