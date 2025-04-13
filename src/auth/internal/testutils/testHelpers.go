package testutils

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestEnvSetup sets environment variables for the tests. The tests assume the
// relevant test database has been created. See `scripts/init_testing.sql` to
// create the testing database.
func TestEnvSetup() {
	os.Setenv("PORT", "3001")
	os.Setenv("DB", "host=localhost user=godiscauth_test password=godiscauth_test dbname=godiscauth_test port=5432 sslmode=disable TimeZone=UTC")

	zerolog.SetGlobalLevel(zerolog.Disabled)
}

// TestDBSetup sets up a test database connection.
func TestDBSetup() *gorm.DB {
	// Silence GORM logs for testing
	gormLogger := logger.New(
		log.New(io.Discard, "", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Connect to test DB
	db, err := gorm.Open(postgres.Open(os.Getenv("DB")), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	return db
}
