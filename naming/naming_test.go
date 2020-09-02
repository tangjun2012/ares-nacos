package naming_test

import (
	"github.com/tangjun2012/ares-nacos/naming"
	"testing"
)

func TestRegisterService(t *testing.T) {
	isSuccess, err := naming.RegisterService("localhost", 8080, "test", "", map[string]string{})
	println(isSuccess, err)
}
