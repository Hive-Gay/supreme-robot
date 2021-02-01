package webapp

import (
	"github.com/Hive-Gay/supreme-robot/models"
	"net/http"
)

type AccordionDashTemplate struct {
	templateCommon

	HiveLinkCount int
	Headers []*models.AccordionHeader
}

func (s *Server)HandleAccordionDashGet(w http.ResponseWriter, r *http.Request) {
	// Init template variables
	tmplVars := &AccordionDashTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	tmplVars.HiveLinkCount, err = s.modelClient.CountHiveHeaderLinks()
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	tmplVars.Headers, err = s.modelClient.ReadAccordionHeaders()
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	err = s.templates.ExecuteTemplate(w, "accordion_dashboard", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}
