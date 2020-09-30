package main

import (
	"cjavellana.me/kong/smitz/internal/pkg/cfg"
	"cjavellana.me/kong/smitz/internal/pkg/cyclops"
	"cjavellana.me/kong/smitz/internal/pkg/ipc"
	"cjavellana.me/kong/smitz/internal/pkg/rpc"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":50051"
)

func main() {
	config := cfg.ReadConfig()
	log.Printf("Cyclops Url: %s \n", config.CyclopsUrl)

	mgtServer := cyclops.New(config)
	err := mgtServer.Register()
	if err != nil {
		log.Printf("Unable to register self in Cyclops, %v\n", err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	ipc.RegisterSmitzServer(s, &rpc.Smitz{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
