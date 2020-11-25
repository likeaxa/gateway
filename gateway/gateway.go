package gateway

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type ObtainMode int //多转发目标时的选择模式

const (
	SelectModeRandom ObtainMode = 1 //随机选择
	SelectModePoll   ObtainMode = 2 //轮询选择
)

type HostInfo struct {
	Target          string     //转发目标域名
	MultiTarget     []string   //有多转发目标的域名集合
	IsMultiTarget   bool       //是否有多转发目标
	MultiTargetMode ObtainMode //多转发目标选择模式
	PoolModeIndex   int        //轮询模式索引
}

//默认转发地址
var DefaultTarget *HostInfo

var HostList map[string]HostInfo

func init() {
	HostList = make(map[string]HostInfo)
}

type GateServer struct{}
type Proxy struct{}

func (hostInfo *HostInfo) GetTarget(req *http.Request) string {
	var route string
	if hostInfo.IsMultiTarget {
		if hostInfo.MultiTargetMode == SelectModeRandom { //随机模式
			route = hostInfo.MultiTarget[rand.Int()%len(hostInfo.MultiTarget)]
		} else if hostInfo.MultiTargetMode == SelectModePoll { //轮询模式
			route = hostInfo.MultiTarget[hostInfo.PoolModeIndex]
			hostInfo.PoolModeIndex++
			hostInfo.PoolModeIndex = hostInfo.PoolModeIndex % len(hostInfo.MultiTarget)
		} else { //未配置或配置错误使用随机模式
			route = hostInfo.MultiTarget[rand.Int()%len(hostInfo.MultiTarget)]
		}
	} else {
		route = hostInfo.Target
	}
	return route
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	in := time.Now()
	//根据配置选择转发到哪
	var route string //转发的目标
	var existRoute = false
	if len(r.Host) == 0 {
		if DefaultTarget != nil {
			route = DefaultTarget.GetTarget(r)
			existRoute = true
		}
	} else if hostInfo, ok := HostList[r.Host]; ok {
		route = hostInfo.GetTarget(r)
		existRoute = true
	}
	if !existRoute {
		if DefaultTarget != nil {
			route = DefaultTarget.GetTarget(r)
			existRoute = true
		} else {
			fmt.Println("获取不到路由")
			return
		}
	}
	//找到转发目标，继续
	if existRoute {
		target, err := url.Parse(route)
		if err != nil {
			fmt.Println("url.Parse失败")
			return
		}

		proxy := newHostReverseProxy(target)
		proxy.ServeHTTP(w, r)
	}
	fmt.Println("耗时：", time.Now().Sub(in).Milliseconds(), "毫秒")
}
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
func newHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

func (s *GateServer) proxy8080() error {

	p := &Proxy{}
	fmt.Println("网关监听端口:8080")
	err := http.ListenAndServe(":8080", p)
	if err != nil {
		panic(err)
	}
	return nil
}
func (s *GateServer) Run() error {
	err := s.proxy8080()
	if err != nil {
		panic(err)
	}
	return nil
}
