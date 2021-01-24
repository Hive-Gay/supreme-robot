package models

type User struct {
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
	Groups            []string `json:"groups"`
	IssuedAt          int64    `json:"iat"`
	Issuer            string   `json:"iss"`
	SessionState      string   `json:"session_state"`
	TokenID           string   `json:"jti"`
	Type              string   `json:"typ"`
}
