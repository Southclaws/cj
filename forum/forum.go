package forum

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	scraper "github.com/Southclaws/go-cloudflare-scraper"
	"github.com/pkg/errors"
	"gopkg.in/xmlpath.v2"
)

// ForumClient provides programatic access to SA:MP forum
type ForumClient struct {
	httpClient *http.Client
}

// UserProfile stores data available on a user's profile page
type UserProfile struct {
	UserName        string
	JoinDate        string
	TotalPosts      int
	Reputation      int
	BioText         string
	DiscordID       string
	VisitorMessages []VisitorMessage
	Errors          []error
}

// VisitorMessage represents a single visitor message available on a user page.
type VisitorMessage struct {
	UserName string
	Message  string
}

// NewForumClient creates a new forum clienta
func NewForumClient() (fc *ForumClient, err error) {
	fc = new(ForumClient)
	scrpr, err := scraper.NewTransport(http.DefaultTransport)
	if err != nil {
		return
	}
	fc.httpClient = &http.Client{Transport: scrpr}
	return
}

// GetUserProfilePage does a HTTP GET on the user's profile page then extracts
// structured information from it.
func (fc *ForumClient) GetUserProfilePage(url string) (UserProfile, error) {
	var result UserProfile

	root, err := fc.GetHTMLRoot(url)
	if err != nil {
		return result, errors.Wrap(err, "failed to get HTML root for user page")
	}

	result.UserName, err = fc.getUserName(root)
	if err != nil {
		return result, errors.Wrap(err, "url did not lead to a valid user page")
	}

	result.JoinDate, err = fc.getJoinDate(root)
	if err != nil {
		result.Errors = append(result.Errors, err)
	}

	result.TotalPosts, err = fc.getTotalPosts(root)
	if err != nil {
		result.Errors = append(result.Errors, err)
	}

	result.Reputation, err = fc.getReputation(root)
	if err != nil {
		result.Errors = append(result.Errors, err)
	}

	result.BioText, err = fc.getUserBio(root)
	if err != nil {
		result.Errors = append(result.Errors, err)
	}

	result.VisitorMessages, err = fc.getFirstTenUserVisitorMessages(root)
	if err != nil {
		result.Errors = append(result.Errors, err)
	}

	result.DiscordID, err = fc.getUserDiscordID(root)
	if err != nil {
		result.Errors = append(result.Errors, err)
	}
	return result, nil
}

// getUserName returns the user profile page owner name
func (fc *ForumClient) getUserName(root *xmlpath.Node) (string, error) {
	var result string

	path := xmlpath.MustCompile(`//span[@class='largetext']//strong//span//strong`)

	result, ok := path.String(root)
	if !ok {
		return result, errors.New("user name xmlpath did not return a result")
	}

	return strings.Trim(result, "\n "), nil
}

// getJoinDate returns the user join date
func (fc *ForumClient) getJoinDate(root *xmlpath.Node) (string, error) {
	var path *xmlpath.Path
	var result string

	path = xmlpath.MustCompile(`//table[@id='profile_desktop']//td[@class='trow1' and contains(text(),'-')]`)

	result, ok := path.String(root)
	if !ok {
		return result, errors.New("join date xmlpath did not return a result")
	}

	return result, nil
}

// getTotalPosts returns the user total posts
func (fc *ForumClient) getTotalPosts(root *xmlpath.Node) (int, error) {
	path := xmlpath.MustCompile(`//table[@id='profile_desktop']//td[@class='trow1' and contains(text(),'posts')]`)

	posts, ok := path.String(root)
	if !ok {
		return 0, errors.New("total posts xmlpath did not return a result")
	}

	posts = strings.Split(posts, " ")[0]
	result, err := strconv.Atoi(posts)
	if err != nil {
		return 0, errors.New("cannot convert posts to integer")
	}

	return result, nil
}

// getTotalPosts returns the user total posts
func (fc *ForumClient) getReputation(root *xmlpath.Node) (int, error) {
	path := xmlpath.MustCompile(`//table[@id='profile_desktop']//td[@class='trow2']//strong[@class='reputation_positive']`)

	reputation, ok := path.String(root)
	if !ok {
		return 0, errors.New("get reputation xmlpath did not return a sresult")
	}

	result, err := strconv.Atoi(reputation)
	if err != nil {
		return 0, errors.Wrap(err, "cannot convert reputation to integer")
	}

	return result, nil
}

// getUserBio returns the bio text.
func (fc *ForumClient) getUserBio(root *xmlpath.Node) (string, error) {
	var result string

	path := xmlpath.MustCompile(`//*[@id="collapseobj_aboutme"]/div/ul/li[1]/dl/dd[1]`)

	result, ok := path.String(root)
	if !ok {
		return result, errors.New("user bio xmlpath did not return a result")
	}

	return result, nil
}

func (fc *ForumClient) getUserDiscordID(root *xmlpath.Node) (string, error) {
	var result string

	path := xmlpath.MustCompile(`//table[@id='profile_desktop']//td[@class='trow1 scaleimages']`)

	result, ok := path.String(root)
	if !ok {
		return result, errors.New("user discord id xmlpath did not resturn a result")
	}

	return result, nil
}

// getFirstTenUserVisitorMessages returns up to ten visitor messages from
func (fc *ForumClient) getFirstTenUserVisitorMessages(root *xmlpath.Node) (result []VisitorMessage, err error) {
	mainPath := xmlpath.MustCompile(`//*[@id="message_list"]/*`)
	userPath := xmlpath.MustCompile(`.//div[2]/div[1]/div/a`)
	textPath := xmlpath.MustCompile(`.//div[2]/div[2]`)

	if !mainPath.Exists(root) {
		return result, errors.New("visitor messages xmlpath did not return a result")
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

// NewPostAlert will call `fn` when the specified user posts or edits a post
func (fc *ForumClient) NewPostAlert(id string, fn func()) {
	ticker := time.NewTicker(time.Second * 10)
	lastPostCount := -1

	go func() {
		for range ticker.C {
			profile, err := fc.GetUserProfilePage("http://forum.sa-mp.com/member.php?u=" + id)
			if err != nil {
				return
			}

			if lastPostCount == -1 {
				lastPostCount = profile.TotalPosts
				continue
			}

			if lastPostCount < profile.TotalPosts {
				fn()
				lastPostCount = profile.TotalPosts
			}
		}
	}()
}
