package main

import (
	"context"
	"example/pb"
	"log"
	"net"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

const (
	USER_ADDDED_QUEUE = "h8_p3_user_added"
)

var (
	DB            *mongo.Database
	UserRepo      UserRepository
	Authenticator AuthInterceptor
)

func initiateDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DB_URL")))

	if err != nil {
		log.Fatal(err.Error())
	}

	DB = client.Database(os.Getenv("DB_NAME"))
	UserRepo = NewUserRepository(DB)
}

func main() {
	initiateDatabase()

	Authenticator = NewAuthInterceptor()
	publisher := NewQueuePublisher(os.Getenv("RABBIT_URL"), USER_ADDDED_QUEUE)

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(Authenticator.AuthtenticateApp()),
	)
	userServer := NewUsersServer(UserRepo, publisher)
	PORT := os.Getenv("USER_SERVICE_PORT")

	pb.RegisterUsersServer(srv, userServer)

	log.Println("Starting RPC server at", PORT)

	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("could not listen to %s: %v", PORT, err)
	}

	log.Fatal(srv.Serve(listener))
}
