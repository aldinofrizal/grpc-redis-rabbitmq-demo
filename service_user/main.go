package main

import (
	"example/pb"
	"log"
	"net"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

func main() {
	srv := grpc.NewServer()
	userServer := NewUsersServer()
	PORT := os.Getenv("USER_SERVICE_PORT")

	pb.RegisterUsersServer(srv, userServer)

	log.Println("Starting RPC server at", PORT)

	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("could not listen to %s: %v", PORT, err)
	}

	log.Fatal(srv.Serve(listener))
}
