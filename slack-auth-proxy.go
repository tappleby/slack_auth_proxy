package main

import "log"
import (
	"github.com/tappleby/slack-auth-proxy/slack"
)

func main() {

	slackClient := slack.NewSlackApi("xoxp-2176048118-2176048120-2250552618-570941")

	userAuth, _ :=  slackClient.Auth.Test()

	log.Println(userAuth)
}
