package handler

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SawitProRecruitment/UserService/helper"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/stretchr/testify/assert"
)

type LoginTestStruct struct {
	name           string
	phone          string
	password       string
	expectedStatus int
	expectError    bool
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	publicKey := &privateKey.PublicKey

	e := echo.New()

	t.Run("ValidCredentials", func(t *testing.T) {
		tc := &LoginTestStruct{
			name:           "ValidCredentials",
			phone:          "+62818288122",
			password:       "password123",
			expectedStatus: http.StatusOK,
			expectError:    false,
		}

		hashedPassword, err := helper.HashPassword(tc.password)
		assert.NoError(t, err)
		mockUser := &model.User{ID: uuid.NewString(), Password: hashedPassword}
		mockRepo.EXPECT().GetUserByPhone(gomock.Any(), tc.phone).Return(mockUser, nil)
		mockRepo.EXPECT().IncrementLoginCount(gomock.Any(), mockUser.ID).Return(nil)

		s := NewServer(NewServerOptions{
			Repository:    mockRepo,
			Logger:        e.Logger,
			JwtPublicKey:  publicKey,
			JwtPrivateKey: privateKey,
		})

		reqBody := map[string]interface{}{
			"phone":    tc.phone,
			"password": tc.password,
		}
		reqJSON, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = s.Login(c)

		if tc.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, tc.expectedStatus, rec.Code)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		tc := &LoginTestStruct{
			name:           "ValidCredentials",
			phone:          "+62818288122",
			password:       "password123",
			expectedStatus: http.StatusOK,
			expectError:    false,
		}

		hashedPassword, err := helper.HashPassword("different_pass")
		assert.NoError(t, err)
		mockUser := &model.User{ID: uuid.NewString(), Password: hashedPassword}
		mockRepo.EXPECT().GetUserByPhone(gomock.Any(), tc.phone).Return(mockUser, nil)

		s := NewServer(NewServerOptions{
			Repository:    mockRepo,
			Logger:        e.Logger,
			JwtPublicKey:  publicKey,
			JwtPrivateKey: privateKey,
		})

		reqBody := map[string]interface{}{
			"phone":    tc.phone,
			"password": tc.password,
		}
		reqJSON, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = s.Login(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
