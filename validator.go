package main

import (
	"github.com/tappleby/slack-auth-proxy/slack"
)

func NewValidator() func(*slack.Auth, *UpstreamConfiguration) bool {
	validator := func(auth *slack.Auth, upstream *UpstreamConfiguration) bool {
		return len(upstream.Users) == 0 || upstream.FindUsername(auth.Username) != ""
	}
	return validator
}
