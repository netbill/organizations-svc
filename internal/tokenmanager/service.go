package tokenmanager

import (
	"time"
)

const Issuer = "organizations-svc"

type Config struct {
	AccountAccess struct {
		SecretKey string `mapstructure:"secret_key" reqquire:"true"`
	} `mapstructure:"account_access" reqquire:"true"`
	Media struct {
		Token struct {
			SecretKey string        `mapstructure:"secret_key" reqquire:"true"`
			TTL       time.Duration `mapstructure:"ttl" reqquire:"true"`
		} `mapstructure:"token" reqquire:"true"`
	} `mapstructure:"media" reqquire:"true"`
}

type Manager struct {
	issuer string

	accessSK string
	uploadSK string

	mediaTTL time.Duration
}

func New(config Config) *Manager {
	return &Manager{
		issuer:   Issuer,
		accessSK: config.AccountAccess.SecretKey,
		uploadSK: config.Media.Token.SecretKey,
		mediaTTL: config.Media.Token.TTL,
	}
}
