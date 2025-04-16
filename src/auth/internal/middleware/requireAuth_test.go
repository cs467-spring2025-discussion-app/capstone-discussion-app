package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/matryer/is"

	"godiscauth/internal/middleware"
	"godiscauth/internal/models"
	"godiscauth/internal/repository"
	"godiscauth/internal/testutils"
	"godiscauth/pkg/config"
)

func TestMiddlewareAuth_RequireAuth(t *testing.T) {
	is := is.New(t)

	testDB := testutils.TestDBSetup()
	tx := testDB.Begin()
	defer tx.Rollback()

	// Setup repositories and middleware
	authMw, err := middleware.NewAuthMiddleware(tx)
	is.NoErr(err)
	sessionRepo, err := repository.NewSessionRepository(tx) // Required for session management
	is.NoErr(err)

	router := gin.New()

	// Create test handler and route
	expectedResp := "test handler called"
	testHandler := func(c *gin.Context) {
		c.String(http.StatusOK, expectedResp)
	}

	router.GET("/protected", authMw.RequireAuth(), testHandler)

	// Register a test user
	email := "TestMiddlewareAuth_RequireAuth@test.com"
	user, err := models.NewUser(email, testutils.TestingPassword)
	is.NoErr(err)
	err = tx.Create(user).Error
	is.NoErr(err)

	// Generate a test token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Unix() + config.TokenExpiration,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv(config.JwtCookieName)))
	is.NoErr(err)

	// Create a session record for this token
	session, err := models.NewSession(
		user.ID,
		tokenString,
		time.Now().Add(time.Hour*24),
	)
	is.NoErr(err)

	// Add session to the database
	err = sessionRepo.CreateSession(session)
	is.NoErr(err)

	t.Run("with valid token", func(t *testing.T) {
		// Make a protected request with the token
		reqWithAuth, err := http.NewRequest("GET", "/protected", nil)
		is.NoErr(err)

		reqWithAuth.AddCookie(&http.Cookie{
			Name:  config.JwtCookieName,
			Value: tokenString,
		})
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, reqWithAuth)

		// Request should be OK with a valid token
		is.Equal(http.StatusOK, rr.Code)
		is.Equal(expectedResp, rr.Body.String())
	})

	t.Run("without token", func(t *testing.T) {
		// Make a request without the token
		reqNoAuth, err := http.NewRequest("GET", "/protected", nil)
		is.NoErr(err)

		rrNoAuth := httptest.NewRecorder()
		router.ServeHTTP(rrNoAuth, reqNoAuth)

		// Request should be unauthorized without a token
		is.Equal(http.StatusUnauthorized, rrNoAuth.Code)
	})

	t.Run("with expired token in db", func(t *testing.T) {
		expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Unix() + config.TokenExpiration*2,
		})
		expiredTokenString, _ := expiredToken.SignedString([]byte(os.Getenv(config.JwtCookieName)))

		// Create a session with an expired token
		expiredSession, err := models.NewSession(
			user.ID,
			expiredTokenString,
			time.Now().Add(-1*time.Hour),
		)
		is.NoErr(err)

		err = sessionRepo.CreateSession(expiredSession)
		is.NoErr(err)

		// Make a request with the expired token
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  config.JwtCookieName,
			Value: expiredTokenString,
		})

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Request should be unauthorized with an expired token
		is.Equal(http.StatusUnauthorized, rr.Code)
	})
}

// TestMiddlewareAuth_RequireAuth_SessionRotation tests the session rotation functionality
func TestMiddlewareAuth_RequireAuth_SessionRotation(t *testing.T) {
	is := is.New(t)
	testDB := testutils.TestDBSetup()
	tx := testDB.Begin()
	defer tx.Rollback()

	// Setup repositories and middleware
	authMw, err := middleware.NewAuthMiddleware(tx)
	is.NoErr(err)
	sessionRepo, err := repository.NewSessionRepository(tx)
	is.NoErr(err)

	router := gin.New()

	// Create test handler and route
	expectedResp := "test handler called"
	testHandler := func(c *gin.Context) {
		c.String(http.StatusOK, expectedResp)
	}
	router.GET("/protected", authMw.RequireAuth(), testHandler)

	// Register a test user
	email := "TestMiddlewareAuth_RequireAuth_SessionRotation@test.com"
	user, err := models.NewUser(email, testutils.TestingPassword)
	is.NoErr(err)
	err = tx.Create(user).Error
	is.NoErr(err)

	// Generate a test token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		// Set to ten minutes from now
		"exp": time.Now().Unix() + 600,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv(config.JwtCookieName)))
	is.NoErr(err)

	// Create a session record for this token
	expiresAt := time.Now().Add(10 * time.Minute)
	session, err := models.NewSession(
		user.ID,
		tokenString,
		expiresAt,
	)
	is.NoErr(err)

	// Set the created_at timestamp to 6 minutes ago to simulate a halfway expired session
	session.CreatedAt = time.Now().Add(-6 * time.Minute)

	// Save the session to the database
	err = sessionRepo.CreateSession(session)
	is.NoErr(err)

	// Check that the session is in the database
	updatedSession, err := sessionRepo.GetUnexpiredSessionByToken(tokenString)
	is.NoErr(err)

	// Calculate halfway point
	halfway := updatedSession.CreatedAt.Add(updatedSession.ExpiresAt.Sub(updatedSession.CreatedAt) / 2)
	// Check that the current time is after the halfway point, i.e. the session is halfway expired
	is.True(time.Now().After(halfway))

	t.Run("rotates halfway expired session", func(t *testing.T) {
		// Make a request with the token
		reqWithAuth, err := http.NewRequest("GET", "/protected", nil)
		is.NoErr(err)
		reqWithAuth.AddCookie(&http.Cookie{
			Name:  config.JwtCookieName,
			Value: tokenString,
		})
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, reqWithAuth)

		// Request should be OK with a valid, unexpired token
		is.Equal(http.StatusOK, rr.Code)
		is.Equal(expectedResp, rr.Body.String())

		// Check for new token in cookie
		var newTokenFromCookie string
		for _, cookie := range rr.Result().Cookies() {
			if cookie.Name == config.JwtCookieName {
				newTokenFromCookie = cookie.Value
				break
			}
		}
		is.True(newTokenFromCookie != "")

		// Check that the new token is different from the old one
		is.True(newTokenFromCookie != "")
		is.True(newTokenFromCookie != tokenString)

		// Check that the old session is deleted
		_, err = sessionRepo.GetUnexpiredSessionByToken(tokenString)
		is.True(err != nil)

		// Check that the new, rotated session is created
		newSession, err := sessionRepo.GetUnexpiredSessionByToken(newTokenFromCookie)
		is.NoErr(err)
		is.Equal(user.ID, newSession.UserID)

		// Rotated session should have a new CreatedAt and ExpiresAt
		is.True(newSession.CreatedAt.After(session.CreatedAt))
		is.True(newSession.ExpiresAt.After(session.ExpiresAt))
		// Rotated session should have the same user ID
		is.Equal(user.ID, newSession.UserID)
		// Rotated session should have a different token
		is.True(newSession.Token != tokenString)
	})
}
