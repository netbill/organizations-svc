package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/netbill/organizations-svc/pkg/log"
	"github.com/netbill/restkit/tokens"
)

type Handlers interface {
	//organization handlers
	CreateOrganization(w http.ResponseWriter, r *http.Request)

	GetOrganization(w http.ResponseWriter, r *http.Request)
	GetOrganizations(w http.ResponseWriter, r *http.Request)
	GetMyOrganizations(w http.ResponseWriter, r *http.Request)

	CreateOrganizationUploadMediaLink(w http.ResponseWriter, r *http.Request)

	UpdateOrganization(w http.ResponseWriter, r *http.Request)

	DeleteOrganization(w http.ResponseWriter, r *http.Request)

	DeleteOrganizationUploadIcon(w http.ResponseWriter, r *http.Request)
	DeleteOrganizationUploadBanner(w http.ResponseWriter, r *http.Request)

	ActivateOrganization(w http.ResponseWriter, r *http.Request)
	DeactivateOrganization(w http.ResponseWriter, r *http.Request)
	SuspendOrganization(w http.ResponseWriter, r *http.Request)
	UnsuspendOrganization(w http.ResponseWriter, r *http.Request)

	GetOrganizationInvites(w http.ResponseWriter, r *http.Request)
	GetOrganizationMembers(w http.ResponseWriter, r *http.Request)

	//OrganizationMember handlers
	GetMember(w http.ResponseWriter, r *http.Request)
	UpdateMember(w http.ResponseWriter, r *http.Request)
	DeleteMember(w http.ResponseWriter, r *http.Request)

	LeaveOrganization(w http.ResponseWriter, r *http.Request)

	//invite handlers
	CreateInvite(w http.ResponseWriter, r *http.Request)
	GetInvite(w http.ResponseWriter, r *http.Request)

	CancelledInvite(w http.ResponseWriter, r *http.Request)
	AcceptInvite(w http.ResponseWriter, r *http.Request)
	DeclineInvite(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	AccountAuth(
		allowedRoles ...string,
	) func(next http.Handler) http.Handler
	Logger(log *log.Logger) func(next http.Handler) http.Handler
	CorsDocs() func(next http.Handler) http.Handler
}

type Server struct {
	handlers    Handlers
	middlewares Middlewares
}

func New(
	middlewares Middlewares,
	handlers Handlers,
) *Server {
	return &Server{
		middlewares: middlewares,
		handlers:    handlers,
	}
}

type Config struct {
	Port              int
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func (s *Server) Run(ctx context.Context, log *log.Logger, cfg Config) {
	auth := s.middlewares.AccountAuth()
	sysadmin := s.middlewares.AccountAuth(tokens.RoleSystemAdmin)

	r := chi.NewRouter()
	r.Use(
		s.middlewares.Logger(log),
		s.middlewares.CorsDocs(),
	)

	r.Route("/organizations-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {

			r.Route("/organizations", func(r chi.Router) {
				r.Get("/", s.handlers.GetOrganizations)
				r.With(auth).Post("/", s.handlers.CreateOrganization)

				r.Route("/{organization_id}", func(r chi.Router) {
					r.Get("/", s.handlers.GetOrganization)
					r.With(auth).Put("/", s.handlers.UpdateOrganization)
					r.With(auth).Delete("/", s.handlers.DeleteOrganization)

					r.With(auth).Route("/media", func(r chi.Router) {
						r.Route("/upload", func(r chi.Router) {
							r.Post("/url", s.handlers.CreateOrganizationUploadMediaLink)

							r.Delete("/icon", s.handlers.DeleteOrganizationUploadIcon)
							r.Delete("/banner", s.handlers.DeleteOrganizationUploadBanner)
						})
					})

					r.With(auth).Post("/activate", s.handlers.ActivateOrganization)
					r.With(auth).Post("/deactivate", s.handlers.DeactivateOrganization)

					r.With(sysadmin).Post("/suspend", s.handlers.SuspendOrganization)
					r.With(sysadmin).Post("/unsuspend", s.handlers.UnsuspendOrganization)

					r.Get("/members", s.handlers.GetOrganizationMembers)
					r.With(auth).Get("/invites", s.handlers.GetOrganizationInvites)
					
					r.With(auth).Delete("/leave", s.handlers.LeaveOrganization)
				})

				r.With(auth).Get("/me", s.handlers.GetMyOrganizations)
			})

			r.Route("/members", func(r chi.Router) {
				r.Route("/{member_id}", func(r chi.Router) {
					r.Get("/", s.handlers.GetMember)
					r.With(auth).Put("/", s.handlers.UpdateMember)
					r.With(auth).Delete("/", s.handlers.DeleteMember)
				})

				//r.Get("/me", s.handlers.GetMyMembers)
			})

			r.With(auth).Route("/invites", func(r chi.Router) {
				r.Post("/", s.handlers.CreateInvite)

				r.Route("/{invite_id}", func(r chi.Router) {
					r.Get("/", s.handlers.GetInvite)
					r.Patch("/accept", s.handlers.AcceptInvite)
					r.Patch("/decline", s.handlers.DeclineInvite)
					r.Patch("/canclled", s.handlers.CancelledInvite)
				})

				//r.Get("/me", s.handlers.GetMyInvites)
			})
		})
	})

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           r,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	log.WithField("port", cfg.Port).Info("starting http service...")

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
		log.Info("shutting down http service...")
	case err := <-errCh:
		if err != nil {
			log.WithError(err).Error("http server error")
		}
	}

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		log.WithError(err).Error("failed to shutdown http server gracefully")
	} else {
		log.Info("http server shutdown gracefully")
	}
}
