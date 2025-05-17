package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	passwordvalidator "github.com/wagslane/go-password-validator"

	_ "godiscauth/docs"
	"godiscauth/internal/database"
	"godiscauth/internal/jobs"
	"godiscauth/internal/server"
	"godiscauth/pkg/config"
	"godiscauth/pkg/logger"
)

// main is the entry point for the auth service. It sets up the logger,
// connects to the database, starts the API server, and start any background jobs
func main() {
	// Ensure session key must be complex for encryption
	if err := passwordvalidator.Validate(os.Getenv(config.SessionKey), config.MinEntropyBits); err != nil {
		log.Fatal().Err(err).Msg("Session secret is not complex enough")
	}

	logger.SetupLogger()

	// Connect and migrate
	db, err := database.NewDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
	}
	err = database.Migrate(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Error migrating database")
	}
	log.Info().Msg(fmt.Sprintf("Connected to postgres database"))

	// Start API server
	log.Info().
		Str(config.AuthServerPort, os.Getenv(config.AuthServerPort)).
		Msg("Starting server")
	apiServer, err := server.NewAPIServer(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing server")
	}
	go apiServer.Run()

	// Create channel to listen for shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start background jobs
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	go jobs.StartJobs(ctx, &wg, db)

	// Block until quit signal
	<-quit
	log.Info().Msg("Server is shutting down...")

	// Signal jobs to terminate
	cancel()

	// Wait with timeout for jobs to finish
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Info().Msg("graceful shutdown")
	case <-time.After(10 * time.Second):
		log.Info().Msg("forcequitting")
	}

	log.Info().Msg("Server exited")
}
