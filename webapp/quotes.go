package webapp

import (
	"github.com/Hive-Gay/go-hivelib"
	"net/http"
)

type QuotesDashTemplate struct {
	templateCommon

	IsmList   *hivelib.IsmList
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

	tmplVars.IsmList, err = s.quotes.IsmsGet()
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	tmplVars.SayerList, err = s.quotes.SayersGet()
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	err = s.templates.ExecuteTemplate(w, "quotes_dashboard", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}
