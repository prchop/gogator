package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gogatorconfig.json"

// Settings is the default config
type Settings struct {
	DBURL    string `json:"db_url"`
	UserName string `json:"current_user_name"`
}

func (cfg *Settings) SetUser(name string) error {
	cfg.UserName = name
	return write(*cfg)
}

// Read config file
func Read() (Settings, error) {
	fpath, err := getConfigFilePath()
	if err != nil {
		return Settings{}, err
	}

	f, err := os.Open(fpath)
	if err != nil {
		return Settings{}, err
	}
	defer f.Close()

	var cfg Settings
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return Settings{}, err
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

func write(cfg Settings) error {
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
