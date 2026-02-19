package app

import (
	"context"

	"github.com/netbill/organizations-svc/internal/rest"
	"github.com/netbill/organizations-svc/internal/rest/controller"
	"github.com/netbill/organizations-svc/internal/rest/middlewares"
	"github.com/netbill/organizations-svc/internal/tokenmanager"
	"github.com/netbill/restkit"
)

func (a *App) RunRest(ctx context.Context, modules *controller.Modules) {
	tokenManager := tokenmanager.New(tokenmanager.Config{
		Issuer:   a.config.Auth.Tokens.Issuer,
		AccessSK: a.config.Auth.Tokens.AccountAccess.SecretKey,
	})

	responser := restkit.NewResponser()
	ctrl := controller.New(modules, responser)
	mdll := middlewares.New(responser, tokenManager)
	router := rest.New(mdll, ctrl)

	router.Run(ctx, nil, rest.Config{
		Port:              a.config.Rest.Port,
		ReadTimeout:       a.config.Rest.Timeouts.Read,
		ReadHeaderTimeout: a.config.Rest.Timeouts.ReadHeader,
		WriteTimeout:      a.config.Rest.Timeouts.Write,
		IdleTimeout:       a.config.Rest.Timeouts.Idle,
	})
}
