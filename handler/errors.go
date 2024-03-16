package handler

import "errors"

var (
	ErrGeneralInternal = errors.New("something went wrong")

	ErrPhoneRequired           = errors.New("phone required")
	ErrNameRequired            = errors.New("name required")
	ErrPasswordRequired        = errors.New("password required")
	ErrInvalidPassword         = errors.New("invalid password")
	ErrInvaInput               = errors.New("invalid input")
	ErrInvalidLoginCredentials = errors.New("phone or password is incorrect")
	ErrUserNotFound            = errors.New("user not found")
	ErrPhoneAlreadyExist       = errors.New("phone already exist")

	ErrRangePhoneNumber         = errors.New("phone numbers must be at minimum 10 characters and maximum 13 characters")
	ErrNotIndonesianPhoneNumber = errors.New("phone numbers must start with the Indonesia country code “+62”.")
	ErrRangeName                = errors.New("full name must be at minimum 3 characters and maximum 60 characters")
	ErrRangePassword            = errors.New("passwords must be minimum 6 characters and maximum 64 characters")
	ErrUnsatisfiedPassword      = errors.New("passwords must be containing at least 1 capital characters AND 1 number AND 1 special (non alpha-numeric) characters")

	ErrInvalidJWT = errors.New("invalid JWT token")
	ErrMissingJWT = errors.New("missing JWT token")
)

type ErrorResponse struct {
	Message     string      `json:"message"`
	ErrorDetail interface{} `json:"errorDetail,omitempty"`
}

type RegisterUserError struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

func Error(err error) ErrorResponse {
	return ErrorResponse{
		Message: err.Error(),
	}
}

func ErrorWithDetail(err error, detail interface{}) ErrorResponse {
	return ErrorResponse{
		Message:     err.Error(),
		ErrorDetail: detail,
	}
}

func InternalError() ErrorResponse {
	return ErrorResponse{
		Message: ErrGeneralInternal.Error(),
	}
}
