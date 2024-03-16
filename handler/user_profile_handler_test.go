package handler

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	e := echo.New()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	publicKey := &privateKey.PublicKey

	t.Run("Success", func(t *testing.T) {
		userProfile := &model.User{ID: uuid.NewString(), Name: "Alif Ramdani", Phone: "+623456789"}

		mockRepo.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(userProfile, nil)

		s := NewServer(NewServerOptions{
			Repository:    mockRepo,
			Logger:        e.Logger,
			JwtPublicKey:  publicKey,
			JwtPrivateKey: privateKey,
		})

		req := httptest.NewRequest(http.MethodGet, "/user/profile", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		jwt, err := s.GenerateJWTToken(userProfile.ID)
		assert.NoError(t, err)

		token := "Bearer " + jwt
		c.Request().Header.Set(echo.HeaderAuthorization, token)

		err = s.GetProfile(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp generated.GetProfileResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.Equal(t, userProfile.Name, resp.Name)
		assert.Equal(t, userProfile.Phone, resp.Phone)
	})
}

func TestUpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	e := echo.New()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	publicKey := &privateKey.PublicKey

	t.Run("Success", func(t *testing.T) {
		userProfile := &model.User{ID: uuid.NewString(), Name: "Alif Ramdani", Phone: "+623456789"}

		phoneInput := "123456789"
		mockInput := generated.UpdateProfileJSONBody{Phone: &phoneInput}

		mockRepo.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).Return(nil, nil)
		mockRepo.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.AssignableToTypeOf(mockInput)).Return(nil)

		s := NewServer(NewServerOptions{
			Repository:    mockRepo,
			Logger:        e.Logger,
			JwtPublicKey:  publicKey,
			JwtPrivateKey: privateKey,
		})

		reqJSON, err := json.Marshal(mockInput)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/user/profile", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		jwt, err := s.GenerateJWTToken(userProfile.ID)
		assert.NoError(t, err)

		token := "Bearer " + jwt
		c.Request().Header.Set(echo.HeaderAuthorization, token)

		err = s.UpdateProfile(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("DuplicatePhone", func(t *testing.T) {
		userProfile := &model.User{ID: uuid.NewString(), Name: "Alif Ramdani", Phone: "+623456789"}

		phoneInput := "123456789"
		mockInput := generated.UpdateProfileJSONBody{Phone: &phoneInput}

		mockRepo.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).Return(userProfile, nil)

		s := NewServer(NewServerOptions{
			Repository:    mockRepo,
			Logger:        e.Logger,
			JwtPublicKey:  publicKey,
			JwtPrivateKey: privateKey,
		})

		reqJSON, err := json.Marshal(mockInput)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/user/profile", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		jwt, err := s.GenerateJWTToken(userProfile.ID)
		assert.NoError(t, err)

		token := "Bearer " + jwt
		c.Request().Header.Set(echo.HeaderAuthorization, token)

		s.UpdateProfile(c)
		assert.Equal(t, http.StatusConflict, rec.Code)
	})
}
