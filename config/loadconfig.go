package config

import (
	"github.com/g4zhuj/hashring"
	"minGateway/gateway"
	"strings"
)

func LoadConfig() error {
	httpProxy := Get().HttpProxy

	for _, v := range httpProxy {
		info := getHostInfo(v)
		if v.Host == "default" {
			//如果定义了default, 遇到未知host走这里
			gateway.DefaultTarget = &info
		} else {
			gateway.HostList[v.Host] = info
			if strings.HasPrefix(v.Host, "www.") {
				if strings.Count(v.Host, ".") == 2 {
					//一级域名，考虑没有带"www"的情况
					gateway.HostList[strings.TrimLeft(v.Host, "www.")] = gateway.HostList[v.Host]
				}
			} else if strings.Count(v.Host, ".") == 1 {
				//排除首位和末位的"."，"."的数量只有一个说明是没有带"www"的一级域名
				gateway.HostList["www."+v.Host] = gateway.HostList[v.Host]
			}
		}
	}
	return nil
}

func getHostInfo(proxyInfo ProxyInfo) gateway.HostInfo {
	var hostInfo gateway.HostInfo
	// Target 里面的空格去掉
	for i := 0; i < len(proxyInfo.Target); i++ {
		proxyInfo.Target[i] = strings.ReplaceAll(proxyInfo.Target[i], " ", "")
	}
	if len(proxyInfo.Target) == 0 {
		panic(proxyInfo.Host + " : len(proxyInfo.Target) == 0")
	} else if len(proxyInfo.Target) == 1 {
		hostInfo = gateway.HostInfo{IsMultiTarget: false, Target: proxyInfo.Target[0]}
	} else {
		//定义了多个目标，使用分流
		targets := proxyInfo.Target
		hostInfo = gateway.HostInfo{IsMultiTarget: true, MultiTarget: targets, MultiTargetMode: gateway.ObtainMode(proxyInfo.ObtainMode)}
		if proxyInfo.ObtainMode == 3 { //哈希模式
			//把节点放到hashring中，同时设置权重
			hostInfo.HashRing = hashring.NewHashRing(100) //vitualSpots=100
			nodeWeight := make(map[string]int)
			for _, target := range targets {
				nodeWeight[target] = 1 //这里简化了，权重都设置为1
			}
			hostInfo.HashRing.AddNodes(nodeWeight)
		}
	}
	return hostInfo
}
