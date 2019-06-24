package main

import (
	"bluemine/config"
	"flag"
	"database/sql"
	"log"

	_ "github.com/cockroachdb/cockroach-go/crdb"
)

func main() {
	var (
		err        error
		configPath string
	)

	flag.StringVar(&configPath, "c", "config.toml", "Path to server configuration")
	flag.Parse()

	config.ParceConfig(configPath)

	db, err := sql.Open("postgres", config.Conf.Postgresql)
	if err != nil {
		log.Fatal("Can't connect to database " + err.Error())
	}
	defer db.Close()

}