package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"godiscauth/internal/database"
	"godiscauth/pkg/logger"
)

func main() {
	logger.SetupLogger()
	_, err := database.NewDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
	}

	log.Info().Msg(fmt.Sprintf("Connected to database: %s", os.Getenv("DB")))
	log.Info().Str("PORT", os.Getenv("PORT")).Msg("Starting server")
}
