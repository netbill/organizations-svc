package tokenmanager

import (
	"fmt"

	"github.com/netbill/restkit/tokens"
)

func (m *Manager) ParseAccountAuthAccessClaims(token string) (tokens.AccountAuthClaims, error) {
	data, err := tokens.ParseAccountJWT(token, m.accessSK)
	if err != nil {
		return tokens.AccountAuthClaims{}, fmt.Errorf("failed to parse access token, cause: %w", err)
	}

	return data, nil
}
