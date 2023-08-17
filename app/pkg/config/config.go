package config

import (
	_ "embed"
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	System string `json:"system"`
	Server struct {
		BindIP           string        `json:"bind_ip"`
		Port             string        `json:"port"`
		Timeout          time.Duration `json:"timeout"`
		ShutdownDuration time.Duration `json:"shutdown_duration"`
		StaticStorage    string        `json:"static_storage"`
		TempStorage      string        `json:"temp_storage"`
	} `json:"server"`
	Session struct {
		IDKey    string `json:"id_key"`
		Lifetime int    `json:"lifetime"`
	} `json:"session"`
}

func GetConfig(cfgPath *string) (*Config, error) {
	file, err := os.Open(*cfgPath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, err
}
