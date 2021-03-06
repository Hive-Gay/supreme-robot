package webapp

import (
	"github.com/Hive-Gay/supreme-robot/util"
	"net/http"
)

func userAuthed(r *http.Request, role string) bool {
	if r.Context().Value(UserKey) != nil {
		user := r.Context().Value(UserKey).(*OAuthUser)
		if util.ContainsString(user.Groups, role) {
			return true
		}
	}
	return false
}
