package main

import (
	"cjavellana.me/kong/smitz/internal/pkg/cfg"
	"cjavellana.me/kong/smitz/internal/pkg/cyclops"
	"fmt"
)

func main() {
	config := cfg.ReadConfig()
	fmt.Printf("Hello %s \n", config.CyclopsUrl)

	mgtServer := cyclops.New(config)
	mgtServer.Register()
}
