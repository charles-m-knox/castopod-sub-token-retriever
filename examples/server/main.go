package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	flagConfig  string
	flagTest    bool
	flagOutFile string
	flagAddr    string
)

func parseFlags() {
	flag.StringVar(&flagConfig, "f", "config.json", "json file to use for loading configuration")
	flag.BoolVar(&flagTest, "test", false, "connect read-only and perform a dry run")
	flag.StringVar(&flagOutFile, "o", "", "a file to write the database query to (can combine with -test to allow manual editing of the query)")
	flag.StringVar(&flagAddr, "addr", "0.0.0.0:19281", "the address (host and port) to listen on")
	flag.Parse()
}

func main() {
	parseFlags()

	var err error
	conf, err = LoadConfig(flagConfig)
	if err != nil {
		log.Fatalf("failed to load config: %v", err.Error())
	}

	DB = getDB(conf.SQLConnectionString, flagTest)

	srv := &http.Server{
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 3 * time.Second,
		Addr:         flagAddr,
	}

	go clearCodes()

	http.Handle("/", http.HandlerFunc(router))

	log.Printf("listening on %v", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to run web server: %v", err.Error())
	}
}
