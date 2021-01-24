package webapp

import (
	"github.com/gorilla/sessions"
	"net/http"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Init Session
	us := r.Context().Value(SessionKey).(*sessions.Session)

	// Set user to nil
	us.Values["user"] = nil
	err := us.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
