package config

import (
	"encoding/json"
	"os"

	"github.com/labstack/gommon/log"
)

type Config struct {
	SPCliPath       string `json:"sp_cli_path"`
	SPCliPort       int    `json:"sp_cli_port"`
	StorageDBPath   string `json:"storage_db_path"`
	ContractAddress string `json:"contract_address"`
	HasGateway      bool   `json:"has_gateway"`
	Port            int    `json:"port"`
}

var StorageConfig Config = Config{
	SPCliPath:       "/home/ton-build/storage/storage-daemon/storage-daemon-cli",
	SPCliPort:       5555,
	StorageDBPath:   "/home/ton-build/storage-db",
	ContractAddress: "",
	Port:            33215,
}

func LoadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Info("Config file not found, using default config")
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	fileConfig := Config{}
	err = decoder.Decode(&fileConfig)
	if err != nil {
		log.Info("Config file is empty, using default config")
		return
	}

	if fileConfig.SPCliPath != "" {
		StorageConfig.SPCliPath = fileConfig.SPCliPath
	}
	if fileConfig.SPCliPort != 0 {
		StorageConfig.SPCliPort = fileConfig.SPCliPort
	}
	if fileConfig.StorageDBPath != "" {
		StorageConfig.StorageDBPath = fileConfig.StorageDBPath
	}
	if fileConfig.ContractAddress != "" {
		StorageConfig.ContractAddress = fileConfig.ContractAddress
	}
	if fileConfig.HasGateway {
		StorageConfig.HasGateway = fileConfig.HasGateway
	}
	if fileConfig.Port != 0 {
		StorageConfig.Port = fileConfig.Port
	}
}
