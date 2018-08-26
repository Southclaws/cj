package forum

import (
	"net/http"

	"github.com/pkg/errors"
	"gopkg.in/xmlpath.v2"
)

// GetHTMLRoot returns page content as goquery.Document
func (fc *ForumClient) GetHTMLRoot(url string) (*xmlpath.Node, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to GET %s", url)
	}

	response, err := fc.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to perform request for %s", url)
	}

	root, err := xmlpath.ParseHTML(response.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse HTML for %s", url)
	}

	return root, nil
}
