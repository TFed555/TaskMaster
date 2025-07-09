package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	_ "errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"notes-service/internal/controllers"
	"notes-service/internal/controllers/mocks"
	"notes-service/internal/models"
	"notes-service/internal/pkg/domain_models"
	"notes-service/internal/pkg/responses"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotesController_GetTodos(t *testing.T) {
	cases := []struct {
		name	string
		requestBody string
		mockSetup func(authMock *mocks.NotesService)
		expectedCode int
		expectedError string
	} {
		{
			name: "Success",
			requestBody: ``,
			mockSetup: func(notesMock *mocks.NotesService){
				notesMock.On("GetTodos", url.Values{}).Return([]models.Todo{}, nil).Once()
			},
			expectedCode: http.StatusAccepted,
		},
		{
			name: "Fail",
			requestBody:  ``,
			mockSetup: func(notesMock *mocks.NotesService){
				notesMock.On("GetTodos", url.Values{}).Return(nil, errors.New("can't fetch data")).Once()
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			authMock := mocks.NewNotesService(t)
			tc.mockSetup(authMock)
			controllerMock := controllers.NewNotesController(authMock)
			req := httptest.NewRequest(http.MethodGet, "/todos", bytes.NewBufferString(tc.requestBody))
			rr := httptest.NewRecorder()
			controllerMock.GetTodos(rr, req)
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

func TestNotesController_CreateTodo(t *testing.T) {
	cases := []struct {
		name	string
		requestBody string
		mockSetup func(authMock *mocks.NotesService)
		expectedCode int
		expectedError string
	} {
		{
			name: "Success",
			requestBody: `{
				"title": "aa",
				"priority": "medium",
				"category": "d",
				"description": "f",
				"createdAt": "2024-06-07T00:00:00Z",
				"userId": 1
			}`,
			mockSetup: func(notesMock *mocks.NotesService){
				userId := 1
				notesMock.On("CreateTask", domain_models.Todo{
					Title: "aa",
					Priority: "medium",
					Category: "d",
					Description: "f",
					CreatedAt: "2024-06-07T00:00:00Z",
					UserId: &userId,
				}).Return(1, nil).Once()
			},
			expectedCode: http.StatusAccepted,
		},
		{
			name: "No title",
			requestBody: `{
				"priority": "medium",
				"category": "d",
				"description": "f",
				"createdAt": "2024-06-07T00:00:00Z",
				"userId": 1
			}`,
			mockSetup: func(notesMock *mocks.NotesService){
				userId := 1
				notesMock.On("CreateTask", domain_models.Todo{
					Priority: "medium",
					Category: "d",
					Description: "f",
					CreatedAt: "2024-06-07T00:00:00Z",
					UserId: &userId,
				}).Return(-1, errors.New("Bad request")).Once()
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			authMock := mocks.NewNotesService(t)
			tc.mockSetup(authMock)
			controllerMock := controllers.NewNotesController(authMock)
			req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(tc.requestBody))
			rr := httptest.NewRecorder()
			controllerMock.CreateTodo(rr, req)
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