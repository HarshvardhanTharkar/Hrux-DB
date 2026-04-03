package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"kv-distributed/internal/api"
	"kv-distributed/internal/datastructures"
	"kv-distributed/internal/indexing"
	"kv-distributed/internal/service"
	"kv-distributed/internal/storage"
)

func main() {
	port := flag.String("port", ":8080", "Server port")
	dataFile := flag.String("data", "kvdata.gob", "Data file path")
	flag.Parse()

	// Initialize layers
	storage := storage.NewStorage()
	indexer := indexing.NewIndexer()
	dataStructs := datastructures.NewDataStructuresService()

	// Create service
	kvService := service.NewKVService(storage, indexer, dataStructs)

	// Create and start API server
	server := api.NewKVServer(kvService)

	// Load existing data if any
	if err := kvService.LoadFromFile(*dataFile); err != nil {
		log.Println("No existing data file found, starting fresh")
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down server...")
		if err := kvService.SaveToFile(*dataFile); err != nil {
			log.Printf("Error saving data: %v", err)
		}
		kvService.Stop()
		os.Exit(0)
	}()

	log.Printf("Starting KV RPC Server on %s", *port)

	// Start the RPC server
	if err := server.Start(*port); err != nil {
		log.Fatal("Server error:", err)
	}
}
