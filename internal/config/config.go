package config

import (
	"os"
	"encoding/json"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DB_url             string  `json:"db_url"`
	Current_user_name  string  `json:"current_user_name"`
}

func (c *Config) SetUser(user string) error {
	c.Current_user_name = user
	config_path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	byte_data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	os.WriteFile(config_path, byte_data, 664)
	return nil
}

func ReadConfig() (Config, error) {
	var conf Config
	config_path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	byte_data, err := os.ReadFile(config_path)
	if err != nil {
		return Config{}, err
	}

	err = json.Unmarshal(byte_data, &conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}


func getConfigFilePath() (string, error) {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	full_conf_file_path := home_dir+"/"+configFileName
	return full_conf_file_path, nil
}
