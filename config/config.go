package config

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type (
	Config struct {
		Postgresql string
		Memcache   string
		Bind       string
	}
)

var (
	Conf Config
)

func ParceConfig(configPath string) {
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatal("Can't open config file " + err.Error())
	}

	defer file.Close()

	contents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Can't read contents from config file " + err.Error())
	}

	if _, err = toml.Decode(string(contents), &Conf); err != nil {
		log.Fatal("Can't parce config file " + err.Error())
	}
}