package app

import (
	"github.com/netbill/logium"
	"github.com/netbill/pgdbx"
)

type App struct {
	db     *pgdbx.DB
	log    *logium.Logger
	config *Config
}
