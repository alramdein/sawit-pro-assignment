package main

import (
	"log"
	"os"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/helper"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		if os.IsNotExist(err) {
			log.Println(".env file not found, using environment variables from Docker")
		} else {
			log.Fatalf("Error loading .env file: %s", err)
		}
	}

	dbDsn := os.Getenv("DATABASE_URL")
	repo := repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})

	publicKey, err := helper.GetRSAPublicKey()
	if err != nil {
		log.Fatal("invalid rsa public key: ", err)
	}

	privateKey, err := helper.GetRSAPrivateKey()
	if err != nil {
		log.Fatal("invalid rsa private key: ", err)
	}

	opts := handler.NewServerOptions{
		Repository:    repo,
		Logger:        echo.New().Logger,
		JwtPublicKey:  publicKey,
		JwtPrivateKey: privateKey,
	}
	return handler.NewServer(opts)
}
