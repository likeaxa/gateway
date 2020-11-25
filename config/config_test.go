package config

import (
	"fmt"
	"miniGateway/gateway"
	"testing"
)

func TestBase(t *testing.T) {
	get := Get()
	proxy := get.HttpProxy
	for k, v := range proxy {
		fmt.Print(k, ":->")
		fmt.Println(v.Host, ":", v.Target, ":", v.ObtainMode)
		fmt.Print("\n")
	}
}
func TestLoadConfig(t *testing.T) {
	err := LoadConfig()
	if err != nil {
		panic(err)
	}
	list := gateway.HostList
	for k, v := range list {
		fmt.Print(k, ":->")
		if v.IsMultiTarget {
			fmt.Print(k, "->", v.MultiTarget, " mode:", v.MultiTargetMode)
		} else {
			fmt.Print(k, "->", v.Target)
		}
		fmt.Print("\n")
	}
}
