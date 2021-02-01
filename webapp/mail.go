package webapp

import (
	"net/http"
)

type MailDashTemplate struct {
	templateCommon
}

func (s *Server)HandleMailDashGet(w http.ResponseWriter, r *http.Request) {
	if !userAuthed(r, groupMailAdmin) {
		s.returnErrorPage(w, r, http.StatusForbidden, "")
		return
	}

	// Init template variables
	tmplVars := &MailDashTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	err = s.templates.ExecuteTemplate(w, "mail_dashboard", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

