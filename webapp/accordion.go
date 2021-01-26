package webapp

import (
	"github.com/Hive-Gay/supreme-robot/models"
	"github.com/gorilla/sessions"
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

	Accordion []*models.AccordionHeader
}

type AccordionHeaderFormTemplate struct {
	templateCommon

	TitleText              string
	FormInputTitleDisabled bool
	FormInputTitleValue    string
	FormButtonSubmitColor  string
	FormButtonSubmitText   string
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
		return
	}

	tmplVars.Accordion, err = models.ReadAccordionHeaders()
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	err = templates.ExecuteTemplate(w, "accordion_dashboard", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

func HandleAccordionHeaderAddGet(w http.ResponseWriter, r *http.Request) {
	// Init template variables
	tmplVars := &AccordionHeaderFormTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
	}

	tmplVars.TitleText = "Add Header"
	tmplVars.FormButtonSubmitColor = "success"
	tmplVars.FormButtonSubmitText = "Add"

	err = templates.ExecuteTemplate(w, "accordion_header_form", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}


func HandleAccordionHeaderAddPost(w http.ResponseWriter, r *http.Request) {
	// parse form data
	err := r.ParseForm()
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	ah := models.AccordionHeader{
		Title: r.Form.Get("title"),
	}

	err = models.CreateAccordionHeader(&ah)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	us := r.Context().Value(SessionKey).(*sessions.Session)
	us.Values["page-alert-success"] = templateAlert{Text: "Header added"}
	err = us.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect home
	http.Redirect(w, r, "/app/accordion", http.StatusFound)
}

func HandleAccordionHeaderEditGet(w http.ResponseWriter, r *http.Request) {
	// Init template variables
	tmplVars := &AccordionHeaderFormTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
	}

	tmplVars.TitleText = "Edit Header"
	tmplVars.FormButtonSubmitColor = "warning"
	tmplVars.FormButtonSubmitText = "Update"

	err = templates.ExecuteTemplate(w, "accordion_header_form", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}