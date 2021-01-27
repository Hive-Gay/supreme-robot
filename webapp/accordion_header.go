package webapp

import (
	"fmt"
	"github.com/Hive-Gay/supreme-robot/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"strconv"
)

type AccordionHeaderTemplate struct {
	templateCommon

	Header *models.AccordionHeader
	Links  []*models.AccordionLink
}

type AccordionHeaderFormTemplate struct {
	templateCommon

	TitleText              string
	FormInputTitleDisabled bool
	FormInputTitleValue    string
	FormButtonSubmitColor  string
	FormButtonSubmitText   string
}

func HandleAccordionHeaderAddGet(w http.ResponseWriter, r *http.Request) {
	// Init template variables
	tmplVars := &AccordionHeaderFormTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
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
		return
	}

	tmplVars.TitleText = "Edit Header"
	tmplVars.FormButtonSubmitColor = "warning"
	tmplVars.FormButtonSubmitText = "Update"

	err = templates.ExecuteTemplate(w, "accordion_header_form", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

func HandleAccordionHeaderGet(w http.ResponseWriter, r *http.Request) {
	// get responder
	vars := mux.Vars(r)

	// Init template variables
	tmplVars := &AccordionHeaderTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	headerID, err := strconv.Atoi(vars["id"])
	if err != nil {
		returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if headerID == 0 {
		tmplVars.Header = &models.AccordionHeader{
			ID:    0,
			Title: "The Hive",
		}
	} else {
		tmplVars.Header, err = models.ReadAccordionHeader(headerID)
		if err != nil {
			returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		if tmplVars.Header == nil {
			returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("header %d not found", headerID))
			return
		}
	}

	err = templates.ExecuteTemplate(w, "accordion_header_view", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}
