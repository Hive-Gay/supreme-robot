package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	ExtHostname string

	LoggerConfig string

	OAuthCallbackHTTPS bool
	OAuthClientID      string
	OAuthClientSecret  string
	OAuthProviderURL   string

	PostgresDsn string

	RedisAddress  string
	RedisDB       int
	RedisPassword string

	Secret string
}

func CollectConfig() *Config {
	var missingEnv []string
	var config Config

	// EXT_HOSTNAME
	config.ExtHostname = os.Getenv("EXT_HOSTNAME")
	if config.ExtHostname == "" {
		missingEnv = append(missingEnv, "EXT_HOSTNAME")
	}

	// LOG_LEVEL
	config.LoggerConfig = os.Getenv("LOG_LEVEL")
	if config.LoggerConfig == "" {
		config.LoggerConfig = "<root>=INFO"
	} else {
		config.LoggerConfig = fmt.Sprintf("<root>=%s", strings.ToUpper(config.LoggerConfig))
	}

	// OAUTH_CALLBACK_HTTPS
	oauthCallbackHttps := os.Getenv("OAUTH_CALLBACK_HTTPS")
	if strings.ToLower(oauthCallbackHttps) == "false" {
		config.OAuthCallbackHTTPS = false
	} else {
		config.OAuthCallbackHTTPS = true
	}

	// OAUTH_CLIENT_ID
	config.OAuthClientID = os.Getenv("OAUTH_CLIENT_ID")
	if config.OAuthClientID == "" {
		missingEnv = append(missingEnv, "OAUTH_CLIENT_ID")
	}

	// OAUTH_CLIENT_SECRET
	config.OAuthClientSecret = os.Getenv("OAUTH_CLIENT_SECRET")
	if config.OAuthClientSecret == "" {
		missingEnv = append(missingEnv, "OAUTH_CLIENT_SECRET")
	}

	// OAUTH_PROVIDER_URL
	config.OAuthProviderURL = os.Getenv("OAUTH_PROVIDER_URL")
	if config.OAuthProviderURL == "" {
		missingEnv = append(missingEnv, "OAUTH_PROVIDER_URL")
	}

	// POSTGRES_DSN
	config.PostgresDsn = os.Getenv("POSTGRES_DSN")
	if config.PostgresDsn == "" {
		missingEnv = append(missingEnv, "POSTGRES_DSN")
	}

	// REDIS_ADDRESS
	config.RedisAddress = os.Getenv("REDIS_ADDRESS")
	if config.RedisAddress == "" {
		missingEnv = append(missingEnv, "REDIS_ADDRESS")
	}

	// REDIS_ADDRESS
	config.RedisPassword = os.Getenv("REDIS_ADDRESS")

	// Validation
	if len(missingEnv) > 0 {
		msg := fmt.Sprintf("Environment variables missing: %v", missingEnv)
		panic(fmt.Sprint(msg))
	}

	return &config
}
