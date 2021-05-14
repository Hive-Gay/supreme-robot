package webapp

import (
	"context"
	"encoding/gob"
	"github.com/Hive-Gay/go-hivelib/clients/v1/quotes"
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
	// data stuff
	redis *redis.Client
	db    *database.Client

	// hive services
	quotes *quotes.Client

	// web stuff
	extHostname    string
	ctx            context.Context
	store          *redisstore.RedisStore
	oauth2Config   oauth2.Config
	oauth2Verifier *oidc.IDTokenVerifier
	router         *mux.Router
	server         *http.Server
	templates      *template.Template
}

func NewServer(cfg *config.Config, rc *redis.Client, mc *database.Client, qc *quotes.Client) (*Server, error) {
	server := Server{
		db:          mc,
		redis:       rc,
		quotes:      qc,
		extHostname: cfg.ExtHostname,
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
	gob.Register(OAuthUser{})
	gob.Register(templateAlert{})

	// Configure Oauth
	callbackURL := &url.URL{
		Scheme: "https",
		Host:   server.extHostname,
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
	protected.HandleFunc("/accordion", server.AccordionDashGetHandler).Methods("GET")
	protected.HandleFunc("/accordion/add", server.AccordionHeaderAddGetHandler).Methods("GET")
	protected.HandleFunc("/accordion/add", server.AccordionHeaderAddPostHandler).Methods("POST")
	protected.HandleFunc("/accordion/{header:[0-9]+}", server.AccordionHeaderGetHandler).Methods("GET")
	protected.HandleFunc("/accordion/{header:[0-9]+}/add", server.AccordionLinkAddGetHandler).Methods("GET")
	protected.HandleFunc("/accordion/{header:[0-9]+}/add", server.AccordionLinkAddPostHandler).Methods("POST")
	protected.HandleFunc("/accordion/{header:[0-9]+}/delete", server.AccordionHeaderDeleteGetHandler).Methods("GET")
	protected.HandleFunc("/accordion/{header:[0-9]+}/delete", server.AccordionHeaderDeletePostHandler).Methods("POST")
	protected.HandleFunc("/accordion/{header:[0-9]+}/edit", server.AccordionHeaderEditGetHandler).Methods("GET")
	protected.HandleFunc("/accordion/{header:[0-9]+}/edit", server.AccordionHeaderEditPostHandler).Methods("POST")
	protected.HandleFunc("/accordion/{header:[0-9]+}/{link:[0-9]+}/delete", server.AccordionLinkDeleteGetHandler).Methods("GET")
	protected.HandleFunc("/accordion/{header:[0-9]+}/{link:[0-9]+}/delete", server.AccordionLinkDeletePostHandler).Methods("POST")
	protected.HandleFunc("/accordion/{header:[0-9]+}/{link:[0-9]+}/edit", server.AccordionLinkEditGetHandler).Methods("GET")
	protected.HandleFunc("/accordion/{header:[0-9]+}/{link:[0-9]+}/edit", server.AccordionLinkEditPostHandler).Methods("POST")

	// Quotes Dashboard
	protected.HandleFunc("/quotes", server.QuotesDashGetHandler).Methods("GET")
	protected.HandleFunc("/quotes/isms", server.QuotesIsmsGetHandler).Methods("GET")
	protected.HandleFunc("/quotes/ism_add", server.QuotesIsmAddGetHandler).Methods("GET")
	protected.HandleFunc("/quotes/ism_add", server.QuotesIsmAddPostHandler).Methods("POST")
	protected.HandleFunc("/quotes/sayers", server.QuotesSayersGetHandler).Methods("GET")
	protected.HandleFunc("/quotes/sayer_add", server.QuotesSayerAddGetHandler).Methods("GET")
	protected.HandleFunc("/quotes/sayer_add", server.QuotesSayerAddPostHandler).Methods("POST")

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
