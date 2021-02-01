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

type AccordionLinkFormTemplate struct {
	templateCommon
	Breadcrumbs *[]templateBreadcrumb

	Header *models.AccordionHeader

	TitleText              string
	FormInputTitleDisabled bool
	FormInputTitleValue    string
	FormInputLinkDisabled  bool
	FormInputLinkValue     string
	FormButtonSubmitColor  string
	FormButtonSubmitText   string
}

func (s *Server)HandleAccordionLinkAddGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Init template variables
	tmplVars := &AccordionLinkFormTemplate{}
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
		tmplVars.Header = &models.AccordionHeader{
			ID:    0,
			Title: "The Hive",
		}
	} else {
		tmplVars.Header, err = s.modelClient.ReadAccordionHeader(headerID)
		if err != nil {
			s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		if tmplVars.Header == nil {
			s.returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("header not found: %d", headerID))
			return
		}
	}

	// Configure Form
	tmplVars.TitleText = fmt.Sprintf("Add Link for %s", tmplVars.Header.Title)
	tmplVars.FormButtonSubmitColor = "success"
	tmplVars.FormButtonSubmitText = "Add"

	// breadcrumbs
	tmplVars.Breadcrumbs = &[]templateBreadcrumb{
		{
			HRef: "/app/accordion",
			Text: "Accordion",
		},
		{
			HRef: fmt.Sprintf("/app/accordion/%d", tmplVars.Header.ID),
			Text: tmplVars.Header.Title,
		},
		{
			Text: tmplVars.TitleText,
		},
	}

	err = s.templates.ExecuteTemplate(w, "accordion_link_form", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

func (s *Server)HandleAccordionLinkAddPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	headerID, err := strconv.Atoi(vars["header"])
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// parse form data
	err = r.ParseForm()
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	al := models.AccordionLink{
		Title: r.Form.Get("title"),
		Link: r.Form.Get("link"),
	}
	if headerID == 0 {
		al.AccordionHeaderID = sql.NullInt32{Valid: false}
	} else {
		al.AccordionHeaderID = sql.NullInt32{Valid: true, Int32: int32(headerID)}
	}

	err = s.modelClient.CreateAccordionLink(&al)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	us := r.Context().Value(SessionKey).(*sessions.Session)
	us.Values["page-alert-success"] = templateAlert{Text: "Link added"}
	err = us.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect home
	http.Redirect(w, r, fmt.Sprintf("/app/accordion/%d", headerID), http.StatusFound)
}

func (s *Server)HandleAccordionLinkDeleteGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Init template variables
	tmplVars := &AccordionLinkFormTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// Get header
	headerID, err := strconv.Atoi(vars["header"])
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	headerIDSQL := sql.NullInt32{Valid: false}

	if headerID == 0 {
		tmplVars.Header = &models.AccordionHeader{
			ID:    0,
			Title: "The Hive",
		}
	} else {
		headerIDSQL.Valid = true
		headerIDSQL.Int32 = int32(headerID)

		tmplVars.Header, err = s.modelClient.ReadAccordionHeader(headerID)
		if err != nil {
			s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		if tmplVars.Header == nil {
			s.returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("header not found: %d", headerID))
			return
		}
	}

	// Get Link
	linkID, err := strconv.Atoi(vars["link"])
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	link, err := s.modelClient.ReadAccordionLink(headerIDSQL, linkID)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if link == nil {
		s.returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("link not found: %d", headerID))
		return
	}

	// Configure Form
	tmplVars.FormInputTitleValue = link.Title
	tmplVars.FormInputTitleDisabled = true
	tmplVars.FormInputLinkValue = link.Link
	tmplVars.FormInputLinkDisabled = true

	tmplVars.TitleText = fmt.Sprintf("Delete Link for %s", tmplVars.Header.Title)
	tmplVars.FormButtonSubmitColor = "danger"
	tmplVars.FormButtonSubmitText = "Delete"

	// breadcrumbs
	tmplVars.Breadcrumbs = &[]templateBreadcrumb{
		{
			HRef: "/app/accordion",
			Text: "Accordion",
		},
		{
			HRef: fmt.Sprintf("/app/accordion/%d", tmplVars.Header.ID),
			Text: tmplVars.Header.Title,
		},
		{
			Text: tmplVars.TitleText,
		},
	}

	err = s.templates.ExecuteTemplate(w, "accordion_link_form", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

func (s *Server)HandleAccordionLinkDeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Get header
	headerID, err := strconv.Atoi(vars["header"])
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	headerIDSQL := sql.NullInt32{Valid: false}

	if headerID != 0 {
		headerIDSQL.Valid = true
		headerIDSQL.Int32 = int32(headerID)

		header, err := s.modelClient.ReadAccordionHeader(headerID)
		if err != nil {
			s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		if header == nil {
			s.returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("header not found: %d", headerID))
			return
		}
	}

	// Get Link
	linkID, err := strconv.Atoi(vars["link"])
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	link, err := s.modelClient.ReadAccordionLink(headerIDSQL, linkID)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if link == nil {
		s.returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("link not found: %d", headerID))
		return
	}

	// Delete link
	err = s.modelClient.DeleteAccordionLink(linkID)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	us := r.Context().Value(SessionKey).(*sessions.Session)
	us.Values["page-alert-success"] = templateAlert{Text: "Link updated"}
	err = us.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect
	http.Redirect(w, r, fmt.Sprintf("/app/accordion/%d", headerID), http.StatusFound)

}

func (s *Server)HandleAccordionLinkEditGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Init template variables
	tmplVars := &AccordionLinkFormTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// Get header
	headerID, err := strconv.Atoi(vars["header"])
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	headerIDSQL := sql.NullInt32{Valid: false}

	if headerID == 0 {
		tmplVars.Header = &models.AccordionHeader{
			ID:    0,
			Title: "The Hive",
		}
	} else {
		headerIDSQL.Valid = true
		headerIDSQL.Int32 = int32(headerID)

		tmplVars.Header, err = s.modelClient.ReadAccordionHeader(headerID)
		if err != nil {
			s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		if tmplVars.Header == nil {
			s.returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("header not found: %d", headerID))
			return
		}
	}

	// Get Link
	linkID, err := strconv.Atoi(vars["link"])
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	link, err := s.modelClient.ReadAccordionLink(headerIDSQL, linkID)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if link == nil {
		s.returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("link not found: %d", headerID))
		return
	}

	// Configure Form
	tmplVars.FormInputTitleValue = link.Title
	tmplVars.FormInputLinkValue = link.Link

	tmplVars.TitleText = fmt.Sprintf("Edit Link for %s", tmplVars.Header.Title)
	tmplVars.FormButtonSubmitColor = "warning"
	tmplVars.FormButtonSubmitText = "Update"

	// breadcrumbs
	tmplVars.Breadcrumbs = &[]templateBreadcrumb{
		{
			HRef: "/app/accordion",
			Text: "Accordion",
		},
		{
			HRef: fmt.Sprintf("/app/accordion/%d", tmplVars.Header.ID),
			Text: tmplVars.Header.Title,
		},
		{
			Text: tmplVars.TitleText,
		},
	}

	err = s.templates.ExecuteTemplate(w, "accordion_link_form", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

func (s *Server)HandleAccordionLinkEditPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Get header
	headerID, err := strconv.Atoi(vars["header"])
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	headerIDSQL := sql.NullInt32{Valid: false}

	if headerID != 0 {
		headerIDSQL.Valid = true
		headerIDSQL.Int32 = int32(headerID)

		header, err := s.modelClient.ReadAccordionHeader(headerID)
		if err != nil {
			s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		if header == nil {
			s.returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("header not found: %d", headerID))
			return
		}
	}

	// Get Link
	linkID, err := strconv.Atoi(vars["link"])
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	link, err := s.modelClient.ReadAccordionLink(headerIDSQL, linkID)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if link == nil {
		s.returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("link not found: %d", headerID))
		return
	}

	// parse form data
	err = r.ParseForm()
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	link.Title = r.Form.Get("title")
	link.Link = r.Form.Get("link")

	err = s.modelClient.UpdateAccordionLink(link)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	us := r.Context().Value(SessionKey).(*sessions.Session)
	us.Values["page-alert-success"] = templateAlert{Text: "Link updated"}
	err = us.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect
	http.Redirect(w, r, fmt.Sprintf("/app/accordion/%d", headerID), http.StatusFound)

}