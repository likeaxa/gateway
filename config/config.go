package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"sync"
)
type ProxyConfig struct {
	HttpProxy map[string]ProxyInfo `toml:"proxyInfo"`
}

type ProxyInfo struct {
	Host       string   `toml:"host"`
	Target     []string `toml:"target"`
	ObtainMode int      `toml:"obtainMode"`
}
var (
	cfg  ProxyConfig
	once sync.Once
)

func Get() *ProxyConfig {
	once.Do(func() {
		cfg = ProxyConfig{}
		filePath := "/Users/yaoxinjian/Documents/视频/go/miniGateway/bin/config/conf.txt"
		if _, err := toml.DecodeFile(filePath, &cfg); err != nil {
			panic(err)
		}
		fmt.Printf("读取配置文件: %s\n", filePath)
	})
	return &cfg
}