package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/matryer/is"

	"godiscauth/internal/handlers"
	"godiscauth/internal/models"
	"godiscauth/internal/repository"
	"godiscauth/internal/server"
	"godiscauth/internal/services"
	"godiscauth/internal/testutils"
	"godiscauth/pkg/apperrors"
)

// TestHandlers_NewUserHandler checks the NewUserHandler constructor returns a valid UsrHandler
func TestHandlers_NewUserHandler(t *testing.T) {
	is := is.New(t)

	t.Run("err on nil user service", func(t *testing.T) {
		uh, err := handlers.NewUserHandler(nil)
		is.Equal(uh, nil)
		is.Equal(err, apperrors.ErrUserServiceIsNil)
	})

	t.Run("creates new user handler", func(t *testing.T) {
		uh := setupUserHandler(t)
		is.True(uh != nil)
	})
}

// TestUserHandle_RegisterUser checks that an http request can add a user to the database
func TestUserHandler_RegisterUser(t *testing.T) {
	is := is.New(t)

	email := "testRegisterUser@test.com"
	password := testutils.TestingPassword // strong password for validator

	t.Run("no email", func(t *testing.T) {
		server := setupServer(t)
		rr, err := makeRequest(
			server.Router,
			"POST",
			"/register",
			UserCredentialsRequest{Password: password},
		)
		is.NoErr(err)
		is.Equal(rr.Code, http.StatusBadRequest)
	})
	t.Run("no password", func(t *testing.T) {
		server := setupServer(t)
		rr, err := makeRequest(
			server.Router,
			"POST",
			"/register",
			UserCredentialsRequest{Email: email},
		)
		is.NoErr(err)
		is.Equal(rr.Code, http.StatusBadRequest)
	})

	t.Run("valid request", func(t *testing.T) {
		server := setupServer(t)
		rr, err := makeRequest(
			server.Router,
			"POST",
			"/register",
			UserCredentialsRequest{Email: email, Password: password},
		)
		is.NoErr(err)
		is.Equal(rr.Code, http.StatusOK)

	// Check user is actually in database
		var user models.User
		result := server.DB.First(&user, "email = ?", email)
		is.NoErr(result.Error)
		is.Equal(user.Email, email)
	})
}

func setupUserHandler(t *testing.T) *handlers.UserHandler {
	t.Helper()

	testDB := testutils.TestDBSetup()
	tx := testDB.Begin()
	t.Cleanup(func() { tx.Rollback() })

	ur, err := repository.NewUserRepository(tx)
	if err != nil {
		t.Fatalf("failed to create user repository: %v", err)
	}
	sr, err := repository.NewSessionRepository(tx)
	if err != nil {
		t.Fatalf("failed to create session repository: %v", err)
	}
	us, err := services.NewUserService(ur, sr)
	if err != nil {
		t.Fatalf("failed to create user service: %v", err)
	}
	uh, err := handlers.NewUserHandler(us)
	if err != nil {
		t.Fatalf("failed to create user handler: %v", err)
	}
	return uh
}

func setupServer(t *testing.T) *server.APIServer {
	testDB := testutils.TestDBSetup()
	tx := testDB.Begin()
	t.Cleanup(func() { tx.Rollback() })

	server, err := server.NewAPIServer(tx)
	if err != nil {
		t.Fatalf("failed to init api server: %v", err)
	}
	server.SetupRoutes()
	return server
}

func makeRequest(router *gin.Engine, method, path string, body any) (*httptest.ResponseRecorder, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		method,
		path,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr, nil
}
