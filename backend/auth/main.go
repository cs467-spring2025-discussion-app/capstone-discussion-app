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
	"gorm.io/gorm"

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
	checkSessionKey()

	logger.SetupLogger()

	db := connectDB()

	startAPIServer(db)

	quit := makeQuitListener()

	wg, cancel := startJobs(db)

	// Block until quit signal
	<-quit
	log.Info().Msg("Server is shutting down...")

	// Signal jobs to terminate
	cancel()

	waitForJobs(wg)

	log.Info().Msg("Server exited")
}

// Ensure session key is complex enough for encryption
func checkSessionKey() {
	if err := passwordvalidator.Validate(os.Getenv(config.SessionKey), config.MinEntropyBits); err != nil {
		log.Fatal().Err(err).Msg("Session secret is not complex enough")
	}
}

// Connect and migrate DB
func connectDB() *gorm.DB {
	db, err := database.NewDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
	}
	err = database.Migrate(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Error migrating database")
	}
	log.Info().Msg(fmt.Sprintf("Connected to postgres database"))
	return db
}

// Start API Server
func startAPIServer(db *gorm.DB) {
	log.Info().
		Str(config.AuthServerPort, os.Getenv(config.AuthServerPort)).
		Msg("Starting server")
	apiServer, err := server.NewAPIServer(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing server")
	}
	go apiServer.Run()
}

func startJobs(db *gorm.DB) (*sync.WaitGroup, context.CancelFunc) {
	// Start background jobs
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	go jobs.StartJobs(ctx, &wg, db)
	return &wg, cancel
}

func makeQuitListener() chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	return quit
}

func waitForJobs(wg *sync.WaitGroup) {
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
}
