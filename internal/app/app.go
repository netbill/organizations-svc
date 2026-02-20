package app

import (
	"github.com/netbill/organizations-svc/internal/config"
	"github.com/netbill/organizations-svc/log"
)

type App struct {
	log    *log.Logger
	config config.Config
}

func New(log *log.Logger, cfg config.Config) *App {
	return &App{
		log:    log,
		config: cfg,
	}
}
