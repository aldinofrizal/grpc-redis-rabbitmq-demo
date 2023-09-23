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

const (
	USER_ADDDED_QUEUE = "h8_p3_user_added"
	USERS_CACHE_KEY   = "h8_p3_users"
)

func init() {
	USER_SERVICE_PORT = os.Getenv("USER_SERVICE_PORT")
	API_PORT = os.Getenv("MAIN_API_PORT")
}

func generateUserService() pb.UsersClient {
	auth := NewUserServiceInterceptor()
	conn, err := grpc.Dial(
		USER_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(auth.EmbedAuthCredentials()),
	)

	if err != nil {
		log.Fatal("could not connect to", USER_SERVICE_PORT, err)
	}

	return pb.NewUsersClient(conn)
}

func main() {
	e := echo.New()

	cache := NewCacheStorage(os.Getenv("REDIS_URL"))

	queueHandler := NewQueueHandler(cache)
	consumer := NewQueueConsumer(os.Getenv("RABBIT_URL"))
	consumer.AddConsumer(USER_ADDDED_QUEUE, queueHandler.UserAddedHandler)

	userService := generateUserService()
	userHandler := NewUserHandler(userService, cache)
	auth := NewAuthenticator(userService, cache)

	e.POST("/users/register", userHandler.Register)
	e.POST("/users/login", userHandler.GetToken)
	e.GET("/users", userHandler.FindAll, auth.Authenticate)

	e.Logger.Fatal(e.Start(API_PORT))
}
