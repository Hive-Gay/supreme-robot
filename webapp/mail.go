package webapp

import (
	"net/http"
)

type MailDashTemplate struct {
	templateCommon
}

func HandleMailDashGet(w http.ResponseWriter, r *http.Request) {
	if !userAuthed(r, groupMailAdmin) {
		returnErrorPage(w, r, http.StatusForbidden, "")
		return
	}

	// Init template variables
	tmplVars := &MailDashTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	err = templates.ExecuteTemplate(w, "mail_dashboard", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

