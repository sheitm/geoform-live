package main

import (
	"github.com/sheitm/ofever/scrape"
	"github.com/sheitm/ofever/storage"
	"log"
	"os"
)

func main(){
	storageDirectory := os.Getenv("STORAGE_DIRECTORY")
	if storageDirectory == "" {
		log.Fatal("environment variable STORAGE_DIRECTORY must be set")
	}
	seasonChan := make(chan *scrape.SeasonFetch)

	storage.Start(storageDirectory, seasonChan)

	// startServer must be last line
	startServer("2112", seasonChan)
}
