package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type Configuration struct {
	DatabaseURL     string `envconfig:"DATABASE_URL" required:"true"`
	TestDatabaseURL string `envconfig:"TEST_DATABASE_URL" default:""`
}

func InitializeTestEnv() *Env {
	var c Configuration
	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := InitializeDB(c.TestDatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	return &Env{db: db}
}

func InitializeProdEnv() *Env {
	var c Configuration
	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := InitializeDB(c.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	return &Env{db: db}
}

func main() {
	serve := flag.Bool("serve", false, "initialize server")
	scrape := flag.Bool("scrape", false, "initialize scrapper")

	flag.Parse()

	if len(os.Args) > 1 {
		if flag.NFlag() != 1 {
			fmt.Println("pass just one argument")
			flag.Usage()
			os.Exit(1)
		}

		e := InitializeProdEnv()
		defer e.db.Close()

		e.migrateDB()

		if *serve {
			e.run()
		}

		if *scrape {
			httpClient := &http.Client{Timeout: 10 * time.Second}
			s := &Scrapper{httpClient: httpClient, db: e.db}

			err := s.Scrape()
			if err != nil {
				log.Fatal(err.Error())
			}
		}

	} else {
		flag.Usage()
	}
}
