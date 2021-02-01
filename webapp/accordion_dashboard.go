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

func HandleAccordionDashGet(w http.ResponseWriter, r *http.Request) {
	// Init template variables
	tmplVars := &AccordionDashTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	tmplVars.HiveLinkCount, err = modelClient.CountHiveHeaderLinks()
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	tmplVars.Headers, err = modelClient.ReadAccordionHeaders()
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	err = templates.ExecuteTemplate(w, "accordion_dashboard", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}
