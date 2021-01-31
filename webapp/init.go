package webapp

import (
	"context"
	"encoding/gob"
	"errors"
	"github.com/Hive-Gay/supreme-robot/models"
	"github.com/coreos/go-oidc"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
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

type contextKey int

const SessionKey contextKey = 0
const UserKey contextKey = 1

const groupMailAdmin = "/Mail Admin"

var adminGroups = []string{
	groupMailAdmin,
}

var (
	ctx            context.Context
	logger         *loggo.Logger
	store          *redistore.RediStore
	oauth2Config   oauth2.Config
	oauth2Verifier *oidc.IDTokenVerifier
	templates      *template.Template
)

func Close() {
	store.Close()
}

func Init(rp *redis.Pool) error {
	newLogger := loggo.GetLogger("web")
	logger = &newLogger

	// Load Templates
	templateDir := pkger.Include("/webapp/templates")
	t, err := compileTemplates(templateDir)
	if err != nil {
		return err
	}
	templates = t

	// Fetch new store.
	Secret := os.Getenv("SECRET")
	if Secret == "" {
		return errors.New("missing env var SECRET")
	}

	store, err = redistore.NewRediStoreWithPool(rp, []byte(Secret))
	if err != nil {
		return err
	}

	// Register models for GOB
	gob.Register(models.User{})
	gob.Register(templateAlert{})

	// Configure Oauth
	AppHostname := os.Getenv("APP_HOSTNAME")
	if AppHostname == "" {
		return errors.New("missing env var APP_HOSTNAME")
	}

	callbackURL := &url.URL{
		Scheme: "https",
		Host:   AppHostname,
		Path:   "/oauth/callback",
	}

	if strings.ToUpper(os.Getenv("OAUTH_CALLBACK_HTTPS")) == "FALSE" {
		callbackURL.Scheme = "http"
	}

	OAuthClientID := os.Getenv("OAUTH_CLIENT_ID")
	if OAuthClientID == "" {
		return errors.New("missing env var OAUTH_CLIENT_ID")
	}

	OAuthClientSecret := os.Getenv("OAUTH_CLIENT_SECRET")
	if OAuthClientSecret == "" {
		return errors.New("missing env var OAUTH_CLIENT_SECRET")
	}

	OAuthProviderURL := os.Getenv("OAUTH_PROVIDER_URL")
	if OAuthProviderURL == "" {
		return errors.New("missing env var OAUTH_PROVIDER_URL")
	}

	ctx = context.Background()
	provider, err := oidc.NewProvider(ctx, OAuthProviderURL)
	if err != nil {
		logger.Errorf("Could not create oidc provider: %s", err.Error())
		return err
	}

	oidcConfig := &oidc.Config{ClientID: OAuthClientID}
	oauth2Verifier = provider.Verifier(oidcConfig)

	// Configure OAuth2 client
	oauth2Config = oauth2.Config{
		ClientID:     OAuthClientID,
		ClientSecret: OAuthClientSecret,
		RedirectURL:  callbackURL.String(),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	// Setup Router
	r := mux.NewRouter()
	r.Use(Middleware)

	// Error Pages
	r.NotFoundHandler = NotFoundHandler()
	r.MethodNotAllowedHandler = MethodNotAllowedHandler()

	// Static Files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(pkger.Dir("/webapp/static"))))

	r.HandleFunc("/", HandleAccordion).Methods("GET")
	r.HandleFunc("/login", HandleLogin).Methods("GET")
	r.HandleFunc("/logout", HandleLogout).Methods("GET")
	r.HandleFunc("/oauth/callback", HandleOauthCallback).Methods("GET")

	// Protected Pages
	protected := r.PathPrefix("/app/").Subrouter()
	protected.Use(MiddlewareRequireAuth)
	protected.HandleFunc("/", GetHome).Methods("GET")

	// Accordion Dashboard
	protected.HandleFunc("/accordion", HandleAccordionDashGet).Methods("GET")
	protected.HandleFunc("/accordion/add", HandleAccordionHeaderAddGet).Methods("GET")
	protected.HandleFunc("/accordion/add", HandleAccordionHeaderAddPost).Methods("POST")
	protected.HandleFunc("/accordion/{header:[0-9]+}", HandleAccordionHeaderGet).Methods("GET")
	protected.HandleFunc("/accordion/{header:[0-9]+}/add", HandleAccordionLinkAddGet).Methods("GET")
	protected.HandleFunc("/accordion/{header:[0-9]+}/add", HandleAccordionLinkAddPost).Methods("POST")
	protected.HandleFunc("/accordion/{header:[0-9]+}/delete", HandleAccordionHeaderDeleteGet).Methods("GET")
	protected.HandleFunc("/accordion/{header:[0-9]+}/delete", HandleAccordionHeaderDeletePost).Methods("POST")
	protected.HandleFunc("/accordion/{header:[0-9]+}/edit", HandleAccordionHeaderEditGet).Methods("GET")
	protected.HandleFunc("/accordion/{header:[0-9]+}/edit", HandleAccordionHeaderEditPost).Methods("POST")
	protected.HandleFunc("/accordion/{header:[0-9]+}/{link:[0-9]+}/delete", HandleAccordionLinkDeleteGet).Methods("GET")
	protected.HandleFunc("/accordion/{header:[0-9]+}/{link:[0-9]+}/delete", HandleAccordionLinkDeletePost).Methods("POST")
	protected.HandleFunc("/accordion/{header:[0-9]+}/{link:[0-9]+}/edit", HandleAccordionLinkEditGet).Methods("GET")
	protected.HandleFunc("/accordion/{header:[0-9]+}/{link:[0-9]+}/edit", HandleAccordionLinkEditPost).Methods("POST")



	logger.Debugf("starting webapp server")
	go func() {
		srv := &http.Server{
			Handler:      r,
			Addr:         ":5000",
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		err := srv.ListenAndServe()
		if err != nil {
			logger.Errorf("Could not start webapp server %s", err.Error())
		}
	}()

	return nil
}
