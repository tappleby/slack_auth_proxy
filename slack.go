package main

import "net/http"

type SlackApi struct {
	httpClient *http.Client


	Token string
	Groups *GroupService
}

func NewSlackApi(token string) *SlackApi {
	api := &SlackApi{ Token: token, httpClient: http.DefaultClient }
	api.Groups = &GroupService{ api: api }

	return api;
}
