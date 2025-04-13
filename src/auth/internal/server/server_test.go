package server_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/matryer/is"

	"godiscauth/internal/server"
	"godiscauth/internal/testutils"
)

// TestMain sets up the test environment for all tests in the `server_test` package.
func TestMain(m *testing.M) {
	testutils.TestEnvSetup()
	os.Exit(m.Run())
}

func TestPingRoute(t *testing.T) {
	is := is.New(t)

	testutils.TestEnvSetup()

	testDB := testutils.TestDBSetup()
	server := server.NewAPIServer(testDB)
	server.SetupRoutes()

	req, _ := http.NewRequest("GET", "/ping", nil)
	rr := httptest.NewRecorder()

	server.Router.ServeHTTP(rr, req)

	is.Equal(http.StatusOK, rr.Code)
	is.Equal("pong", rr.Body.String())
}
