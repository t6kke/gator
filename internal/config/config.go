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

func (c *Config) SetUser(user string) {
	c.Current_user_name = user
	config_path := getConfigFilePath()
	byte_data, _ := json.Marshal(c)  //should do error handling
	os.WriteFile(config_path, byte_data, 664)
}

func ReadConfig() Config {
	var conf Config
	config_path := getConfigFilePath()
	byte_data, _ := os.ReadFile(config_path)  //should do error handling
	json.Unmarshal(byte_data, &conf)  //should do error handling
	return conf
}


func getConfigFilePath() string {
	home_dir, _ := os.UserHomeDir()  //should do error handling
	full_conf_file_path := home_dir+"/"+configFileName
	return full_conf_file_path
}
