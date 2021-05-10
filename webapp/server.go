package webapp

import (
	"context"
	"encoding/gob"
	"github.com/Hive-Gay/supreme-robot/config"
	"github.com/Hive-Gay/supreme-robot/database"
	"github.com/Hive-Gay/supreme-robot/redis"
	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/pkger"
	"github.com/rbcervilla/redisstore/v8"
	"golang.org/x/oauth2"
	"html/template"
	"net/http"
	"net/url"
	"time"
)

type Server struct {
	webapphostname string
	ctx            context.Context

	redis *redis.Client
	db    *database.Client

	store          *redisstore.RedisStore
	oauth2Config   oauth2.Config
	oauth2Verifier *oidc.IDTokenVerifier
	router         *mux.Router
	server         *http.Server
	templates      *template.Template
}

func NewServer(cfg *config.Config, rc *redis.Client, mc *database.Client) (*Server, error) {
	server := Server{
		db:             mc,
		redis:          rc,
		webapphostname: cfg.ExtHostname,
	}

	// Load Templates
	templateDir := pkger.Include("/webapp/templates")
	t, err := compileTemplates(templateDir)
	if err != nil {
		return nil, err
	}
	server.templates = t

	// Fetch new store.
	server.store, err = redisstore.NewRedisStore(rc.Client())
	if err != nil {
		logger.Errorf("create redis store: %s", err.Error())
		return nil, err
	}

	server.store.KeyPrefix(redis.KeySession)
	server.store.Options(sessions.Options{
		Path:   "/",
		Domain: cfg.ExtHostname,
		MaxAge: 86400 * 60,
	})

	// Register database for GOB
	gob.Register(database.User{})
	gob.Register(templateAlert{})

	// Configure Oauth
	callbackURL := &url.URL{
		Scheme: "https",
		Host:   server.webapphostname,
		Path:   "/oauth/callback",
	}

	if !cfg.OAuthCallbackHTTPS {
		callbackURL.Scheme = "http"
	}

	server.ctx = context.Background()
	provider, err := oidc.NewProvider(server.ctx, cfg.OAuthProviderURL)
	if err != nil {
		logger.Errorf("Could not create oidc provider: %s", err.Error())
		return nil, err
	}

	oidcConfig := &oidc.Config{ClientID: cfg.OAuthClientID}
	server.oauth2Verifier = provider.Verifier(oidcConfig)

	// Configure OAuth2 client
	server.oauth2Config = oauth2.Config{
		ClientID:     cfg.OAuthClientID,
		ClientSecret: cfg.OAuthClientSecret,
		RedirectURL:  callbackURL.String(),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	// Setup Router
	server.router = mux.NewRouter()
	server.router.Use(server.Middleware)

	// Error Pages
	server.router.NotFoundHandler = server.NotFoundHandler()
	server.router.MethodNotAllowedHandler = server.MethodNotAllowedHandler()

	// Static Files
	server.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(pkger.Dir("/webapp/static"))))

	server.router.HandleFunc("/", server.HandleAccordion).Methods("GET")
	server.router.HandleFunc("/login", server.HandleLogin).Methods("GET")
	server.router.HandleFunc("/logout", HandleLogout).Methods("GET")
	server.router.HandleFunc("/oauth/callback", server.HandleOauthCallback).Methods("GET")

	// Protected Pages
	protected := server.router.PathPrefix("/app/").Subrouter()
	protected.Use(server.MiddlewareRequireAuth)
	protected.HandleFunc("/", server.GetHome).Methods("GET")

	// Accordion Dashboard
	protected.HandleFunc("/accordion", server.HandleAccordionDashGet).Methods("GET")
	protected.HandleFunc("/accordion/add", server.HandleAccordionHeaderAddGet).Methods("GET")
	protected.HandleFunc("/accordion/add", server.HandleAccordionHeaderAddPost).Methods("POST")
	protected.HandleFunc("/accordion/{header:[0-9]+}", server.HandleAccordionHeaderGet).Methods("GET")
	protected.HandleFunc("/accordion/{header:[0-9]+}/add", server.HandleAccordionLinkAddGet).Methods("GET")
	protected.HandleFunc("/accordion/{header:[0-9]+}/add", server.HandleAccordionLinkAddPost).Methods("POST")
	protected.HandleFunc("/accordion/{header:[0-9]+}/delete", server.HandleAccordionHeaderDeleteGet).Methods("GET")
	protected.HandleFunc("/accordion/{header:[0-9]+}/delete", server.HandleAccordionHeaderDeletePost).Methods("POST")
	protected.HandleFunc("/accordion/{header:[0-9]+}/edit", server.HandleAccordionHeaderEditGet).Methods("GET")
	protected.HandleFunc("/accordion/{header:[0-9]+}/edit", server.HandleAccordionHeaderEditPost).Methods("POST")
	protected.HandleFunc("/accordion/{header:[0-9]+}/{link:[0-9]+}/delete", server.HandleAccordionLinkDeleteGet).Methods("GET")
	protected.HandleFunc("/accordion/{header:[0-9]+}/{link:[0-9]+}/delete", server.HandleAccordionLinkDeletePost).Methods("POST")
	protected.HandleFunc("/accordion/{header:[0-9]+}/{link:[0-9]+}/edit", server.HandleAccordionLinkEditGet).Methods("GET")
	protected.HandleFunc("/accordion/{header:[0-9]+}/{link:[0-9]+}/edit", server.HandleAccordionLinkEditPost).Methods("POST")

	return &server, nil
}

func (s *Server) Close() {
	err := s.server.Close()
	if err != nil {
		logger.Warningf("closing server: %s", err.Error())
	}
}

func (s *Server) ListenAndServe() error {

	s.server = &http.Server{
		Handler:      s.router,
		Addr:         ":5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return s.server.ListenAndServe()
}
