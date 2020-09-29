package main

import (
	"cjavellana.me/kong/smitz/internal/pkg/cfg"
	"cjavellana.me/kong/smitz/internal/pkg/cyclops"
	"log"
)

func main() {
	config := cfg.ReadConfig()
	log.Printf("Cyclops Url: %s \n", config.CyclopsUrl)
	
	mgtServer := cyclops.New(config)
	err := mgtServer.Register()
	if err != nil {
		log.Fatalf("Unable to register self in Cyclops, %v\n", err)
	}
}
