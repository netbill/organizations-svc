package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/netbill/organizations-svc/internal/media"
	"github.com/netbill/organizations-svc/pkg/log"
	"github.com/netbill/restkit/tokens"
)

type OrgController interface {
	Create(w http.ResponseWriter, r *http.Request)

	Get(w http.ResponseWriter, r *http.Request)
	GetList(w http.ResponseWriter, r *http.Request)

	CreateUploadMediaLink(w http.ResponseWriter, r *http.Request)

	Update(w http.ResponseWriter, r *http.Request)

	Delete(w http.ResponseWriter, r *http.Request)

	DeleteUploadMedia(w http.ResponseWriter, r *http.Request)

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
}

type InviteController interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)

	GetList(w http.ResponseWriter, r *http.Request)

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
	ResolverUrl(resolver *media.Resolver) func(next http.Handler) http.Handler
}

type Server struct {
	middlewares Middlewares
	org         OrgController
	member      MemberController
	invite      InviteController

	log           *log.Logger
	mediaResolver *media.Resolver
	config        Config
}

type ServerDeps struct {
	Middlewares Middlewares

	Org    OrgController
	Member MemberController
	Invite InviteController

	Log           *log.Logger
	MediaResolver *media.Resolver
	Config        Config
}

type Config struct {
	Port              int
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func NewServer(deps ServerDeps) *Server {
	return &Server{
		middlewares:   deps.Middlewares,
		org:           deps.Org,
		member:        deps.Member,
		invite:        deps.Invite,
		log:           deps.Log,
		mediaResolver: deps.MediaResolver,
		config:        deps.Config,
	}
}

func (s *Server) Run(ctx context.Context) {
	auth := s.middlewares.AccountAuth()
	sysadmin := s.middlewares.AccountAuth(tokens.RoleSystemAdmin)

	r := chi.NewRouter()
	r.Use(
		s.middlewares.Logger(s.log),
		s.middlewares.ResolverUrl(s.mediaResolver),
		s.middlewares.CorsDocs(),
	)

	r.Route("/organizations-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {

			r.Route("/organizations", func(r chi.Router) {
				r.Get("/", s.org.GetList)
				r.With(auth).Post("/", s.org.Create)

				r.Route("/{organization_id}", func(r chi.Router) {
					r.Get("/", s.org.Get)
					r.With(auth).Patch("/", s.org.Update)
					r.With(auth).Delete("/", s.org.Delete)

					r.With(auth).Route("/media", func(r chi.Router) {
						r.Post("/", s.org.CreateUploadMediaLink)
						r.Delete("/", s.org.DeleteUploadMedia)
					})

					r.With(auth).Patch("/activate", s.org.Activate)
					r.With(auth).Patch("/deactivate", s.org.Deactivate)

					r.With(sysadmin).Patch("/suspend", s.org.Suspend)
					r.With(sysadmin).Patch("/unsuspend", s.org.Unsuspend)
				})
			})

			r.Route("/members", func(r chi.Router) {
				r.Route("/{member_id}", func(r chi.Router) {
					r.Get("/", s.member.Get)
					r.With(auth).Patch("/", s.member.Update)
					r.With(auth).Delete("/", s.member.Delete)
				})

				r.Get("/", s.member.GetList)
			})

			r.With(auth).Route("/invites", func(r chi.Router) {
				r.Post("/", s.invite.Create)
				r.Get("/", s.invite.GetList)

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
		Addr:              fmt.Sprintf(":%d", s.config.Port),
		Handler:           r,
		ReadTimeout:       s.config.ReadTimeout,
		ReadHeaderTimeout: s.config.ReadHeaderTimeout,
		WriteTimeout:      s.config.WriteTimeout,
		IdleTimeout:       s.config.IdleTimeout,
	}

	s.log.WithField("port", s.config.Port).Info("starting http service...")

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
		s.log.Info("shutting down http service...")
	case err := <-errCh:
		if err != nil {
			s.log.WithError(err).Error("http server error")
		}
	}

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		s.log.WithError(err).Error("failed to shutdown http server gracefully")
	} else {
		s.log.Info("http server shutdown gracefully")
	}
}
