package webapp

import (
	"database/sql"
	"fmt"
	"github.com/Hive-Gay/supreme-robot/database"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"strconv"
)

type AccordionHeaderTemplate struct {
	templateCommon
	Breadcrumbs *[]templateBreadcrumb

	Header *database.AccordionHeader
	Links  []*database.AccordionLink
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

func (s *Server) HandleAccordionHeaderAddGet(w http.ResponseWriter, r *http.Request) {
	// Init template variables
	tmplVars := &AccordionHeaderFormTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
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

	err = s.templates.ExecuteTemplate(w, "accordion_header_form", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

func (s *Server) HandleAccordionHeaderAddPost(w http.ResponseWriter, r *http.Request) {
	// parse form data
	err := r.ParseForm()
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	ah := database.AccordionHeader{
		Title: r.Form.Get("title"),
	}

	err = ah.Create(s.modelClient)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
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

func (s *Server) HandleAccordionHeaderDeleteGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Init template variables
	tmplVars := &AccordionHeaderFormTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	headerID, err := strconv.Atoi(vars["header"])
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	header, err := s.modelClient.ReadAccordionHeader(headerID)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if header == nil {
		s.returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("header not found: %d", headerID))
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

	err = s.templates.ExecuteTemplate(w, "accordion_header_form", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

func (s *Server) HandleAccordionHeaderDeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	headerID, err := strconv.Atoi(vars["header"])
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	header, err := s.modelClient.ReadAccordionHeader(headerID)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if header == nil {
		s.returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("header not found: %d", headerID))
		return
	}

	err = header.Delete(s.modelClient)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
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

func (s *Server) HandleAccordionHeaderEditGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Init template variables
	tmplVars := &AccordionHeaderFormTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	headerID, err := strconv.Atoi(vars["header"])
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	header, err := s.modelClient.ReadAccordionHeader(headerID)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if header == nil {
		s.returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("header not found: %d", headerID))
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

	err = s.templates.ExecuteTemplate(w, "accordion_header_form", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

func (s *Server) HandleAccordionHeaderEditPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	headerID, err := strconv.Atoi(vars["header"])
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	header, err := s.modelClient.ReadAccordionHeader(headerID)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if header == nil {
		s.returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("header not found: %d", headerID))
		return
	}

	// parse form data
	err = r.ParseForm()
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	header.Title = r.Form.Get("title")

	err = header.Update(s.modelClient)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
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

func (s *Server) HandleAccordionHeaderGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Init template variables
	tmplVars := &AccordionHeaderTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	headerID, err := strconv.Atoi(vars["header"])
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if headerID == 0 {
		tmplVars.Header = &database.AccordionHeader{
			ID:    0,
			Title: "The Hive",
		}

		tmplVars.Links, err = s.modelClient.ReadAccordionLinks(sql.NullInt32{Valid: false})
		if err != nil {
			s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		tmplVars.Header, err = s.modelClient.ReadAccordionHeader(headerID)
		if err != nil {
			s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		if tmplVars.Header == nil {
			s.returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("header %d not found", headerID))
			return
		}

		tmplVars.Links, err = s.modelClient.ReadAccordionLinks(sql.NullInt32{Valid: true, Int32: int32(tmplVars.Header.ID)})
		if err != nil {
			s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
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

	err = s.templates.ExecuteTemplate(w, "accordion_header_view", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}
