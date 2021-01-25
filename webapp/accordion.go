package webapp

import (
	"net/http"
)

type AccordionTemplate struct {
	Accordion []AccordionHeader
}

type AccordionHeader struct {
	Title string
	Links []AccordionLink
}

type AccordionLink struct {
	Title string
	Link  string
}

type AccordionDashTemplate struct {
	templateCommon
}

func HandleAccordion(w http.ResponseWriter, r *http.Request) {
	// Init template variables
	tmplVars := &AccordionTemplate{}

	err := templates.ExecuteTemplate(w, "accordion", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

func HandleAccordionDashGet(w http.ResponseWriter, r *http.Request) {
	// Init template variables
	tmplVars := &AccordionDashTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
	}

	err = templates.ExecuteTemplate(w, "accordion_dashboard", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}

}
