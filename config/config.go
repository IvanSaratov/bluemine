package config

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

//Conf storing main datas
var Conf struct {
	Postgresql   string
	Memcache     string
	Bind         string
	BindTLS      string
	Host         string
	LdapServer   string
	LdapBaseDN   string
	LdapUser     string
	LdapPassword string
	SessionKey   string
}

//ParseConfig to parse .toml config
func ParseConfig(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	if _, err = toml.Decode(string(contents), &Conf); err != nil {
		return err
	}

	log.Println("Config parsed!")

	return nil
}
