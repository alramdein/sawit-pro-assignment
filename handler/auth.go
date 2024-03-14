package handler

import (
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/labstack/echo/v4"
)

func (s *Server) Register(ctx echo.Context) error {
	var resp generated.RegisterResponse

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) Login(ctx echo.Context) error {
	var resp generated.LoginResponse

	return ctx.JSON(http.StatusOK, resp)
}
