package tracker

import "net/http"

///rest/issue?{project}&{summary}&{description}&{attachments}&{permittedGroup}

type youtrack struct {
}

func (u *youtrack) SetAccessToken(token string) {

	panic("implement me")
}

func (u *youtrack) BuildRequest(user, title, body, project string) (*http.Request, error) {
	panic("implement me")
}

func (u *youtrack) Create(r *http.Request) (string, error) {
	panic("implement me")
}
