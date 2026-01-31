package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/netbill/logium"
	"github.com/netbill/restkit/tokens/roles"
)

type Handlers interface {
	//Organization handlers
	CreateOrganization(w http.ResponseWriter, r *http.Request)

	GetOrganization(w http.ResponseWriter, r *http.Request)
	GetOrganizations(w http.ResponseWriter, r *http.Request)
	GetMyOrganizations(w http.ResponseWriter, r *http.Request)

	OpenUpdateOrganizationSession(w http.ResponseWriter, r *http.Request)
	UpdateOrganization(w http.ResponseWriter, r *http.Request)
	DeleteUploadOrganizationIcon(w http.ResponseWriter, r *http.Request)
	DeleteUploadOrganizationBanner(w http.ResponseWriter, r *http.Request)

	ActivateOrganization(w http.ResponseWriter, r *http.Request)
	DeactivateOrganization(w http.ResponseWriter, r *http.Request)

	GetOrganizationInvites(w http.ResponseWriter, r *http.Request)
	GetOrganizationMembers(w http.ResponseWriter, r *http.Request)
	GetOrganizationRoles(w http.ResponseWriter, r *http.Request)

	//OrganizationMember handlers
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
	AccountAuth() func(http.Handler) http.Handler
	AccountRoleGrant(allowedRoles map[string]bool) func(http.Handler) http.Handler
	UpdateOrganization() func(http.Handler) http.Handler
}

type Service struct {
	handlers    Handlers
	middlewares Middlewares
	log         *logium.Logger
}

func New(
	log *logium.Logger,
	middlewares Middlewares,
	handlers Handlers,
) *Service {
	return &Service{
		log:         log,
		middlewares: middlewares,
		handlers:    handlers,
	}
}

type Config struct {
	Port              string
	TimeoutRead       time.Duration
	TimeoutReadHeader time.Duration
	TimeoutWrite      time.Duration
	TimeoutIdle       time.Duration
}

func (s *Service) Run(ctx context.Context, cfg Config) {
	auth := s.middlewares.AccountAuth()
	sysadmin := s.middlewares.AccountRoleGrant(map[string]bool{
		roles.SystemAdmin: true,
	})
	updOrganization := s.middlewares.UpdateOrganization()

	r := chi.NewRouter()

	r.Route("/organizations-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {

			r.With(auth).Route("/organizations", func(r chi.Router) {
				r.Get("/", s.handlers.GetOrganizations)
				r.Post("/", s.handlers.CreateOrganization)

				r.Route("/{organization_id}", func(r chi.Router) {
					r.Get("/", s.handlers.GetOrganization)
					r.With(auth).Route("/update-session", func(r chi.Router) {
						r.Post("/", s.handlers.OpenUpdateOrganizationSession)

						r.With(updOrganization).Put("/confirm", s.handlers.UpdateOrganization)
						r.With(updOrganization).Delete("/upload-icon", s.handlers.DeleteUploadOrganizationIcon)
						r.With(updOrganization).Delete("/upload-banner", s.handlers.DeleteUploadOrganizationBanner)
					})

					r.Patch("/activate", s.handlers.ActivateOrganization)
					r.Patch("/deactivate", s.handlers.DeactivateOrganization)
					r.Get("/members", s.handlers.GetOrganizationMembers)
					r.Get("/invites", s.handlers.GetOrganizationInvites)
					r.Route("/roles", func(r chi.Router) {
						r.Get("/", s.handlers.GetOrganizationRoles)
						r.Put("/ranks", s.handlers.UpdateRolesRanks)
					})
				})

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
				r.Route("/organizations", func(r chi.Router) {
				})
			})
		})
	})

	srv := &http.Server{
		Addr:              cfg.Port,
		Handler:           r,
		ReadTimeout:       cfg.TimeoutRead,
		ReadHeaderTimeout: cfg.TimeoutReadHeader,
		WriteTimeout:      cfg.TimeoutWrite,
		IdleTimeout:       cfg.TimeoutIdle,
	}

	s.log.Infof("starting REST service on %s", cfg.Port)

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
		s.log.Warnf("shutting down REST service...")
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
		s.log.Warnf("REST server stopped")
	}
}
