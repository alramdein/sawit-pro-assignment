package handler

import (
	"crypto/rsa"

	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Repository    repository.RepositoryInterface
	Logger        echo.Logger
	JwtPublicKey  *rsa.PublicKey
	JwtPrivateKey *rsa.PrivateKey
}

type NewServerOptions struct {
	Repository    repository.RepositoryInterface
	Logger        echo.Logger
	JwtPublicKey  *rsa.PublicKey
	JwtPrivateKey *rsa.PrivateKey
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository:    opts.Repository,
		Logger:        opts.Logger,
		JwtPublicKey:  opts.JwtPublicKey,
		JwtPrivateKey: opts.JwtPrivateKey,
	}
}
