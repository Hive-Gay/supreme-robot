package webapp

import (
	"database/sql"
	"fmt"
	"github.com/Hive-Gay/supreme-robot/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"strconv"
)

type AccordionHeaderTemplate struct {
	templateCommon
	Breadcrumbs *[]templateBreadcrumb

	Header *models.AccordionHeader
	Links  []*models.AccordionLink
}

type AccordionHeaderFormTemplate struct {
	templateCommon
	Breadcrumbs *[]templateBreadcrumb

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

	// Configure Form
	tmplVars.TitleText = "Add Header"
	tmplVars.FormButtonSubmitColor = "success"
	tmplVars.FormButtonSubmitText = "Add"

	// breadcrumbs
	tmplVars.Breadcrumbs = &[]templateBreadcrumb{
		{
			HRef: "/app/accordion",
			Text: "Accordion",
		},
		{
			Text: tmplVars.TitleText,
		},
	}

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

	err = modelClient.CreateAccordionHeader(&ah)
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

func HandleAccordionHeaderDeleteGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Init template variables
	tmplVars := &AccordionHeaderFormTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	headerID, err := strconv.Atoi(vars["header"])
	if err != nil {
		returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	header, err := modelClient.ReadAccordionHeader(headerID)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if header == nil {
		returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("header not found: %d", headerID))
		return
	}

	// Configure Form
	tmplVars.FormInputTitleValue = header.Title
	tmplVars.FormInputTitleDisabled = true

	tmplVars.TitleText = "Delete Header"
	tmplVars.FormButtonSubmitColor = "danger"
	tmplVars.FormButtonSubmitText = "Delete"

	// breadcrumbs
	tmplVars.Breadcrumbs = &[]templateBreadcrumb{
		{
			HRef: "/app/accordion",
			Text: "Accordion",
		},
		{
			Text: tmplVars.TitleText,
		},
	}

	err = templates.ExecuteTemplate(w, "accordion_header_form", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

func HandleAccordionHeaderDeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	headerID, err := strconv.Atoi(vars["header"])
	if err != nil {
		returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	header, err := modelClient.ReadAccordionHeader(headerID)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if header == nil {
		returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("header not found: %d", headerID))
		return
	}

	err = modelClient.DeleteAccordionHeader(headerID)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	us := r.Context().Value(SessionKey).(*sessions.Session)
	us.Values["page-alert-success"] = templateAlert{Text: "Header deleted"}
	err = us.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect
	http.Redirect(w, r, "/app/accordion", http.StatusFound)
}

func HandleAccordionHeaderEditGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Init template variables
	tmplVars := &AccordionHeaderFormTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	headerID, err := strconv.Atoi(vars["header"])
	if err != nil {
		returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	header, err := modelClient.ReadAccordionHeader(headerID)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if header == nil {
		returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("header not found: %d", headerID))
		return
	}

	// Configure Form
	tmplVars.FormInputTitleValue = header.Title

	tmplVars.TitleText = "Edit Header"
	tmplVars.FormButtonSubmitColor = "warning"
	tmplVars.FormButtonSubmitText = "Update"

	// breadcrumbs
	tmplVars.Breadcrumbs = &[]templateBreadcrumb{
		{
			HRef: "/app/accordion",
			Text: "Accordion",
		},
		{
			Text: tmplVars.TitleText,
		},
	}

	err = templates.ExecuteTemplate(w, "accordion_header_form", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

func HandleAccordionHeaderEditPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	headerID, err := strconv.Atoi(vars["header"])
	if err != nil {
		returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	header, err := modelClient.ReadAccordionHeader(headerID)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if header == nil {
		returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("header not found: %d", headerID))
		return
	}

	// parse form data
	err = r.ParseForm()
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	header.Title =  r.Form.Get("title")

	err = modelClient.UpdateAccordionHeaders(header)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	us := r.Context().Value(SessionKey).(*sessions.Session)
	us.Values["page-alert-success"] = templateAlert{Text: "Header updated"}
	err = us.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect
	http.Redirect(w, r, "/app/accordion", http.StatusFound)
}

func HandleAccordionHeaderGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Init template variables
	tmplVars := &AccordionHeaderTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	headerID, err := strconv.Atoi(vars["header"])
	if err != nil {
		returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if headerID == 0 {
		tmplVars.Header = &models.AccordionHeader{
			ID:    0,
			Title: "The Hive",
		}

		tmplVars.Links, err = modelClient.ReadAccordionLinks(sql.NullInt32{Valid: false})
		if err != nil {
			returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		tmplVars.Header, err = modelClient.ReadAccordionHeader(headerID)
		if err != nil {
			returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		if tmplVars.Header == nil {
			returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("header %d not found", headerID))
			return
		}

		tmplVars.Links, err = modelClient.ReadAccordionLinks(sql.NullInt32{Valid: true, Int32: int32(tmplVars.Header.ID)})
		if err != nil {
			returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// breadcrumbs
	tmplVars.Breadcrumbs = &[]templateBreadcrumb{
		{
			HRef: "/app/accordion",
			Text: "Accordion",
		},
		{
			Text: tmplVars.Header.Title,
		},
	}

	err = templates.ExecuteTemplate(w, "accordion_header_view", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}
