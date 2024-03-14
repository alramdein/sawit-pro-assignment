package handler

import (
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/labstack/echo/v4"
)

func (s *Server) GetProfile(ctx echo.Context) error {
	var resp generated.GetProfileResponse

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) UpdateProfile(ctx echo.Context) error {
	var resp generated.UpdateProfileResponse

	return ctx.JSON(http.StatusOK, resp)
}
