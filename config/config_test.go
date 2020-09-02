package config_test

import (
	"github.com/tangjun2012/ares-nacos/config"
	"github.com/tidwall/gjson"
	"os"
	"testing"
)

func TestEnv(t *testing.T) {
	println(os.Getenv("TEST"))
}

func TestConfig(t *testing.T) {
	println(config.GetString("test"))
}

func TestConfigString(t *testing.T) {
	data := gjson.Parse(`{"a":123}`)
	config.InitByJson(&data)
	println(config.GetString("a"))
}
