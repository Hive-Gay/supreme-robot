package webapp

import (
	"github.com/Hive-Gay/go-hivelib"
	"github.com/gorilla/sessions"
	"net/http"
)

type QuotesDashTemplate struct {
	templateCommon

	IsmList   *hivelib.IsmList
	SayerList *hivelib.SayerList
}

type QuotesIsmFormTemplate struct {
	templateCommon
	Breadcrumbs *[]templateBreadcrumb

	TitleText  string
	FormId     *templateFormInput
	FormText   *templateFormInput
	FormTts    *templateFormInput
	FormSubmit *templateFormButton
}

type QuotesIsmsTemplate struct {
	templateCommon
	Breadcrumbs *[]templateBreadcrumb

	IsmList *hivelib.IsmList
}

type QuotesSayerFormTemplate struct {
	templateCommon
	Breadcrumbs *[]templateBreadcrumb

	TitleText  string
	FormId     *templateFormInput
	FormUuid   *templateFormInput
	FormSubmit *templateFormButton
}

type QuotesSayersTemplate struct {
	templateCommon
	Breadcrumbs *[]templateBreadcrumb

	SayerList *hivelib.SayerList
}

func (s *Server) QuotesDashGetHandler(w http.ResponseWriter, r *http.Request) {
	// Init template variables
	tmplVars := &QuotesDashTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	tmplVars.IsmList, err = s.quotes.IsmsPageGet(1, 5, "created_at", false)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	tmplVars.SayerList, err = s.quotes.SayersPageGet(1, 5, "created_at", false)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	err = s.templates.ExecuteTemplate(w, "quotes_dashboard", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

func (s *Server) QuotesIsmAddGetHandler(w http.ResponseWriter, r *http.Request) {
	// Init template variables
	tmplVars := &QuotesIsmFormTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// breadcrumbs
	tmplVars.Breadcrumbs = &[]templateBreadcrumb{
		{
			HRef: "/app/quotes",
			Text: "Dashboard",
		},
		{
			Text: "Quotes",
			HRef: "/app/quotes/isms",
		},
		{
			Text: "Add Quote",
		},
	}

	// Configure Form
	tmplVars.TitleText = "Add Quote"
	tmplVars.FormId = &templateFormInput{
		ID:          "id",
		Name:        "id",
		Placeholder: "ID",
		Required:    true,
	}
	tmplVars.FormText = &templateFormInput{
		ID:          "text",
		Name:        "text",
		Placeholder: "Text",
		Required:    true,
	}
	tmplVars.FormTts = &templateFormInput{
		ID:          "tts",
		Name:        "tts",
		Placeholder: "TTS",
	}
	tmplVars.FormSubmit = &templateFormButton{
		Color: "success",
		Text:  "Add",
	}

	err = s.templates.ExecuteTemplate(w, "quotes_ism_form", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

func (s *Server) QuotesIsmAddPostHandler(w http.ResponseWriter, r *http.Request) {
	// parse form data
	err := r.ParseForm()
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// create hive object
	ism := hivelib.Ism{
		Id:   r.Form.Get("id"),
		Text: r.Form.Get("text"),
		Tts:  r.Form.Get("tts"),
	}

	logger.Debugf("ism %#v", &ism)

	err = s.quotes.IsmPost(&ism)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	us := r.Context().Value(SessionKey).(*sessions.Session)
	us.Values["page-alert-success"] = templateAlert{Text: "Quote added"}
	err = us.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect home
	http.Redirect(w, r, "/app/quotes/isms", http.StatusFound)
}

func (s *Server) QuotesIsmsGetHandler(w http.ResponseWriter, r *http.Request) {
	// Init template variables
	tmplVars := &QuotesIsmsTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	tmplVars.IsmList, err = s.quotes.IsmsGet()
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// breadcrumbs
	tmplVars.Breadcrumbs = &[]templateBreadcrumb{
		{
			HRef: "/app/quotes",
			Text: "Dashboard",
		},
		{
			Text: "Quotes",
		},
	}

	err = s.templates.ExecuteTemplate(w, "quotes_isms_view", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

func (s *Server) QuotesSayerAddGetHandler(w http.ResponseWriter, r *http.Request) {
	// Init template variables
	tmplVars := &QuotesSayerFormTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// breadcrumbs
	tmplVars.Breadcrumbs = &[]templateBreadcrumb{
		{
			HRef: "/app/quotes",
			Text: "Dashboard",
		},
		{
			Text: "Sayers",
			HRef: "/app/quotes/sayers",
		},
		{
			Text: "Add Sayer",
		},
	}

	// Configure Form
	tmplVars.TitleText = "Add Sayer"
	tmplVars.FormId = &templateFormInput{
		ID:          "id",
		Name:        "id",
		Placeholder: "ID",
		Required:    true,
	}
	tmplVars.FormUuid = &templateFormInput{
		ID:          "uuid",
		Name:        "uuid",
		Placeholder: "UUID",
	}
	tmplVars.FormSubmit = &templateFormButton{
		Color: "success",
		Text:  "Add",
	}

	err = s.templates.ExecuteTemplate(w, "quotes_sayer_form", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

func (s *Server) QuotesSayerAddPostHandler(w http.ResponseWriter, r *http.Request) {
	// parse form data
	err := r.ParseForm()
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// create hive object
	sayer := hivelib.Sayer{
		Id:   r.Form.Get("id"),
		Uuid: r.Form.Get("uuid"),
	}

	logger.Debugf("post %#v", &sayer)

	err = s.quotes.SayerPost(&sayer)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	us := r.Context().Value(SessionKey).(*sessions.Session)
	us.Values["page-alert-success"] = templateAlert{Text: "Sayer added"}
	err = us.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect home
	http.Redirect(w, r, "/app/quotes/sayers", http.StatusFound)
}

func (s *Server) QuotesSayersGetHandler(w http.ResponseWriter, r *http.Request) {
	// Init template variables
	tmplVars := &QuotesSayersTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	tmplVars.SayerList, err = s.quotes.SayersGet()
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// breadcrumbs
	tmplVars.Breadcrumbs = &[]templateBreadcrumb{
		{
			HRef: "/app/quotes",
			Text: "Dashboard",
		},
		{
			Text: "Sayers",
		},
	}

	err = s.templates.ExecuteTemplate(w, "quotes_sayers_view", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}
