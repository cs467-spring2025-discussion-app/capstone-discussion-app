package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	passwordvalidator "github.com/wagslane/go-password-validator"

	"godiscauth/internal/database"
	"godiscauth/internal/server"
	"godiscauth/pkg/config"
	"godiscauth/pkg/logger"
)

// main is the entry point for the auth service. It sets up the logger, connects to the database, and starts the API server.
func main() {
	// Ensure session key must be complex for encryption
	if err := passwordvalidator.Validate(os.Getenv(config.SessionKey), config.MinEntropyBits); err != nil {
		log.Fatal().Err(err).Msg("Session secret is not complex enough")
	}

	logger.SetupLogger()
	db, err := database.NewDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
	}
	err = database.Migrate(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Error migrating database")
	}

	log.Info().Msg(fmt.Sprintf("Connected to postgres database"))
	log.Info().Str(config.AuthServerPort, os.Getenv(config.AuthServerPort)).Msg("Starting server")
	apiServer, err := server.NewAPIServer(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing server")
	}
	apiServer.Run()
}
