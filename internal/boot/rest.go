package boot

import (
	"github.com/netbill/organizations-svc/internal/rest"
	"github.com/netbill/organizations-svc/internal/rest/controller"
	"github.com/netbill/organizations-svc/internal/rest/middlewares"
	"github.com/netbill/organizations-svc/internal/tokenmanager"
	"github.com/netbill/restkit"
)



func (c *Config) Rest(modules controller.Modules) *rest.Server {
	token := tokenmanager.New(tokenmanager.Config{
		c.Service.
	})

	responser := restkit.NewResponser()
	ctrl := controller.New(modules, responser)
	mdll := middlewares.New(responser, tokenManager)

	return rest.New(mdll, ctrl)
}
