package main

import (
	"os"
	"solana_rpc/api"
	"solana_rpc/config"
	"solana_rpc/daemon"
)

func main() {
	if os.Args[len(os.Args)-1] != "-d" {
		os.Args = append(os.Args, "-d")
	}
	daemon.Background("./out.log", true)

	config.Load()
	api.StartHttpServer(config.HttpPort)

	daemon.WaitForKill()
}
