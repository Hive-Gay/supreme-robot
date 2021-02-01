package webapp

import (
	"database/sql"
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

func HandleAccordion(w http.ResponseWriter, r *http.Request) {
	// Init template variables
	tmplVars := &AccordionTemplate{}

	// Add The Hive to the Top
	theHive := AccordionHeader{
		Title: "The Hive",
	}

	hiveLinks, err := modelClient.ReadAccordionLinks(sql.NullInt32{Valid: false})
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	for _, link := range hiveLinks {
		theHive.Links = append(theHive.Links, AccordionLink{
			Title: link.Title,
			Link: link.Link,
		})
	}

	tmplVars.Accordion = []AccordionHeader{
		theHive,
	}

	// Get other headers
	mHeaders, err := modelClient.ReadAccordionHeaders()
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	for _, mHeader := range mHeaders {
		header := AccordionHeader{
			Title: mHeader.Title,
		}

		headerLinks, err := modelClient.ReadAccordionLinks(sql.NullInt32{Valid: true, Int32: int32(mHeader.ID)})
		if err != nil {
			returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		for _, link := range headerLinks {
			header.Links = append(header.Links, AccordionLink{
				Title: link.Title,
				Link: link.Link,
			})
		}

		tmplVars.Accordion = append(tmplVars.Accordion, header)
	}




	err = templates.ExecuteTemplate(w, "accordion", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}
