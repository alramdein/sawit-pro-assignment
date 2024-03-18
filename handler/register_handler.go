package handler

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/helper"
	"github.com/SawitProRecruitment/UserService/helper/constants"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s *Server) Register(ctx echo.Context) error {
	var req generated.RegisterJSONBody
	var resp generated.RegisterResponse

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Error(err))
	}

	valid, errDetail := validateRegisterInput(req)
	if !valid {
		return ctx.JSON(http.StatusBadRequest, ErrorWithDetail(ErrInvaInput, errDetail))
	}

	user, err := s.Repository.GetUserByPhone(ctx.Request().Context(), req.Phone)
	if err != nil {
		s.Logger.Error(err)
		return ctx.JSON(http.StatusInternalServerError, InternalError())
	}
	if user != nil {
		return ctx.JSON(http.StatusConflict, Error(ErrPhoneAlreadyExist))
	}

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		s.Logger.Error(err)
		return ctx.JSON(http.StatusInternalServerError, InternalError())
	}

	newUserId := uuid.NewString()

	err = s.Repository.InsertUser(ctx.Request().Context(), repository.InsertUserInput{
		Id:       newUserId,
		Name:     req.Name,
		Phone:    req.Phone,
		Password: hashedPassword,
	})
	if err != nil {
		s.Logger.Error(err)
		return ctx.JSON(http.StatusInternalServerError, InternalError())
	}

	resp.Id = newUserId
	return ctx.JSON(http.StatusOK, resp)
}

func validateRegisterInput(req generated.RegisterJSONBody) (isValid bool, errDetail *RegisterUserError) {
	var errFields RegisterUserError
	var errFieldsName []string
	var errFieldsPhone []string
	var errFieldsPassword []string

	if req.Phone == "" {
		errFieldsPhone = append(errFieldsPhone, ErrPhoneRequired.Error())
	}
	if req.Name == "" {
		errFieldsName = append(errFieldsName, ErrNameRequired.Error())
	}
	if req.Password == "" {
		errFieldsPassword = append(errFieldsPassword, ErrPasswordRequired.Error())
	}

	if len(req.Phone) < 10 || len(req.Phone) > 13 {
		errFieldsPhone = append(errFieldsPhone, ErrRangePhoneNumber.Error())
	}

	dialCode := req.Phone[:3]
	if dialCode != constants.INA_DIAL_CODE {
		errFieldsPhone = append(errFieldsPhone, ErrNotIndonesianPhoneNumber.Error())
	}

	if len(req.Name) < 3 || len(req.Phone) > 60 {
		errFieldsName = append(errFieldsName, ErrRangeName.Error())
	}

	if len(req.Password) < 6 || len(req.Password) > 64 {
		errFieldsPassword = append(errFieldsPassword, ErrRangePassword.Error())
	}

	// Regular expressions to check for at least 1 capital letter, 1 number, and 1 special character
	hasCapital := regexp.MustCompile(`[A-Z]`).MatchString
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString
	hasSpecial := regexp.MustCompile(`[^A-Za-z0-9]`).MatchString

	if !hasCapital(req.Password) || !hasNumber(req.Password) || !hasSpecial(req.Password) {
		errFieldsPassword = append(errFieldsPassword, ErrUnsatisfiedPassword.Error())
	}

	if len(errFieldsName) > 0 || len(errFieldsPhone) > 0 || len(errFieldsPassword) > 0 {
		errFields.Name = strings.Join(errFieldsName, ",")
		errFields.Phone = strings.Join(errFieldsPhone, ",")
		errFields.Password = strings.Join(errFieldsPassword, ",")
		return false, &errFields
	}

	return true, nil
}
