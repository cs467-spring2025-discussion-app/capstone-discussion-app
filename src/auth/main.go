package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"godiscauth/internal/database"
	"godiscauth/internal/server"
	"godiscauth/pkg/logger"
)

// main is the entry point for the auth service. It sets up the logger, connects to the database, and starts the API server.
func main() {
	logger.SetupLogger()
	db, err := database.NewDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
	}
	err = database.Migrate(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Error migrating database")
	}

	log.Info().Msg(fmt.Sprintf("Connected to database: %s", os.Getenv("DB")))
	log.Info().Str("PORT", os.Getenv("PORT")).Msg("Starting server")
	apiServer, err := server.NewAPIServer(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing server")
	}
	apiServer.Run()
}
