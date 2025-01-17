package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseURL string
}

type Server interface {
	Config() *Config
}

type Broker struct {
	config *Config
	router *mux.Router
}

func (b *Broker) Config() *Config {
	return b.config
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("Port is required")
	}

	if config.JWTSecret == "" {
		return nil, errors.New("Secret is required")
	}

	if config.DatabaseURL == "" {
		return nil, errors.New("Database URL is required")
	}

	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
	}

	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) error {

	b.router = mux.NewRouter()
	binder(b, b.router)
	log.Println("Server is running on port %s", b.Config().Port)

	if err := http.ListenAndServe(b.Config().Port, b.router); err != nil {
		log.Fatal("ListenAndServe", err)
	}

	return nil
}
