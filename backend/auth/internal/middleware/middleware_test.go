package middleware_test

import (
	"io"
	"os"
	"testing"

	"github.com/gin-gonic/gin"

	"godiscauth/internal/testutils"
)

type Request struct {
	Email    string
	Password string
}

func TestMain(m *testing.M) {
	testutils.TestEnvSetup()

	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = io.Discard
	os.Exit(m.Run())
}
