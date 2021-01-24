package webapp

import (
	"github.com/Hive-Gay/supreme-robot/util"
	"github.com/gorilla/sessions"
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Init Session
	us := r.Context().Value(SessionKey).(*sessions.Session)

	newState := util.RandString(16)
	us.Values["oauth-state"] = newState
	err := us.Save(r, w)
	if err != nil {
		logger.Warningf("Could not save session: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, oauth2Config.AuthCodeURL(newState), http.StatusFound)
}
