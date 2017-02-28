package main

import (
	"fmt"
	"net/http"
	"strings"

	"strconv"

	"gopkg.in/xmlpath.v2"
)

// UserProfile stores data available on a user's profile page
type UserProfile struct {
	UserName        string
	JoinDate        string
	TotalPosts      string
	Reputation      int
	BioText         string
	VisitorMessages []VisitorMessage
	Errors          []error
}

// VisitorMessage represents a single visitor message available on a user page.
type VisitorMessage struct {
	UserName string
	Message  string
}

// GetUserProfilePage does a HTTP GET on the user's profile page then extracts
// structured information from it.
func (app App) GetUserProfilePage(url string) (UserProfile, error) {
	var err error
	var result UserProfile
	var req *http.Request
	var response *http.Response

	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return result, err
	}
	if response, err = app.httpClient.Do(req); err != nil {
		return result, err
	}

	root, err := xmlpath.ParseHTML(response.Body)
	if err != nil {
		return result, err
	}

	result.UserName, err = app.getUserName(root)
	if err != nil {
		return result, fmt.Errorf("url did not lead to a valid user page")
	}

	result.JoinDate, err = app.getJoinDate(root)
	if err != nil {
		result.Errors = append(result.Errors, err)
	}

	result.TotalPosts, err = app.getTotalPosts(root)
	if err != nil {
		result.Errors = append(result.Errors, err)
	}

	result.Reputation, err = app.getReputation(strings.TrimPrefix(url, "http://forum.sa-mp.com/member.php?u="))
	if err != nil {
		result.Errors = append(result.Errors, err)
	}

	result.BioText, err = app.getUserBio(root)
	if err != nil {
		result.Errors = append(result.Errors, err)
	}

	result.VisitorMessages, err = app.getFirstTenUserVisitorMessages(root)
	if err != nil {
		result.Errors = append(result.Errors, err)
	}

	return result, nil
}

// getUserName returns the user profile page owner name
func (app App) getUserName(root *xmlpath.Node) (string, error) {
	var result string

	path := xmlpath.MustCompile(`//*[@id="username_box"]/h1`)

	result, ok := path.String(root)
	if !ok {
		return result, fmt.Errorf("user name xmlpath did not return a result")
	}

	return strings.Trim(result, "\n "), nil
}

// getJoinDate returns the user join date
func (app App) getJoinDate(root *xmlpath.Node) (string, error) {
	var result string

	path := xmlpath.MustCompile(`//*[@id="collapseobj_stats"]/div/fieldset[3]/ul/li[2]`)

	result, ok := path.String(root)
	if !ok {
		return result, fmt.Errorf("user name xmlpath did not return a result")
	}

	return strings.TrimPrefix(result, "Join Date: "), nil
}

// getTotalPosts returns the user total posts
func (app App) getTotalPosts(root *xmlpath.Node) (string, error) {
	var result string

	path := xmlpath.MustCompile(`//*[@id="collapseobj_stats"]/div/fieldset[1]/ul/li[1]`)

	result, ok := path.String(root)
	if !ok {
		return result, fmt.Errorf("user name xmlpath did not return a result")
	}

	return strings.TrimPrefix(result, "Total Posts: "), nil
}

// getTotalPosts returns the user total posts
func (app App) getReputation(forumUserID string) (int, error) {
	document, ok := app.GetHTMLDocument(fmt.Sprintf("http://forum.sa-mp.com/search.php?do=finduser&u=%s", forumUserID))
	if !ok {
		return 0, fmt.Errorf("cannot get user's posts")
	}

	// Get the first post from the list.
	href := document.Find(`td[class="alt1"] > div[class="alt2"] > div > em > a`).AttrOr("href", "")

	// If we have a valid post, search in it for user's reputation.
	if len(href) > 0 {
		document, ok := app.GetHTMLDocument(fmt.Sprintf("http://forum.sa-mp.com/%s", href))
		if !ok {
			return 0, fmt.Errorf("cannot get user's post in a topic")
		}

		// Get the table for that post.
		reputation := document.Find(fmt.Sprintf(`table[id="%s"] > tbody > tr[valign="top"] > td[class="alt2"] > div`, strings.Split(href, "#")[1]))

		// Get "div" with index "4" from "td[class="alt2"]" and then get the "div" for reputation.
		return strconv.Atoi(strings.TrimPrefix(reputation.Eq(4).Find("div").Eq(2).Text(), "Reputation: "))
	}

	return 0, fmt.Errorf("cannot find a valid post for %s", forumUserID)
}

// getUserBio returns the bio text.
func (app App) getUserBio(root *xmlpath.Node) (string, error) {
	var result string

	path := xmlpath.MustCompile(`//*[@id="collapseobj_aboutme"]/div/ul/li[1]/dl/dd[1]`)

	result, ok := path.String(root)
	if !ok {
		return result, fmt.Errorf("user bio xmlpath did not return a result")
	}

	return result, nil
}

// getFirstTenUserVisitorMessages returns up to ten visitor messages from
func (app App) getFirstTenUserVisitorMessages(root *xmlpath.Node) ([]VisitorMessage, error) {
	var result []VisitorMessage

	mainPath := xmlpath.MustCompile(`//*[@id="message_list"]/*`)
	userPath := xmlpath.MustCompile(`.//div[2]/div[1]/div/a`)
	textPath := xmlpath.MustCompile(`.//div[2]/div[2]`)

	if !mainPath.Exists(root) {
		return result, fmt.Errorf("visitor messages xmlpath did not return a result")
	}

	var ok bool
	var user string
	var text string

	messageBlock := mainPath.Iter(root)

	for messageBlock.Next() {
		user, ok = userPath.String(messageBlock.Node())
		text, ok = textPath.String(messageBlock.Node())

		if !ok {
			continue
		}

		result = append(result, VisitorMessage{UserName: user, Message: text})
	}

	return result, nil
}
