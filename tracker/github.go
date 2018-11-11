package tracker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// CreateIssueRequest structure used for issue creation
type CreateIssueRequest struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Labels    []string `json:"labels"`
	Assignees []string `json:"assignees"`
	Milestone uint64   `json:"milestone,omitempty"`
}

// CreateIssueResponse structure
type CreateIssueResponse struct {
	Url string `json:"html_url"`
}

type creator struct {
	token     string
	endpoint  string // https://github.com/settings/tokens
	createURL string
	client    http.Client
}

const defaultEndpoint = "https://api.github.com"
const createIssueURL = "/repos/%s/%s/issues"

func NewGitIssueCreator() *creator {
	return &creator{}
}

// Endpoint setter for git repository
func (c *creator) Endpoint(ept string) {
	c.endpoint = ept
}

// Url assembler
func (c *creator) Url(endpoint, user, project string) string {
	if endpoint == "" {
		endpoint = defaultEndpoint
	}
	return endpoint + fmt.Sprintf(createIssueURL, user, project)
}

// SetAccessToken setter
func (c *creator) SetAccessToken(tkn string) {
	c.token = tkn
}

// AccessToken getter
func (c *creator) AccessToken() string {
	return c.token
}

// Request builder
func (c *creator) BuildRequest(user, title, body, project string) (*http.Request, error) {
	w := &bytes.Buffer{}
	data := CreateIssueRequest{
		Title:     title,
		Body:      body,
		Assignees: []string{},
		Labels:    []string{},
	}
	b, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("git.BuildRequest request structure error :" + err.Error())
	}
	log.Print("git.buildRequest body ", string(b))
	if _, err := w.Write(b); err != nil {
		return nil, errors.New("git.BuildRequest request write error :" + err.Error())
	}

	request, err := http.NewRequest(http.MethodPost, c.Url(c.endpoint, user, project), w)
	// all trackers have own headers //todo(ivanphotosoul).todo Move that part to another place
	request.Header.Add("Accept", "application/vnd.github.v3+json")
	request.Header.Add("Authorization", fmt.Sprintf("token %s", c.AccessToken()))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", "Awesome-Octocat-App")
	request.Header.Add("Access-Control-Allow-Origin", "*")

	return request, err
}

// Create issue
func (c *creator) Create(r *http.Request) (issue string, err error) {
	log.Printf("r.RequestURI %#v", r.RequestURI)
	log.Printf("r.Host %#v", r.Host)
	log.Printf("r.Header: %#v", r.Header)

	resp, err := c.client.Do(r)
	if err != nil {
		return "", errors.New("git.create: error after request: " + err.Error())
	}
	defer resp.Body.Close()

	log.Print("git.response.status ", resp.StatusCode)
	log.Print(r.RequestURI)
	if resp.StatusCode != http.StatusCreated {
		if resp.ContentLength < 1 {
			return "", errors.New("git.create: response body is empty")
		}
		buff := bytes.NewBuffer(make([]byte, resp.ContentLength))
		buff.ReadFrom(resp.Body)

		if err != nil {
			return issue, errors.New("git.create: could not read response body: " + err.Error())
		}
		return issue, errors.New("git.create: response body: " + buff.String())
	}

	return resp.Header.Get("Location"), err
}
