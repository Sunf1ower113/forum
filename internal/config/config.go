package config

import (
	"encoding/json"
	"os"
)

// type Database struct {
// 	DbName string `json:"db_name"`
// }

// type Listner struct {
// 	Protocol     string `json:"protocol"`
// 	Ip           string `json:"ip"`
// 	Port         string `json:"port"`
// 	IdleTimeout  int    `json:"idle_timeout"`
// 	WriteTimeout int    `json:"write_timeout"`
// 	ReadTimeout  int    `json:"read_timeout"`
// }

type Config struct {
	Listner struct {
		Protocol     string `json:"protocol"`
		Ip           string `json:"ip"`
		Port         string `json:"port"`
		IdleTimeout  int    `json:"idle_timeout"`
		WriteTimeout int    `json:"write_timeout"`
		ReadTimeout  int    `json:"read_timeout"`
	} `json:"listner"`
	Database struct {
		DbDriver string `json:"db_driver"`
		DbName string `json:"db_name"`
		Config string `json:"config"`
	} `json:"database"`
}

func LoadConfiguration(file string) (cfg *Config, err error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return
	}

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return
	}
	return
}
