package main

import (
	"example/pb"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	USER_SERVICE_PORT string
	API_PORT          string
)

func init() {
	USER_SERVICE_PORT = os.Getenv("USER_SERVICE_PORT")
	API_PORT = os.Getenv("MAIN_API_PORT")
}

func generateUserService() pb.UsersClient {
	conn, err := grpc.Dial(
		USER_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatal("could not connect to", USER_SERVICE_PORT, err)
	}

	return pb.NewUsersClient(conn)
}

func main() {
	e := echo.New()

	userService := generateUserService()
	userHandler := NewUserHandler(userService)

	e.POST("/users/register", userHandler.Register)
	e.GET("/users", userHandler.FindAll)

	e.Logger.Fatal(e.Start(API_PORT))
}
