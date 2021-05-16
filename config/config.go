package config

import "github.com/jinzhu/configor"

type Config struct {
	Twitter struct {
		Token            string
		Secret           string
		RequestURI       string
		AuthorizationURI string
		TokenRequestURI  string
		CallbackURI      string
	}
}

func New() *Config {
	config := new(Config)

	err := configor.Load(config, "./config/config.yaml")
	if err != nil {
		return nil
	}

	return config
}
