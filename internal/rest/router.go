package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/umisto/cities-svc/internal"
	"github.com/umisto/logium"
	"github.com/umisto/restkit/roles"
)

type Handlers interface {
	CreateAgglomeration(w http.ResponseWriter, r *http.Request)
	UpdateAgglomeration(w http.ResponseWriter, r *http.Request)
	GetAgglomeration(w http.ResponseWriter, r *http.Request)
	ActivateAgglomeration(w http.ResponseWriter, r *http.Request)
	DeactivateAgglomeration(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	Auth() func(http.Handler) http.Handler
	SystemRoleGrant(allowedRoles map[string]bool) func(http.Handler) http.Handler
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
	sysadmin := s.middlewares.SystemRoleGrant(map[string]bool{
		roles.SystemAdmin: true,
	})

	r := chi.NewRouter()

	r.Route("/cities-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/{city_slug}", nil)

			r.Route("/agglomerations", func(r chi.Router) {
				r.Get("/", nil)
				r.With(sysadmin).Post("/", nil)

				r.With(auth).Route("/{agglomeration_id}", func(r chi.Router) {
					r.Get("/", nil)
					r.Put("/", nil)

					r.Patch("/activate", nil)
					r.Patch("/deactivate", nil)

					r.Get("/cities", nil)
					r.Get("/members", nil)
					r.Get("/roles", nil)
				})
			})

			r.Route("/cities", func(r chi.Router) {
				r.With(auth, sysadmin).Post("/", nil)

				r.Route("/{city_id}", func(r chi.Router) {
					r.Get("/", nil)
					r.With(auth, sysadmin).Delete("/", nil)

					r.With(auth).Put("/", nil)
					r.With(auth).Patch("/slug", nil)
					r.With(auth).Patch("/activate", nil)
					r.With(auth).Patch("/deactivate", nil)
				})
			})

			r.Route("/memebers", func(r chi.Router) {
				r.Route("/{member_id}", func(r chi.Router) {
					r.Get("/", nil)
					r.With(auth).Put("/", nil)
					r.With(auth).Delete("/", nil)
				})
			})

			r.Route("/invite", func(r chi.Router) {
				r.Post("/", nil)

				r.With(auth).Route("/{invite_id}", func(r chi.Router) {
					r.Get("/", nil)
					r.Patch("/accept", nil)
					r.Patch("/decline", nil)
				})
			})

			r.Route("/roles", func(r chi.Router) {
				r.Post("/", nil)

				r.Route("/{role_id}", func(r chi.Router) {
					r.Get("/", nil)
					r.Put("/", nil)
					r.Delete("/", nil)
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
