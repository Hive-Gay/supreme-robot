package webapp

import (
	"context"
	"encoding/gob"
	"errors"
	"github.com/Hive-Gay/supreme-robot/jobs"
	"github.com/Hive-Gay/supreme-robot/models"
	"github.com/Hive-Gay/supreme-robot/twilio"
	"github.com/coreos/go-oidc"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"golang.org/x/oauth2"
	"gopkg.in/boj/redistore.v1"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Server struct {
	apphostname    string
	enqueuer       *jobs.Enqueuer
	ctx            context.Context
	modelClient    *models.Client
	store          *redistore.RediStore
	oauth2Config   oauth2.Config
	oauth2Verifier *oidc.IDTokenVerifier
	router         *mux.Router
	templates      *template.Template
	twilioClient   *twilio.Client
}

func NewServer(redisAddress string, mc *models.Client, e *jobs.Enqueuer, tc *twilio.Client) (*Server, error) {
	server := Server{
		modelClient:  mc,
		enqueuer:     e,
		twilioClient: tc,
	}

	// Load Templates
	templateDir := pkger.Include("/webapp/templates")
	t, err := compileTemplates(templateDir)
	if err != nil {
		return nil, err
	}
	server.templates = t

	// Fetch new store.
	Secret := os.Getenv("SECRET")
	if Secret == "" {
		return nil, errors.New("missing env var SECRET")
	}

	server.store, err = redistore.NewRediStoreWithPool(&redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisAddress)
		},
	}, []byte(Secret))
	if err != nil {
		return nil, err
	}

	// Register models for GOB
	gob.Register(models.User{})
	gob.Register(templateAlert{})

	// Configure Oauth
	server.apphostname = os.Getenv("APP_HOSTNAME")
	if server.apphostname == "" {
		return nil, errors.New("missing env var APP_HOSTNAME")
	}

	callbackURL := &url.URL{
		Scheme: "https",
		Host:   server.apphostname,
		Path:   "/oauth/callback",
	}

	if strings.ToUpper(os.Getenv("OAUTH_CALLBACK_HTTPS")) == "FALSE" {
		callbackURL.Scheme = "http"
	}

	OAuthClientID := os.Getenv("OAUTH_CLIENT_ID")
	if OAuthClientID == "" {
		return nil, errors.New("missing env var OAUTH_CLIENT_ID")
	}

	OAuthClientSecret := os.Getenv("OAUTH_CLIENT_SECRET")
	if OAuthClientSecret == "" {
		return nil, errors.New("missing env var OAUTH_CLIENT_SECRET")
	}

	OAuthProviderURL := os.Getenv("OAUTH_PROVIDER_URL")
	if OAuthProviderURL == "" {
		return nil, errors.New("missing env var OAUTH_PROVIDER_URL")
	}

	server.ctx = context.Background()
	provider, err := oidc.NewProvider(server.ctx, OAuthProviderURL)
	if err != nil {
		logger.Errorf("Could not create oidc provider: %s", err.Error())
		return nil, err
	}

	oidcConfig := &oidc.Config{ClientID: OAuthClientID}
	server.oauth2Verifier = provider.Verifier(oidcConfig)

	// Configure OAuth2 client
	server.oauth2Config = oauth2.Config{
		ClientID:     OAuthClientID,
		ClientSecret: OAuthClientSecret,
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
	server.router.HandleFunc("/webhook/sms", server.HandleWebhookSMSPost).Methods("POST")

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

	// Mail Dashboard
	protected.HandleFunc("/admin/mail", server.HandleMailDashGet).Methods("GET")

	return &server, nil
}

func (s *Server) Run() error {

	srv := &http.Server{
		Handler:      s.router,
		Addr:         ":5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return srv.ListenAndServe()
}
