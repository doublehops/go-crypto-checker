package main

import (
	"encoding/json"
	"fmt"
	"os"

	"btcmwatch.local/btcmarketsmodule"
)

type apiConfig struct {
	ApiKey     string
	PrivateKey string
}

func main() {

	config, err := getConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	btcmodule := btcmarketsmodule.CreateInstance(config.ApiKey, config.PrivateKey)
	btcmodule.PrintCurrencies()
}

func getConfig() (apiConfig, error) {
	configFile := "config/config.json"
	configuration := apiConfig{}

	file, err := os.Open(configFile)
	if err != nil {
		return configuration, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		return configuration, err
	}

	return configuration, nil
}
