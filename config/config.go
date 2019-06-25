package config

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

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

func ParceConfig(configPath string) error {
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

	return nil
}
