package config

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

//Conf storing main datas
var Conf struct {
	DBHost        string `toml:"dbHost"`
	Memcache      string `toml:"Memcache"`
	Port          string `toml:"Port"`
	Host          string `toml:"Host"`
	LdapServer    string `toml:"ldapServer"`
	LdapBaseDN    string `toml:"ldapBaseDN"`
	LdapUser      string `toml:"ldapUser"`
	LdapPassword  string `toml:"ldapPassword"`
	SessionKey    string `toml:"sessionKey"`
	EncryptionKey string `toml:"encryptionKey"`
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

	return nil
}
