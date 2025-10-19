package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gogatorconfig.json"

type Config struct {
	DBURL    string `json:"db_url"`
	UserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(name string) error {
	cfg.UserName = name
	return write(*cfg)
}

func Read() (Config, error) {
	fpath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	f, err := os.Open(fpath)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func getConfigFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fpath := filepath.Join(dir, configFileName)
	return fpath, nil
}

func write(cfg Config) error {
	fpath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	f, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(cfg); err != nil {
		return err
	}
	return nil
}
