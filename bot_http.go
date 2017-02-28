package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// GetHTMLDocument returns page content as goquery.Document
func (app App) GetHTMLDocument(url string) (*goquery.Document, bool) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("NewRequest error: %s", err.Error())
		return nil, false
	}

	response, err := app.httpClient.Do(request)
	if err != nil {
		log.Printf("httpClient.Do error: %s", err.Error())
		return nil, false
	}

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Printf("NewDocumentFromReader error: %s", err.Error())
		return nil, false
	}

	return document, true
}
