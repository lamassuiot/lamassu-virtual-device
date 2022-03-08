package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Device struct {
		Id                  string `json:"id"`
		Alias               string `json:"alias"`
		EnrollingProperties struct {
			Subject struct {
				Country            string `json:"country"`
				State              string `json:"state"`
				Locality           string `json:"locality"`
				Organization       string `json:"organization"`
				OrganizationalUnit string `json:"organization_unit"`
			} `json:"subject"`
		} `json:"enrolling_properties"`
	} `json:"device"`
	Dms struct {
		DmsEstServer string `json:"est_server"`
	} `json:"dms"`
	Aws struct {
		IotCoreEndpoint string `json:"iot_core_endpoint"`
		TestLambda      string `json:"test_lambda"`
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
