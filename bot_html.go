package main

import (
	"net/http"

	xmlpath "gopkg.in/xmlpath.v2"
)

// GetHTMLRoot returns page content as goquery.Document
func (app App) GetHTMLRoot(url string) (*xmlpath.Node, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := app.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	root, err := xmlpath.ParseHTML(response.Body)
	if err != nil {
		return nil, err
	}

	return root, nil
}
