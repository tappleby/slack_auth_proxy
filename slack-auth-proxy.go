package main

import (
	"log"
	"net"
	"net/http"
	"strings"
	"github.com/tappleby/slack-auth-proxy/slack"
)

func main() {

	config, err := LoadConfiguration()

	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	log.Println(config.Upstreams[0].HostURL)

	listener, err := net.Listen("tcp", config.ServerAddr)
	if err != nil {
		log.Fatalf("FATAL: listen (%s) failed - %s", config.ServerAddr, err.Error())
	}
	log.Printf("listening on %s", config.ServerAddr)

	oauthClient := slack.NewOAuthClient(config.ClientId, config.ClientSecret, config.RedirectUri)
	oauthClient.TeamId = config.SlackTeam

	oauthServer := NewOauthServer(oauthClient, config.Upstreams)

	server := &http.Server{Handler: oauthServer}
	err = server.Serve(listener)
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
	log.Printf("ERROR: http.Serve() - %s", err.Error())
	}

	log.Printf("HTTP: closing %s", listener.Addr().String())
}
