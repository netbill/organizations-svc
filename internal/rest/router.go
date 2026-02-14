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
	UpdateOrganizationMediaContent() func(next http.Handler) http.Handler
	Logger(log *logium.Entry) func(next http.Handler) http.Handler
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
	Port     string `mapstructure:"port"`
	Timeouts struct {
		Read       time.Duration `mapstructure:"read"`
		ReadHeader time.Duration `mapstructure:"read_header"`
		Write      time.Duration `mapstructure:"write"`
		Idle       time.Duration `mapstructure:"idle"`
	} `mapstructure:"timeouts"`
}

func (s *Server) Run(ctx context.Context, log *logium.Entry, cfg Config) {
	auth := s.middlewares.AccountAuth()
	updOrganization := s.middlewares.UpdateOrganizationMediaContent()

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
						r.With(updOrganization).Delete("/", s.handlers.CancelUpdateOrganization)

						r.With(updOrganization).Put("/confirm", s.handlers.ConfirmUpdateOrganization)
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
		})
	})

	srv := &http.Server{
		Addr:              cfg.Port,
		Handler:           r,
		ReadTimeout:       cfg.Timeouts.Read,
		ReadHeaderTimeout: cfg.Timeouts.ReadHeader,
		WriteTimeout:      cfg.Timeouts.Write,
		IdleTimeout:       cfg.Timeouts.Idle,
	}

	log.Infof("starting http service on %s", cfg.Port)

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
		log.Warnf("shutting down http service...")
	case err := <-errCh:
		if err != nil {
			log.Errorf("http server error: %v", err)
		}
	}

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		log.Errorf("http shutdown error: %v", err)
	} else {
		log.Warnf("http server stopped")
	}
}
