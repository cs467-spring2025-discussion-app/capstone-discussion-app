package database_test

import (
	"testing"

	"github.com/matryer/is"

	"godiscauth/internal/database"
)

func TestConnectToDB(t *testing.T) {
	is := is.New(t)
	testDB := database.NewDB()
	t.Run("connects", func(t *testing.T) {
		is.True(testDB != nil)
	})
}
