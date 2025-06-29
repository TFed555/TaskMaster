package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"auth-service/internal/controllers"
	"auth-service/internal/controllers/mocks"
	"auth-service/internal/models"

	// "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthController_Registration(t *testing.T) {
	cases := []struct {
		name          string
		requestBody   string
		mockSetup     func(authMock *mocks.AuthService)
		expectedCode  int
		expectedError string
	}{
		{
			name: "Success",
			requestBody: `{
				"login": "testuser",
				"email": "user@example.com",
				"password": "password123"
			}`,
			mockSetup: func(authMock *mocks.AuthService) {
				authMock.On("Register", "testuser", "user@example.com", "password123").
					Return(&models.User{ID: 1}, &models.Tokens{
						AccessToken:  "access_token_example",
						RefreshToken: "refresh_token_example",
					}, nil).
					Once()
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "No password",
			requestBody: `{
				"login": "testuser",
				"email": "user@example.com"
			}`,
			mockSetup:    func(authMock *mocks.AuthService) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Constraint email",
			requestBody: `{
				"login": "testuser",
				"email": "user@example.com",
				"password": "password123"
			}`,
			mockSetup: func(authMock *mocks.AuthService) {
				authMock.On("Register", "testuser", "user@example.com", "password123").
					Return(nil, nil, errors.New("some error")).
					Once()
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "Registration failed",
		},

	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			authMock := mocks.NewAuthService(t)
			tc.mockSetup(authMock)
			controllerMock := controllers.NewAuthController(authMock)
			req := httptest.NewRequest(http.MethodPost, "/registration", bytes.NewBufferString(tc.requestBody))
			rr := httptest.NewRecorder()
			controllerMock.Register(rr, req)
			require.Equal(t, tc.expectedCode, rr.Code)

			if tc.expectedError != "" {
				var resp controllers.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Contains(t, resp.Message, tc.expectedError)
			}
		})
	}
}

func TestAuthController_Login(t *testing.T) {
	cases := [] struct {
		name string
		requestBody string
		mockSetup func (authMock *mocks.AuthService)
		expectedCode int
		expectedError string
	} {
		{
			name: "Success",
			requestBody: `{
				"email": "user@example.com",
				"password": "password123"
			}`,
			mockSetup: func(authMock *mocks.AuthService) {
				authMock.On("Authorize", "user@example.com", "password123").
					Return(&models.User{ID: 1, Email: "user@example.com", Login: "user"}, &models.Tokens{
						AccessToken:  "access_token_example",
						RefreshToken: "refresh_token_example",
					}, 
					nil,
					true).
					Once()
			},
			expectedCode: http.StatusAccepted,
		},
		{ 
			name: "Wrong email",
			requestBody: `{
				"email": "user@example.com",
				"password": "password123"
			}`,
			mockSetup: func(authMock *mocks.AuthService) {
				authMock.On("Authorize", "user@example.com", "password123").
					Return(nil,
					nil,
					nil,
					false).
					Once()
			},
			expectedCode: http.StatusBadRequest,
			expectedError: "User does not exists",
		}, 
		// {

		// }
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T){
			t.Parallel()
			authMock := mocks.NewAuthService(t)
			tc.mockSetup(authMock)
			controllerMock := controllers.NewAuthController(authMock)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tc.requestBody))
			rr := httptest.NewRecorder()
			controllerMock.Authorize(rr, req)
			require.Equal(t, tc.expectedCode, rr.Code)

			if tc.expectedError != "" {
				var resp controllers.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Contains(t, resp.Message, tc.expectedError)
			}

		})
	}
}
