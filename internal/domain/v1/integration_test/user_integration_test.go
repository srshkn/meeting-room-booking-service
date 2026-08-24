package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	v1GenAPI "mrb-service/internal/generated/v1"
)

func TestRegisterUserIntegration(t *testing.T) {

	httpHandler := testServerApp.Handler()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/user/register",
		strings.NewReader(`{
			"username": "Ben",
			"email": "benben@example.com",
			"password": "secret321",
			"confirmation": "secret321"
		}`),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	httpHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"status = %d, want %d; body = %s",
			rec.Code,
			http.StatusCreated,
			rec.Body.String(),
		)
	}

	var response v1GenAPI.PostUserRegister201JSONResponse

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.User == nil {
		t.Fatal("response user is nil")
	}

	userResponse := response.User

	if userResponse.Username != "Ben" {
		t.Errorf(
			"username = %q, want %q",
			userResponse.Username,
			"Ben",
		)
	}

	ctx := t.Context()

	user, err := testDB.GetUserByID(ctx, userResponse.Id)
	if err != nil {
		t.Fatalf("query user: %v", err)
	}

	if userResponse.Id == uuid.Nil {
		t.Fatal("response id is empty")
	}

	if user.Email != "benben@example.com" {
		t.Errorf(
			"email = %q, want %q",
			user.Email,
			"benben@example.com",
		)
	}

	if user.Username != "Ben" {
		t.Errorf(
			"database username = %q, want %q",
			user.Username,
			"Ben",
		)
	}

}
