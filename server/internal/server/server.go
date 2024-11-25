package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"server/internal/database"
	"server/internal/logger"
	"server/internal/mailer"
)

type Server struct {
	port int

	db     database.Service
	log    *logger.Logger
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,

		db:     database.New(),
		log:    logger.New(),
		mailer: mailer.New(),
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
