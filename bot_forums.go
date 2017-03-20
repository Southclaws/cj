package main

import (
	"fmt"
	"strings"

	"strconv"

	"gopkg.in/xmlpath.v2"
)

// UserProfile stores data available on a user's profile page
type UserProfile struct {
	UserName        string
	JoinDate        string
	TotalPosts      int
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
	var result UserProfile

	root, err := app.GetHTMLRoot(url)
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
	var path *xmlpath.Path
	var result string

	path = xmlpath.MustCompile(`//*[@id="collapseobj_stats"]/div/*/ul/*[contains(.,'Join Date: ')]`)

	result, ok := path.String(root)
	if !ok {
		return result, fmt.Errorf("join date xmlpath did not return a result")
	}

	return strings.TrimPrefix(result, "Join Date: "), nil
}

// getTotalPosts returns the user total posts
func (app App) getTotalPosts(root *xmlpath.Node) (int, error) {
	path := xmlpath.MustCompile(`//*[@id="collapseobj_stats"]/div/fieldset[1]/ul/li[1]`)

	posts, ok := path.String(root)
	if !ok {
		return 0, fmt.Errorf("total posts xmlpath did not return a result")
	}

	posts = strings.TrimPrefix(posts, "Total Posts: ")
	posts = strings.Replace(posts, ",", "", -1)

	result, err := strconv.Atoi(posts)
	if err != nil {
		return 0, fmt.Errorf("cannot convert posts to integer")
	}

	return result, nil
}

// getTotalPosts returns the user total posts
func (app App) getReputation(forumUserID string) (int, error) {
	root, err := app.GetHTMLRoot(fmt.Sprintf("http://forum.sa-mp.com/search.php?do=finduser&u=%s", forumUserID))
	if err != nil {
		return 0, fmt.Errorf("cannot get user's posts")
	}

	path := xmlpath.MustCompile(`//td[@class="alt1"]/div[@class="alt2"]/div/em/a/@href`)

	// Get the first post from the list.
	href, ok := path.String(root)
	if !ok {
		return 0, fmt.Errorf("cannot get user posts")
	}

	// If we have a valid post, search in it for user's reputation.
	root, err = app.GetHTMLRoot(fmt.Sprintf("http://forum.sa-mp.com/%s", href))
	if err != nil {
		return 0, fmt.Errorf("cannot get user's post in a topic")
	}

	path = xmlpath.MustCompile(fmt.Sprintf(`//table[@id="%s"]/tbody/tr[@valign="top"]/td[@class="alt2"]/*/*[contains(text(),'Reputation: ')]`, strings.Split(href, "#")[1]))

	// Get the table for that post.
	reputation, ok := path.String(root)
	if !ok {
		return 0, fmt.Errorf("cannot get reputation field from post")
	}

	reputation = strings.TrimPrefix(reputation, "Reputation: ")
	reputation = strings.Replace(reputation, ",", "", -1)

	result, err := strconv.Atoi(reputation)
	if err != nil {
		return 0, fmt.Errorf("cannot convert reputation to integer")
	}

	return result, nil
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
		if !ok {
			continue
		}
		text, ok = textPath.String(messageBlock.Node())
		if !ok {
			continue
		}

		result = append(result, VisitorMessage{UserName: user, Message: text})
	}

	return result, nil
}
