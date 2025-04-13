package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"godiscauth/internal/database"
	"godiscauth/pkg/logger"
)

// main is the entry point for the auth service and sets up the logger, connects to the database,
// and starts the server.
// The service expects the following environment variables:
// - DB: The connection string to the database.
// - PORT: The port on which the service will run.
func main() {
	logger.SetupLogger()
	_, err := database.NewDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
	}

	log.Info().Msg(fmt.Sprintf("Connected to database: %s", os.Getenv("DB")))
	log.Info().Str("PORT", os.Getenv("PORT")).Msg("Starting server")
}
