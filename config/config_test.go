package config_test

import (
	"github.com/swift9/ares-nacos/config"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestEnv(t *testing.T) {
	println(os.Getenv("TEST"))
}

func TestConfigString(t *testing.T) {
	println(config.GetString("test", "test"))
}

func TestGetVal(t *testing.T) {
	s := "a ${s.d  },${s.dd}"
	reg := regexp.MustCompile(`\$\{(\s+)?(\S)+(\s+)?\}`)
	i := reg.FindAllString(s, -1)
	for _, p := range i {
		regexp.MustCompile(`\$\{(\s+)?(\S)+(\s+)?\}`)
		s = strings.ReplaceAll(s, p, "88")
	}
	println(s)
}

func TestGetVal2(t *testing.T) {
	s := "${s1.s.s.a_d  }"
	reg := regexp.MustCompile(`(\w+(\.)?)+`)
	i := reg.FindAllString(s, -1)
	println(i)
}

func TestGetVal3(t *testing.T) {
	println(config.GetValue("1.${test}"))
}

func TestGetVal4(t *testing.T) {
	println(config.GetString("a.b"))
	println(config.GetString("a.c"))
}
