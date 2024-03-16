package handler

import (
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) Login(ctx echo.Context) error {
	var req generated.LoginJSONBody
	var resp generated.LoginResponse

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, Error(err))
	}

	user, err := s.Repository.GetUserByPhone(ctx.Request().Context(), req.Phone)
	if err != nil {
		s.Logger.Error(err)
		return ctx.JSON(http.StatusBadRequest, Error(ErrInvalidLoginCredentials))
	}
	if user == nil {
		return ctx.JSON(http.StatusBadRequest, Error(ErrInvalidLoginCredentials)) // obfuscate
	}

	valid := validatePassword(user.Password, req.Password)
	if !valid {
		return ctx.JSON(http.StatusBadRequest, Error(ErrInvalidLoginCredentials))
	}

	token, err := s.GenerateJWTToken(user.ID)
	if err != nil {
		s.Logger.Error(err)
		return ctx.JSON(http.StatusInternalServerError, InternalError())
	}

	if err := s.Repository.IncrementLoginCount(ctx.Request().Context(), user.ID); err != nil {
		s.Logger.Error(err)
	}

	resp.Id = user.ID
	resp.JwtToken = token

	return ctx.JSON(http.StatusOK, resp)
}

func validatePassword(hashedPasswordFromDB string, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPasswordFromDB), []byte(inputPassword))
	if err != nil {
		return false
	}
	return true
}
