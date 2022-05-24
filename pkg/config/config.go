package config

import (
	"encoding/json"
	"os"
)

type PrivateKeyMetadata struct {
	KeyType string `json:"type"`
	KeyBits int    `json:"bits"`
}
type PrivateKeyMetadataWithStregth struct {
	KeyType     string `json:"type"`
	KeyBits     int    `json:"bits"`
	KeyStrength string `json:"strength"`
}
type Config struct {
	CertificatesDirectory string `json:"certificates_dir"`
	Devmanager            struct {
		EstServer string `json:"est_server"`
		Cert      string `json:"cert"`
	} `json:"devmanager"`
	Aws struct {
		IotCoreEndpointCACertFile string `json:"iot_core_ca_file"`
		IotCoreEndpoint           string `json:"iot_core_endpoint"`
		TestLambda                string `json:"test_lambda"`
	} `json:"aws"`
}

func NewConfig() (Config, error) {
	f, err := os.Open("config.json")
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&cfg)
	return cfg, err
}
