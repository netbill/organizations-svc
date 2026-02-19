package tokenmanager

type Config struct {
	Issuer   string
	AccessSK string
}

type Manager struct {
	issuer string

	accessSK string
}

func New(config Config) *Manager {
	return &Manager{
		issuer:   config.Issuer,
		accessSK: config.AccessSK,
	}
}
