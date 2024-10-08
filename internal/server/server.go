package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"

	"arczed/internal/server/configs"
	"arczed/internal/server/database"
)

type Server struct {
	port   int
	db     database.Service
	config configs.Config
}

func NewServer() *http.Server {

	var config configs.Config
	// Load environment variables into the config struct
	if err := env.Parse(&config); err != nil {
		log.Fatalf("failed to parse env vars: %v", err)
	}

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:   port,
		db:     database.New(&config),
		config: config,
	}
	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
