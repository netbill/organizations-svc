package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/netbill/logium"
)

type Handlers interface {
	//organization handlers
	CreateOrganization(w http.ResponseWriter, r *http.Request)

	GetOrganization(w http.ResponseWriter, r *http.Request)
	GetOrganizations(w http.ResponseWriter, r *http.Request)
	GetMyOrganizations(w http.ResponseWriter, r *http.Request)

	OpenUpdateOrganizationSession(w http.ResponseWriter, r *http.Request)
	ConfirmUpdateOrganization(w http.ResponseWriter, r *http.Request)
	DeleteUploadOrganizationIcon(w http.ResponseWriter, r *http.Request)
	DeleteUploadOrganizationBanner(w http.ResponseWriter, r *http.Request)
	CancelUpdateOrganization(w http.ResponseWriter, r *http.Request)

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

	//invite handlers
	CreateInvite(w http.ResponseWriter, r *http.Request)
	GetInvite(w http.ResponseWriter, r *http.Request)
	DeleteInvite(w http.ResponseWriter, r *http.Request)
	AcceptInvite(w http.ResponseWriter, r *http.Request)
	DeclineInvite(w http.ResponseWriter, r *http.Request)

	//role handlers
	CreateRole(w http.ResponseWriter, r *http.Request)
	GetRole(w http.ResponseWriter, r *http.Request)
	UpdateRole(w http.ResponseWriter, r *http.Request)
	DeleteRole(w http.ResponseWriter, r *http.Request)

	UpdateRolesRanks(w http.ResponseWriter, r *http.Request)

	UpdateRolePermissions(w http.ResponseWriter, r *http.Request)
	GetAllPermissions(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	AccountAuth(
		allowedRoles ...string,
	) func(next http.Handler) http.Handler
	UpdateOrganization() func(next http.Handler) http.Handler
}

type Router struct {
	handlers    Handlers
	middlewares Middlewares
	log         *logium.Logger
}

func New(
	log *logium.Logger,
	middlewares Middlewares,
	handlers Handlers,
) *Router {
	return &Router{
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

func (rt *Router) Run(ctx context.Context, cfg Config) {
	auth := rt.middlewares.AccountAuth()
	updOrganization := rt.middlewares.UpdateOrganization()

	r := chi.NewRouter()

	r.Route("/organizations-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {

			r.With(auth).Route("/organizations", func(r chi.Router) {
				r.Get("/", rt.handlers.GetOrganizations)
				r.Post("/", rt.handlers.CreateOrganization)

				r.Route("/{organization_id}", func(r chi.Router) {
					r.Get("/", rt.handlers.GetOrganization)
					
					r.With(auth).Route("/update-session", func(r chi.Router) {
						r.Post("/", rt.handlers.OpenUpdateOrganizationSession)
						r.With(updOrganization).Delete("/", rt.handlers.CancelUpdateOrganization)

						r.With(updOrganization).Put("/confirm", rt.handlers.ConfirmUpdateOrganization)
						r.With(updOrganization).Delete("/upload-icon", rt.handlers.DeleteUploadOrganizationIcon)
						r.With(updOrganization).Delete("/upload-banner", rt.handlers.DeleteUploadOrganizationBanner)
					})

					r.Patch("/activate", rt.handlers.ActivateOrganization)
					r.Patch("/deactivate", rt.handlers.DeactivateOrganization)
					r.Get("/members", rt.handlers.GetOrganizationMembers)
					r.Get("/invites", rt.handlers.GetOrganizationInvites)
					r.Route("/roles", func(r chi.Router) {
						r.Get("/", rt.handlers.GetOrganizationRoles)
						r.Put("/ranks", rt.handlers.UpdateRolesRanks)
					})
				})

				r.Get("/me", rt.handlers.GetMyOrganizations)
			})

			r.With(auth).Route("/members", func(r chi.Router) {
				r.Route("/{member_id}", func(r chi.Router) {
					r.Get("/", rt.handlers.GetMember)
					r.Put("/", rt.handlers.UpdateMember)
					r.Delete("/", rt.handlers.DeleteMember)

					r.Route("/roles/{role_id}", func(r chi.Router) {
						r.Post("/", rt.handlers.MemberAddRole)
						r.Delete("/", rt.handlers.MemberRemoveRole)
					})
				})
			})

			r.With(auth).Route("/invites", func(r chi.Router) {
				r.Post("/", rt.handlers.CreateInvite)

				r.Route("/{invite_id}", func(r chi.Router) {
					r.Get("/", rt.handlers.GetInvite)
					r.Patch("/accept", rt.handlers.AcceptInvite)
					r.Patch("/decline", rt.handlers.DeclineInvite)
				})
			})

			r.With(auth).Route("/roles", func(r chi.Router) {
				r.Post("/", rt.handlers.CreateRole)
				r.Get("/permissions", rt.handlers.GetAllPermissions)

				r.Route("/{role_id}", func(r chi.Router) {
					r.Get("/", rt.handlers.GetRole)
					r.Put("/", rt.handlers.UpdateRole)
					r.Delete("/", rt.handlers.DeleteRole)

					r.Put("/permissions", rt.handlers.UpdateRolePermissions)
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

	rt.log.Infof("starting REST service on %s", cfg.Port)

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
		rt.log.Warnf("shutting down REST service...")
	case err := <-errCh:
		if err != nil {
			rt.log.Errorf("REST server error: %v", err)
		}
	}

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		rt.log.Errorf("REST shutdown error: %v", err)
	} else {
		rt.log.Warnf("REST server stopped")
	}
}
