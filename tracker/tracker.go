package tracker

import (
	"errors"
	"net/http"
	"text/template"
)

type CreatorType string

const (
	GitCreator      CreatorType = "git"
	YoutrackCreator CreatorType = "youtrack"
)

// Creator represent interface for task creation
type Creator interface {
	Create(title, body string) error
}

// IssueCreator interface
type IssueCreator interface {
	SetAccessToken(token string)
	BuildRequest(user, title, body, project string) (*http.Request, error)
	Create(r *http.Request) (string, error)
}

// Issue common structure
type Issue struct {
	username, kind, project, description, title string
}

// IssueTemplate
type IssueTemplate struct {
	title template.Template
	body  template.Template
}

// Create git issue
func Create(title, body, user, project string, ic IssueCreator) (string, error) {
	if ic == nil {
		return "", errors.New("issueCreator is nil")
	}
	r, err := ic.BuildRequest(user, title, body, project)
	if err != nil {
		return "", err
	}
	return ic.Create(r)
}

// New issue creator with given type
func New(kind CreatorType) IssueCreator {
	if kind == GitCreator {
		return NewGitIssueCreator()
	}
	return nil
}
