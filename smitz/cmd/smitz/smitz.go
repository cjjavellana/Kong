package main

import (
	"cjavellana.me/kong/smitz/internal/pkg/cfg"
	"cjavellana.me/kong/smitz/internal/pkg/cyclops"
	"cjavellana.me/kong/smitz/internal/pkg/ipc"
	"cjavellana.me/kong/smitz/internal/pkg/rpc"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	config := cfg.ReadConfig()
	log.Printf("Cyclops Url: %s \n", config.CyclopsUrl)

	mgtServer := cyclops.New(config)
	err := mgtServer.Register()
	if err != nil {
		log.Printf("Unable to register self in Cyclops, %v\n", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.SmitzPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	ipc.RegisterAdminServer(s, rpc.New(config.KongAdminUrl))
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
