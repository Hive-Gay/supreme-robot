package webapp

import (
	"net/http"
)

type HomeTemplate struct {
	templateCommon
}

func GetHome(w http.ResponseWriter, r *http.Request) {

	// Init template variables
	tmplVars := &HomeTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	err = templates.ExecuteTemplate(w, "home", tmplVars)
	if err != nil {
		logger.Errorf("could not render home template: %s", err.Error())
	}

}
