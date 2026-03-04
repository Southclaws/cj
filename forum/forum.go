package forum

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/xmlpath.v2"
)

// ForumClient provides programatic access to SA:MP forum
type ForumClient struct {
	httpClient *http.Client
}

// UserProfile stores data available on a user's profile page
type UserProfile struct {
	UserName   string
	JoinDate   string
	TotalPosts int
	Reputation int
	DiscordID  string
	Errors     []error
}

// NewForumClient creates a new forum clienta
func NewForumClient() (fc *ForumClient, err error) {
	fc = new(ForumClient)
	fc.httpClient = &http.Client{}
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

	result.DiscordID, err = fc.getUserDiscordID(root)
	if err != nil {
		result.Errors = append(result.Errors, err)
	}
	return result, nil
}

// getUserName returns the user profile page owner name
func (fc *ForumClient) getUserName(root *xmlpath.Node) (string, error) {
	paths := []*xmlpath.Path{
		xmlpath.MustCompile(`//span[@class='largetext']//strong//span`),
		xmlpath.MustCompile(`//span[@class='largetext']//strong`),
		xmlpath.MustCompile(`//span[@class='active']`),
		xmlpath.MustCompile(`//div[contains(@class,'profile-side')]`),
		xmlpath.MustCompile(`//title`),
	}

	for _, path := range paths {
		result, ok := path.String(root)
		if !ok {
			continue
		}

		result = strings.TrimSpace(result)
		if result == "" {
			continue
		}

		if strings.Contains(result, "Profile of ") {
			result = strings.TrimSpace(strings.TrimPrefix(result, "Profile of "))
		} else if strings.Contains(result, " - Profile of ") {
			parts := strings.SplitN(result, " - Profile of ", 2)
			if len(parts) == 2 {
				result = strings.TrimSpace(parts[1])
			}
		} else if strings.Contains(result, "'s Forum Info") {
			result = strings.TrimSpace(strings.TrimSuffix(result, "'s Forum Info"))
		} else if strings.Contains(result, "'s Profile") {
			result = strings.TrimSpace(strings.TrimSuffix(result, "'s Profile"))
		}

		if strings.Contains(result, " - ") {
			parts := strings.Split(result, " - ")
			if len(parts) > 0 {
				result = strings.TrimSpace(parts[0])
			}
		}

		if result != "" && result != "open.mp forum" {
			return result, nil
		}
	}

	return "", errors.New("user name xmlpath did not return a result")
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

func (fc *ForumClient) getUserDiscordID(root *xmlpath.Node) (string, error) {
	var result string

	path := xmlpath.MustCompile(`//table[@id='profile_desktop']//td[@class='trow1 scaleimages']`)

	result, ok := path.String(root)
	if !ok {
		return result, errors.New("user discord id xmlpath did not resturn a result")
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
