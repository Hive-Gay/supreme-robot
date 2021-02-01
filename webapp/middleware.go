package webapp

import (
	"context"
	"github.com/Hive-Gay/supreme-robot/models"
	"github.com/Hive-Gay/supreme-robot/util"
	"github.com/gorilla/sessions"
	"net/http"
	"time"
)

const layoutCombined = "02/Jan/2006:15:04:05 -0700"

type ResponseWriterX struct {
	http.ResponseWriter
	status     int
	bodyLength int
}

func (w *ResponseWriterX) Write(b []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(b)
	w.bodyLength += n
	return
}

func (r *ResponseWriterX) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.status = status
	return
}

func (s *Server)Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wx := &ResponseWriterX{
			ResponseWriter: w,
			status:         200,
			bodyLength:     0,
		}

		// Init Session
		us, err := s.store.Get(r, "supreme-robot")
		if err != nil {
			logger.Infof("got %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), SessionKey, us)

		// Retrieve our user and type-assert it
		val := us.Values["user"]
		var user = models.User{}
		var ok bool
		if user, ok = val.(models.User); ok {
			ctx = context.WithValue(ctx, UserKey, &user)
		}

		// Do Request
		next.ServeHTTP(wx, r.WithContext(ctx))

		// Log It
		duration := time.Since(start)
		logger.Debugf("%s - %s [%s] \"%s %s %s\" %d %d \"%s\" \"%s\" rt=%d",
			r.RemoteAddr,
			"-",
			start.Format(layoutCombined),
			r.Method,
			r.URL.Path,
			r.Proto,
			wx.status,
			wx.bodyLength,
			r.Referer(),
			r.UserAgent(),
			duration.Milliseconds(),
		)
	})
}

func (s *Server)MiddlewareRequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		us := r.Context().Value(SessionKey).(*sessions.Session)

		if r.Context().Value(UserKey) == nil {
			// Save current page
			us.Values["login-redirect"] = r.URL.Path
			err := us.Save(r, w)
			if err != nil {
				s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
				return
			}

			// redirect to login
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		user := r.Context().Value(UserKey).(*models.User)
		now := time.Now()
		if user.ExpiresAt < now.Unix() {
			// Save current page
			us.Values["login-redirect"] = r.URL.Path
			err := us.Save(r, w)
			if err != nil {
				s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
				return
			}

			// redirect to login
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Check for UWU Crew group
		UWUCrew := util.ContainsString(user.Groups, groupUWUCrew)
		if !UWUCrew {
			s.returnErrorPage(w, r, http.StatusUnauthorized, "Ask Tyr to join the UWU Crew")
			return
		}


		next.ServeHTTP(w, r)
	})
}
