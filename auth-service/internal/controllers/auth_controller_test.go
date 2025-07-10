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
	"auth-service/internal/pkg/domain_models"
	"auth-service/internal/pkg/responses"
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
				"name": "testuser",
				"email": "user@example.com",
				"password": "password123"
			}`,
			mockSetup: func(authMock *mocks.AuthService) {
				authMock.On("Register", domain_models.RegisterParams{Name:"testuser", Email:"user@example.com", Password:"password123"}).
					Return(1, domain_models.Tokens{
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
				"name": "testuser",
				"email": "user@example.com"
			}`,
			mockSetup:    func(authMock *mocks.AuthService) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Constraint email",
			requestBody: `{
				"name": "testuser",
				"email": "user@example.com",
				"password": "password123"
			}`,
			mockSetup: func(authMock *mocks.AuthService) {
				authMock.On("Register", domain_models.RegisterParams{Name:"testuser", Email:"user@example.com", Password:"password123"}).
					Return(1, domain_models.Tokens{}, errors.New("some error")).
					Once()
			},
			expectedCode:  http.StatusInternalServerError,
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
				var resp responses.ErrorResponse
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
				authMock.On("Authorize", domain_models.AuthorizeParams{Email: "user@example.com", Password: "password123"}).
					Return(models.User{ID: 1, Email: "user@example.com", Login: "user"}, domain_models.Tokens{
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
				authMock.On("Authorize", domain_models.AuthorizeParams{Email: "user@example.com", Password: "password123"}).
					Return(models.User{},
					domain_models.Tokens{},
					errors.New("User not found"),
					false).
					Once()
			},
			expectedCode: http.StatusNotFound,
			expectedError: "User not found",
		}, 
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
				var resp responses.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Contains(t, resp.Message, tc.expectedError)
			}

		})
	}
}

func TestAuthController_Refresh(t *testing.T) {
	cases := [] struct {
		name string
		cookie string
		mockSetup func (authMock *mocks.AuthService)
		expectedCode int
		expectedError string
		expectedTokens bool
	} {
		{
			name: "Success",
			cookie: "refresh_token=valid_refresh_token",
			mockSetup: func(authMock *mocks.AuthService) {
				authMock.On("ValidateToken", "valid_refresh_token").
					Return(true, "").
					Once()

				authMock.On("Refresh", "valid_refresh_token").
					Return("new_access_token", "new_refresh_token", nil).
					Once()
			},
			expectedCode: http.StatusAccepted,
			expectedTokens: true,
		},
		{
			name:   "No token",
			cookie: "",
			mockSetup: func(authMock *mocks.AuthService) {
				authMock.On("ValidateToken", "").
					Return(false, "Access denied").
					Once()
			},
			expectedCode:  http.StatusForbidden,
			expectedError: "Access denied",
		},
		{
			name:   "Invalid token",
			cookie: "refresh_token=invalid_token",
			mockSetup: func(authMock *mocks.AuthService) {
				authMock.On("ValidateToken", "invalid_token").
					Return(false, "Access denied").
					Once()
			},
			expectedCode:  http.StatusForbidden,
			expectedError: "Access denied",
		},
		{
			name:   "Refresh error",
			cookie: "refresh_token=expired_token",
			mockSetup: func(authMock *mocks.AuthService) {
				authMock.On("ValidateToken", "expired_token").
					Return(true, "").
					Once()
				authMock.On("Refresh", "expired_token").
					Return("", "", errors.New("token expired")).
					Once()
			},
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T){
			t.Parallel()
			authMock := mocks.NewAuthService(t)
			tc.mockSetup(authMock)
			controllerMock := controllers.NewAuthController(authMock)
			req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
			if tc.cookie != "" {
				req.Header.Set("Cookie", tc.cookie)
			}

			rr := httptest.NewRecorder()
			controllerMock.Refresh(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code)

			if tc.expectedError != "" {
				require.Contains(t, rr.Body.String(), tc.expectedError)
			}
			
			if tc.expectedTokens {
				var resp responses.RefreshResponse
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Equal(t, "new_access_token", resp.AccessToken)
				require.Equal(t, "new_refresh_token", resp.RefreshToken)
				cookies := rr.Result().Cookies()
				require.Len(t, cookies, 2)
				var accessCookie, refreshCookie *http.Cookie
				for _, cookie := range cookies {
					if cookie.Name == "access_token" {
						accessCookie = cookie
					} else if cookie.Name == "refresh_token" {
						refreshCookie = cookie
					}
				}
				require.NotNil(t, accessCookie)
				require.NotNil(t, refreshCookie)
				require.Equal(t, "new_access_token", accessCookie.Value)
				require.Equal(t, "new_refresh_token", refreshCookie.Value)
			}
		})
	}
}