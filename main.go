package main

import (
	"minGateway/config"
	"minGateway/gateway"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	//运行服务
	srv := new(gateway.GateServer)
	err = srv.Run()

	if err != nil {
		panic(err)
	}
}
