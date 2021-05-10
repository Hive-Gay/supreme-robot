package webapp

import (
	"net/http"
)

type HomeTemplate struct {
	templateCommon
}

func (s *Server) GetHome(w http.ResponseWriter, r *http.Request) {

	// Init template variables
	tmplVars := &HomeTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	err = s.templates.ExecuteTemplate(w, "home", tmplVars)
	if err != nil {
		logger.Errorf("could not render home template: %s", err.Error())
	}

}
