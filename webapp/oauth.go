package webapp

import (
	"encoding/json"
	"github.com/gorilla/sessions"
	"net/http"
)

type OAuthUser struct {
	Subject           string   `json:"sub"`
	Email             *string  `json:"email"`
	Name              *string  `json:"name"`
	PreferredUsername *string  `json:"preferred_username"`
	Audience          string   `json:"aud"`
	AuthnCtxClsRef    string   `json:"acr"`
	AuthorizedParty   string   `json:"azp"`
	AuthTime          int64    `json:"auth_time"`
	EmailVerified     bool     `json:"email_verified"`
	ExpiresAt         int64    `json:"exp"`
	Groups            []string `json:"group"`
	IssuedAt          int64    `json:"iat"`
	Issuer            string   `json:"iss"`
	SessionState      string   `json:"session_state"`
	TokenID           string   `json:"jti"`
	Type              string   `json:"typ"`
}

func (s *Server) HandleOauthCallback(w http.ResponseWriter, r *http.Request) {
	// Init Session
	us := r.Context().Value(SessionKey).(*sessions.Session)

	// get state
	val := us.Values["oauth-state"]
	var state string
	var ok bool
	if state, ok = val.(string); !ok {
		// redirect home page if no login-redirect
		logger.Warningf("Invalid State")
		http.Error(w, "state invalid", http.StatusBadRequest)
		return
	}
	if r.URL.Query().Get("state") != state {
		logger.Warningf("States don't match")
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	// display error
	if r.URL.Query().Get("error") != "" {
		http.Error(w, r.URL.Query().Get("error_description"), http.StatusForbidden)
		return
	}

	// process returned oauth data
	user, err := s.processCallback(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not logged in", http.StatusForbidden)
		return
	}

	// Insert into database or update existing record0
	//err = user.Upsert()
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	// Save user to session
	us.Values["oauth-state"] = nil
	us.Values["user"] = user
	err = us.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect to last page
	val = us.Values["login-redirect"]
	var loginRedirect string
	if loginRedirect, ok = val.(string); !ok {
		// redirect home page if no login-redirect
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	http.Redirect(w, r, loginRedirect, http.StatusFound)
	return
}

func (s *Server) processCallback(r *http.Request) (*OAuthUser, error) {
	oauth2Token, err := s.oauth2Config.Exchange(s.ctx, r.URL.Query().Get("code"))
	if err != nil {
		return nil, err
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, err
	}
	idToken, err := s.oauth2Verifier.Verify(s.ctx, rawIDToken)
	if err != nil {
		return nil, err
	}

	IDTokenClaims := new(json.RawMessage)
	if err := idToken.Claims(&IDTokenClaims); err != nil {
		return nil, err
	}

	logger.Tracef("Response From OAUTH: %s", IDTokenClaims)

	user := OAuthUser{}
	if err := json.Unmarshal(*IDTokenClaims, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
