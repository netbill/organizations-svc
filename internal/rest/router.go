package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/netbill/logium"
	"github.com/netbill/organizations-svc/internal"
	"github.com/netbill/restkit/roles"
)

type Handlers interface {
	//Organization handlers
	CreateOrganization(w http.ResponseWriter, r *http.Request)

	GetOrganization(w http.ResponseWriter, r *http.Request)
	GetOrganizations(w http.ResponseWriter, r *http.Request)
	GetMyOrganizations(w http.ResponseWriter, r *http.Request)

	UpdateOrganization(w http.ResponseWriter, r *http.Request)

	SuspendOrganization(w http.ResponseWriter, r *http.Request)
	ActivateOrganization(w http.ResponseWriter, r *http.Request)
	DeactivateOrganization(w http.ResponseWriter, r *http.Request)

	GetOrganizationInvites(w http.ResponseWriter, r *http.Request)
	GetOrganizationMembers(w http.ResponseWriter, r *http.Request)
	GetOrganizationRoles(w http.ResponseWriter, r *http.Request)

	//Member handlers
	GetMember(w http.ResponseWriter, r *http.Request)
	UpdateMember(w http.ResponseWriter, r *http.Request)
	DeleteMember(w http.ResponseWriter, r *http.Request)

	MemberAddRole(w http.ResponseWriter, r *http.Request)
	MemberRemoveRole(w http.ResponseWriter, r *http.Request)

	//Invite handlers
	CreateInvite(w http.ResponseWriter, r *http.Request)
	GetInvite(w http.ResponseWriter, r *http.Request)
	DeleteInvite(w http.ResponseWriter, r *http.Request)
	AcceptInvite(w http.ResponseWriter, r *http.Request)
	DeclineInvite(w http.ResponseWriter, r *http.Request)

	//Role handlers
	CreateRole(w http.ResponseWriter, r *http.Request)
	GetRole(w http.ResponseWriter, r *http.Request)
	UpdateRole(w http.ResponseWriter, r *http.Request)
	DeleteRole(w http.ResponseWriter, r *http.Request)

	UpdateRolesRanks(w http.ResponseWriter, r *http.Request)

	UpdateRolePermissions(w http.ResponseWriter, r *http.Request)
	GetAllPermissions(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	Auth() func(http.Handler) http.Handler
	RoleGrant(allowedRoles map[string]bool) func(http.Handler) http.Handler
}

type Service struct {
	handlers    Handlers
	middlewares Middlewares
	log         logium.Logger
}

func New(
	log logium.Logger,
	middlewares Middlewares,
	handlers Handlers,
) *Service {
	return &Service{
		log:         log,
		middlewares: middlewares,
		handlers:    handlers,
	}
}

func (s *Service) Run(ctx context.Context, cfg internal.Config) {
	auth := s.middlewares.Auth()
	sysadmin := s.middlewares.RoleGrant(map[string]bool{
		roles.SystemAdmin: true,
	})

	r := chi.NewRouter()

	r.Route("/organizations-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {

			r.With(auth).Route("/organizations", func(r chi.Router) {
				r.Get("/", s.handlers.GetOrganizations)
				r.Post("/", s.handlers.CreateOrganization)

				r.Route("/{organization_id}", func(r chi.Router) {
					r.Get("/", s.handlers.GetOrganization)
					r.Put("/", s.handlers.UpdateOrganization)

					r.Patch("/activate", s.handlers.ActivateOrganization)
					r.Patch("/deactivate", s.handlers.DeactivateOrganization)

					r.Get("/members", s.handlers.GetOrganizationMembers)
					r.Get("/invites", s.handlers.GetOrganizationInvites)
					r.Route("/roles", func(r chi.Router) {
						r.Get("/", s.handlers.GetOrganizationRoles)
						r.Put("/ranks", s.handlers.UpdateRolesRanks)
					})
				})

				//TODO
				r.Get("/me", s.handlers.GetMyOrganizations)
			})

			r.With(auth).Route("/members", func(r chi.Router) {
				r.Route("/{member_id}", func(r chi.Router) {
					r.Get("/", s.handlers.GetMember)
					r.Put("/", s.handlers.UpdateMember)
					r.Delete("/", s.handlers.DeleteMember)

					r.Route("/roles/{role_id}", func(r chi.Router) {
						r.Post("/", s.handlers.MemberAddRole)
						r.Delete("/", s.handlers.MemberRemoveRole)
					})
				})
			})

			r.With(auth).Route("/invites", func(r chi.Router) {
				r.Post("/", s.handlers.CreateInvite)

				r.Route("/{invite_id}", func(r chi.Router) {
					r.Get("/", s.handlers.GetInvite)
					r.Patch("/accept", s.handlers.AcceptInvite)
					r.Patch("/decline", s.handlers.DeclineInvite)
				})
			})

			r.With(auth).Route("/roles", func(r chi.Router) {
				r.Post("/", s.handlers.CreateRole)
				r.Get("/permissions", s.handlers.GetAllPermissions)

				r.Route("/{role_id}", func(r chi.Router) {
					r.Get("/", s.handlers.GetRole)
					r.Put("/", s.handlers.UpdateRole)
					r.Delete("/", s.handlers.DeleteRole)

					r.Put("/permissions", s.handlers.UpdateRolePermissions)
				})
			})

			r.With(auth, sysadmin).Route("/admin", func(r chi.Router) {
				r.With(auth, sysadmin).Route("/organizations", func(r chi.Router) {
					r.Patch("/", s.handlers.SuspendOrganization)
				})
			})
		})
	})

	srv := &http.Server{
		Handler:           r,
		Addr:              cfg.Rest.Port,
		ReadTimeout:       cfg.Rest.Timeouts.Read,
		ReadHeaderTimeout: cfg.Rest.Timeouts.ReadHeader,
		WriteTimeout:      cfg.Rest.Timeouts.Write,
		IdleTimeout:       cfg.Rest.Timeouts.Idle,
	}

	s.log.Infof("starting REST service on %s", cfg.Rest.Port)

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		} else {
			errCh <- nil
		}
	}()

	select {
	case <-ctx.Done():
		s.log.Info("shutting down REST service...")
	case err := <-errCh:
		if err != nil {
			s.log.Errorf("REST server error: %v", err)
		}
	}

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		s.log.Errorf("REST shutdown error: %v", err)
	} else {
		s.log.Info("REST server stopped")
	}
}
