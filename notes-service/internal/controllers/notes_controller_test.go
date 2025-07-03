package controllers_test

import (
	"bytes"
	"encoding/json"
	_ "errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"notes-service/internal/controllers"
	"notes-service/internal/controllers/mocks"
	"notes-service/internal/models"
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
		
	}

	for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			authMock := mocks.NewNotesService(t)
			tc.mockSetup(authMock)
			controllerMock := controllers.NewNotesController(authMock)
			req := httptest.NewRequest(http.MethodPost, "/registration", bytes.NewBufferString(tc.requestBody))
			rr := httptest.NewRecorder()
			controllerMock.GetTodos(rr, req)
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