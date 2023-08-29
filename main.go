package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"yyyoichi/Collo-API/proto/collo"

	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	srv := &Server{}

	collo.RegisterColloServiceServer(s, srv)
	log.Printf("start gRPC server port: %v", port)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
