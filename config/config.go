package config

import (
	"encoding/json"
	"log"
	"os"
)

var (
	HttpPort  int    // http service port
	SolanaRpc string // solana endpoint url
)

// Load load config file
func Load() {
	file, err := os.Open("config/config.json")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	type Config struct {
		HttpPort  int    `json:"http_port"`
		SolanaRpc string `json:"solana_rpc"`
	}
	all := &Config{}
	if err = json.NewDecoder(file).Decode(all); err != nil {
		log.Panic(err)
	}
	HttpPort = all.HttpPort
	SolanaRpc = all.SolanaRpc
}
