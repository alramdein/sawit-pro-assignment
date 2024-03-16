package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	e := echo.New()

	t.Run("ValidRegistration", func(t *testing.T) {
		mockRepo.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(nil).Do(func(ctx context.Context, input repository.InsertUserInput) {
			assert.NotEmpty(t, input.Name)
			assert.NotEmpty(t, input.Phone)
			assert.NotEmpty(t, input.Password)
		})

		s := NewServer(NewServerOptions{
			Repository: mockRepo,
			Logger:     e.Logger,
		})

		reqBody := map[string]interface{}{
			"name":     "Alif",
			"phone":    "+62812345670",
			"password": "Password123!",
		}
		reqJSON, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = s.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InvalidPhoneNumberLength", func(t *testing.T) {
		s := NewServer(NewServerOptions{
			Repository: mockRepo,
			Logger:     e.Logger,
		})

		reqBody := map[string]interface{}{
			"name":     "Alif",
			"phone":    "+62812234324234324",
			"password": "Password123!",
		}
		reqJSON, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		s.Register(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("DialCodeMismatch", func(t *testing.T) {
		s := NewServer(NewServerOptions{
			Repository: mockRepo,
			Logger:     e.Logger,
		})

		reqBody := map[string]interface{}{
			"name":     "Alif",
			"phone":    "+14155552671", // US dial code instead of INA
			"password": "Password123!",
		}
		reqJSON, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		s.Register(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidNameLength", func(t *testing.T) {
		s := NewServer(NewServerOptions{
			Repository: mockRepo,
			Logger:     e.Logger,
		})

		reqBody := map[string]interface{}{
			"name":     "A",
			"phone":    "+62812345670",
			"password": "Password123!",
		}
		reqJSON, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		s.Register(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidPasswordLength", func(t *testing.T) {
		s := NewServer(NewServerOptions{
			Repository: mockRepo,
			Logger:     e.Logger,
		})

		reqBody := map[string]interface{}{
			"name":     "Alif",
			"phone":    "+62812345670",
			"password": "Pass",
		}
		reqJSON, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		s.Register(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidPasswordFormat", func(t *testing.T) {
		s := NewServer(NewServerOptions{
			Repository: mockRepo,
			Logger:     e.Logger,
		})

		reqBody := map[string]interface{}{
			"name":     "Alif",
			"phone":    "+62812345670",
			"password": "password123",
		}
		reqJSON, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		s.Register(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
