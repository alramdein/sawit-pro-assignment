package handler

import (
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/labstack/echo/v4"
)

func (s *Server) GetProfile(ctx echo.Context) error {
	claims, err := s.HasAccess(ctx)
	if err != nil {
		s.Logger.Error(err)
		return ctx.JSON(http.StatusForbidden, Error(err))
	}

	userProfile, err := s.Repository.GetUserByID(ctx.Request().Context(), claims.UserId)
	if err != nil {
		s.Logger.Error(err)
		return ctx.JSON(http.StatusInternalServerError, Error(err))
	}
	if userProfile == nil {
		return ctx.JSON(http.StatusNotFound, Error(ErrUserNotFound))
	}

	resp := generated.GetProfileResponse{
		Name:  userProfile.Name,
		Phone: userProfile.Phone,
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) UpdateProfile(ctx echo.Context) error {
	claims, err := s.HasAccess(ctx)
	if err != nil {
		s.Logger.Error(err)
		return ctx.JSON(http.StatusForbidden, Error(err))
	}

	var input generated.UpdateProfileJSONBody
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, Error(err))
	}

	if input.Phone != nil && *input.Phone != "" {
		user, err := s.Repository.GetUserByPhone(ctx.Request().Context(), *input.Phone)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, InternalError())
		}
		if user != nil {
			return ctx.JSON(http.StatusConflict, Error(ErrPhoneAlreadyExist))
		}
	}

	if err := s.Repository.UpdateUser(ctx.Request().Context(), claims.UserId, input); err != nil {
		s.Logger.Error(err)
		return ctx.JSON(http.StatusInternalServerError, Error(err))
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "Profile updated successfully"})
}
