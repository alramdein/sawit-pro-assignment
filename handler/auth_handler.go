package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/SawitProRecruitment/UserService/model"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func (s *Server) GenerateJWTToken(userID string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(s.JwtPrivateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Server) HasAccess(ctx echo.Context) (claims model.JwtClaims, err error) {
	authHeader := ctx.Request().Header.Get("Authorization")
	if authHeader == "" {
		return claims, ErrMissingJWT
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.JwtPublicKey, nil
	})
	if err != nil {
		return claims, ErrInvalidJWT
	}

	if !token.Valid {
		return claims, ErrInvalidJWT
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.Logger.Error("failed to parse claims")
		return claims, ErrInvalidJWT
	}

	if sub, ok := mapClaims["sub"].(string); ok {
		claims.UserId = sub
	} else {
		s.Logger.Error("invalid user id in jwt token")
		return claims, ErrInvalidJWT
	}

	if exp, ok := mapClaims["exp"].(string); ok {
		claims.ExpiresAt = exp
	}

	return claims, nil
}
