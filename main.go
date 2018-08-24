package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/steelx/urlshortnr/config"
	"github.com/steelx/urlshortnr/handlers"
	"github.com/steelx/urlshortnr/storages"
)

func main() {
	configPath := flag.String("config", "./config/config.json", "Path for the config json file")
	flag.Parse()
	// Set use storage, select [Postgres, Filesystem, Redis ...]
	storage := &storages.Postgres{}

	// Read config
	config, err := config.FromFile(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	// Init storage
	if err = storage.Init(config); err != nil {
		log.Fatal(err)
	}

	// Defers
	defer storage.Close()

	// Handlers
	http.Handle("/", handlers.New(config.Options.Prefix, storage))

	// Create a server
	server := &http.Server{Addr: fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)}

	// Check for a closing signal
	go func() {
		// Graceful shutdown
		sigquit := make(chan os.Signal, 1)
		signal.Notify(sigquit, os.Interrupt, os.Kill)

		sig := <-sigquit
		log.Printf("caught sig: %+v", sig)
		log.Printf("Gracefully shutting down server...")

		if err := server.Shutdown(nil); err != nil {
			log.Printf("Unable to shut down server: %v", err)
		} else {
			log.Printf("Server stopped")
		}
	}()

	// Start server
	log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Println(err.Error())
	} else {
		log.Println("Server closed!")
	}
}
