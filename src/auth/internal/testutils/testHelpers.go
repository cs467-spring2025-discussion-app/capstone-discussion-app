package testutils

import (
	"os"

	"github.com/rs/zerolog"
)

// TestEnvSetup sets environment variables for the tests. The tests assume the
// relevant test database has been created. See `scripts/init_testing.sql` to
// create the testing database.
func TestEnvSetup() {
	os.Setenv("PORT", "3001")
	os.Setenv("DB", "host=localhost user=godiscauth_test password=godiscauth_test dbname=godiscauth_test port=5432 sslmode=disable TimeZone=UTC")

	zerolog.SetGlobalLevel(zerolog.Disabled)
}
