package main

import (
	"electrumx-proxy-go/config"
	"electrumx-proxy-go/router"
	"electrumx-proxy-go/ws"
	"log"
)

func main() {
	config.InitConf()
	ws.InitWebSocket(config.Conf.ElectrumxServer)
	api := router.InitMasterRouter()
	err := api.Run(config.Conf.ServerAddress)
	if err != nil {
		log.Fatal(err)
		return
	}
}
