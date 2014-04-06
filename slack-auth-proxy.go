package main

import (
	"log"
	"net"
	"net/http"
	"strings"
	"flag"
	"github.com/tappleby/slack-auth-proxy/slack"
)

var (
	httpAddr                = flag.String("http-address", "127.0.0.1:4180", "<addr>:<port> to listen on for HTTP clients")
	clientID                = flag.String("client-id", "", "Slack Oauth client id")
	clientSecret            = flag.String("client-secret", "", "Slack oauth client secret")
	slackTeamId		        = flag.String("slack-team", "", "authenticate against the given slack team id")
)


func main() {

	flag.Parse()

	if *clientID == "" {
		log.Fatal("missing --client-id")
	}
	if *clientSecret == "" {
		log.Fatal("missing --client-secret")
	}

	listener, err := net.Listen("tcp", *httpAddr)
	if err != nil {
		log.Fatalf("FATAL: listen (%s) failed - %s", *httpAddr, err.Error())
	}
	log.Printf("listening on %s", *httpAddr)

	oauthClient := slack.NewOAuthClient(*clientID, *clientSecret, "http://127.0.0.1:4180/oauth2/callback")

	if *slackTeamId != "" {
		oauthClient.TeamId = *slackTeamId
	}

	oauthServer := NewOauthServer(oauthClient)

	server := &http.Server{Handler: oauthServer}
	err = server.Serve(listener)
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
	log.Printf("ERROR: http.Serve() - %s", err.Error())
	}

	log.Printf("HTTP: closing %s", listener.Addr().String())
}
