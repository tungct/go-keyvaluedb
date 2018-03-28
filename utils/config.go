package utils

import (
	"os"
	"log"
	"github.com/BurntSushi/toml"
)

type ConfigRedis struct {
	REDIS_ADDR                  string
	PASSWORD                    string
	DB			    int
}

type ConfigServer struct{
	IP 			       string
	PORT                           string
	PATHLEVELDB		       string
}


//function load all config of ConfigRedis struct
func LoadConfigRedis(nameFileConfig string) ConfigRedis{
	var config ConfigRedis
	_, err := os.Stat(nameFileConfig)
	if err != nil {
		log.Fatal("Config file is missing: ", nameFileConfig)
	}

	if _, err := toml.DecodeFile(nameFileConfig, &config); err != nil {
		log.Fatal("Wrong parameters!")
		log.Fatal(err)
	}
	return config
}

//function load all config of ConfigRedis struct
func LoadConfigServer(nameFileConfig string) ConfigServer{
	var config ConfigServer
	_, err := os.Stat(nameFileConfig)
	if err != nil {
		log.Fatal("Config file is missing: ", nameFileConfig)
	}

	if _, err := toml.DecodeFile(nameFileConfig, &config); err != nil {
		log.Fatal("Wrong parameters!")
		log.Fatal(err)
	}
	return config
}