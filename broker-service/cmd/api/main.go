package main

import (
	"fmt"
	"github.com/loidinhm31/go-micro/common"
	"log"
	"net/http"
)

type Config struct{}

func main() {
	app := Config{}

	log.Printf("Starting broker service on port %s\n", common.BrokerPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", common.BrokerPort),
		Handler: app.routes(),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
