package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	ElectrumxServer string
	ServerAddress   string
}

var Conf Config

func InitConf() {
	if _, err := toml.DecodeFile("config.toml", &Conf); err != nil {
		log.Fatal(err)
	}
}
