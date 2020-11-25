package config

func GetMockConfigData() []ProxyInfo {
	configList := []ProxyInfo{
		{
			Host:       "127.0.0.1:8080",
			Target:     []string{"http://127.0.0.1:7788", "http://127.0.0.1:7788"},
			ObtainMode: 2,
		}, {
			Host:       "default",
			Target:     []string{"http://127.0.0.1:7788", "http://127.0.0.1:7789"},
			ObtainMode: 1,
		},
	}
	return configList
}
