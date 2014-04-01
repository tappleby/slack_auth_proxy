package main

import (
	"log"
	"github.com/tappleby/slack-auth-proxy/slack"
)

func main() {
	oauthClient := slack.NewOAuthClient("foo", "bar", "http://127.0.0.1:4180/oauth2/callback")

	loginUrl := oauthClient.LoginUrl("")
	log.Println(loginUrl.String())

	authToken, err := oauthClient.RedeemCode("2176048118.2259982886.06aedf88fc");

	if err != nil {
		panic(err)
	}

	log.Println(authToken)
}
