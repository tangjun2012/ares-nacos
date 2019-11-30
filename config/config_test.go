package config_test

import (
	"github.com/swift9/ares-nacos/config"
	"os"
	"testing"
)

func TestEnv(t *testing.T) {
	println(os.Getenv("TEST"))
}

func TestConfigString(t *testing.T) {
	println(config.GetString("test", "test"))
}
