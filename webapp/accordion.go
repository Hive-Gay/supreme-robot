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

func (s *Server) HandleAccordion(w http.ResponseWriter, r *http.Request) {
	logger.Tracef("Starting HandleAccordion")
	// Init template variables
	tmplVars := &AccordionTemplate{}

	// Add The Hive to the Top
	theHive := AccordionHeader{
		Title: "The Hive",
	}

	hiveLinks, err := s.db.ReadAccordionLinks(sql.NullInt32{Valid: false})
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	for _, link := range hiveLinks {
		theHive.Links = append(theHive.Links, AccordionLink{
			Title: link.Title,
			Link:  link.Link,
		})
	}

	tmplVars.Accordion = []AccordionHeader{
		theHive,
	}

	// Get other headers
	mHeaders, err := s.db.ReadAccordionHeaders()
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	for _, mHeader := range mHeaders {
		header := AccordionHeader{
			Title: mHeader.Title,
		}

		headerLinks, err := s.db.ReadAccordionLinks(sql.NullInt32{Valid: true, Int32: int32(mHeader.ID)})
		if err != nil {
			s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		for _, link := range headerLinks {
			header.Links = append(header.Links, AccordionLink{
				Title: link.Title,
				Link:  link.Link,
			})
		}

		tmplVars.Accordion = append(tmplVars.Accordion, header)
	}

	err = s.templates.ExecuteTemplate(w, "accordion", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}
