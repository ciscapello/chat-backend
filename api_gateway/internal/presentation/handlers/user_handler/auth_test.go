package userhandler

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ciscapello/api_gateway/internal/common/mocks"
	"github.com/ciscapello/api_gateway/internal/presentation/response"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

type reqBody struct {
	Email string `json:"email"`
}

type testCase struct {
	name                 string
	input                string
	expectedStatus       int
	expectedErrorMessage string
}

func TestAuth(t *testing.T) {
	testCases := []testCase{
		{
			name:                 "Valid email",
			input:                "test@example.com",
			expectedStatus:       http.StatusOK,
			expectedErrorMessage: "",
		},
		{
			name:                 "Invalid email without @",
			input:                "testexample.com",
			expectedStatus:       http.StatusBadRequest,
			expectedErrorMessage: "Invalid email",
		},
		{
			name:                 "Invalid email with invalid domain",
			input:                "test@example",
			expectedStatus:       http.StatusBadRequest,
			expectedErrorMessage: "Invalid email",
		},
		{
			name:                 "Empty email",
			input:                "",
			expectedStatus:       http.StatusBadRequest,
			expectedErrorMessage: "Invalid email",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(reqBody{
				Email: tc.input,
			})

			if err != nil {
				t.Error("cannot marshal")
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockIUserService(ctrl)
			service.EXPECT().Authentication(gomock.Any()).Return(uuid.New(), nil).AnyTimes()
			jwtMan := mocks.NewMockIJwtManager(ctrl)
			responder := response.Responder{}
			logger := slog.New(slog.NewJSONHandler(io.Writer(os.Stdout), nil))
			handler := New(service, logger, jwtMan, responder)

			request := httptest.NewRequest(http.MethodPost, "/users/auth", bytes.NewBuffer(b))
			recorder := httptest.NewRecorder()
			handler.Auth(recorder, request)

			if recorder.Code != tc.expectedStatus {
				t.Errorf("expected status %v; got %v", tc.expectedStatus, recorder.Code)
			}

		})
	}

}
