package config

type ProxyConfig struct {
	HttpProxy map[string]ProxyInfo
}

type ProxyInfo struct {
	Host       string
	Target     []string
	ObtainMode int
}

func Get() *ProxyConfig {
	data := GetMockConfigData()
	cfg := ProxyConfig{}
	cfg.HttpProxy = make(map[string]ProxyInfo)
	for i := 0; i < len(data); i++ {
		info := data[i]
		cfg.HttpProxy[info.Host] = info
	}
	return &cfg
}
