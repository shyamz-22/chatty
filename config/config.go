package config

import (
	"encoding/json"
	"log"
	"io/ioutil"
	"os"
)

type AuthConfig struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	CallbackUrl  string `json:"callback_url"`
	DiscoveryUrl string `json:"discovery_url"`
	SecretKey    string `json:"secret_key"`
	CookieAge    int    `json:"cookie_age"`
}

type Configuration struct {
	Auth    AuthConfig `json:"auth"`
	Profile string     `json:"profile"`
}

type ConfigurationLoader interface {
	Load() *Configuration
}

type JsonConfigLoader struct {
	data []byte
}

func NewJsonConfigLoader(fileName string) *JsonConfigLoader {
	dataBytes, err := ioutil.ReadFile("configuration.json")

	if err != nil {
		log.Printf("Invalid config file: %v \n", err)
		os.Exit(1)
	}
	return &JsonConfigLoader{
		data: dataBytes,
	}
}

func (jcl *JsonConfigLoader) Load() *Configuration {

	var config = &Configuration{}
	err := json.Unmarshal(jcl.data, config)

	if err != nil {
		log.Printf("Invalid configuration: %v \n", err)
	}

	return config
}
