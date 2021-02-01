package webapp

import (
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"net/http"
)

func (s *Server)HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Init Session
	us := r.Context().Value(SessionKey).(*sessions.Session)

	newState := uuid.New().String()
	us.Values["oauth-state"] = newState
	err := us.Save(r, w)
	if err != nil {
		logger.Warningf("Could not save session: %s", err.Error())
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, s.oauth2Config.AuthCodeURL(newState), http.StatusFound)
}
