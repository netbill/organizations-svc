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

type OrgController interface {
	Create(w http.ResponseWriter, r *http.Request)

	Get(w http.ResponseWriter, r *http.Request)
	GetList(w http.ResponseWriter, r *http.Request)
	GetMyList(w http.ResponseWriter, r *http.Request)

	CreateUploadMediaLink(w http.ResponseWriter, r *http.Request)

	Update(w http.ResponseWriter, r *http.Request)

	Delete(w http.ResponseWriter, r *http.Request)

	DeleteUploadIcon(w http.ResponseWriter, r *http.Request)
	DeleteUploadBanner(w http.ResponseWriter, r *http.Request)

	Activate(w http.ResponseWriter, r *http.Request)
	Deactivate(w http.ResponseWriter, r *http.Request)
	Suspend(w http.ResponseWriter, r *http.Request)
	Unsuspend(w http.ResponseWriter, r *http.Request)
}

type MemberController interface {
	GetList(w http.ResponseWriter, r *http.Request)

	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)

	LeaveFromOrg(w http.ResponseWriter, r *http.Request)
}

type InviteController interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)

	GetMyList(w http.ResponseWriter, r *http.Request)

	Cancelled(w http.ResponseWriter, r *http.Request)
	Accept(w http.ResponseWriter, r *http.Request)
	Decline(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	AccountAuth(
		allowedRoles ...string,
	) func(next http.Handler) http.Handler
	Logger(log *log.Logger) func(next http.Handler) http.Handler
	CorsDocs() func(next http.Handler) http.Handler
}

type Controllers struct {
}

type Server struct {
	org         OrgController
	member      MemberController
	invite      InviteController
	middlewares Middlewares
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
				r.Get("/", s.org.Get)
				r.With(auth).Post("/", s.org.Create)

				r.Route("/{organization_id}", func(r chi.Router) {
					r.Get("/", s.org.Get)
					r.With(auth).Put("/", s.org.Update)
					r.With(auth).Delete("/", s.org.Delete)

					r.With(auth).Route("/media", func(r chi.Router) {
						r.Route("/upload", func(r chi.Router) {
							r.Post("/url", s.org.CreateUploadMediaLink)

							r.Delete("/icon", s.org.DeleteUploadIcon)
							r.Delete("/banner", s.org.DeleteUploadBanner)
						})
					})

					r.With(auth).Patch("/activate", s.org.Activate)
					r.With(auth).Patch("/deactivate", s.org.Deactivate)

					r.With(sysadmin).Patch("/suspend", s.org.Suspend)
					r.With(sysadmin).Patch("/unsuspend", s.org.Unsuspend)

					r.With(auth).Get("/invites", s.invite.GetMyList)

					r.With(auth).Delete("/leave", s.member.LeaveFromOrg)
				})
			})

			r.Route("/members", func(r chi.Router) {
				r.Route("/{member_id}", func(r chi.Router) {
					r.Get("/", s.member.Get)
					r.With(auth).Put("/", s.member.Update)
					r.With(auth).Delete("/", s.member.Delete)
				})

				r.Get("/", s.member.GetList)
			})

			r.With(auth).Route("/invites", func(r chi.Router) {
				r.Post("/", s.invite.Create)
				r.Get("/me", s.invite.GetMyList)

				r.Route("/{invite_id}", func(r chi.Router) {
					r.Get("/", s.invite.Get)
					r.Patch("/accept", s.invite.Accept)
					r.Patch("/decline", s.invite.Decline)
					r.Patch("/cancel", s.invite.Cancelled)
				})
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
