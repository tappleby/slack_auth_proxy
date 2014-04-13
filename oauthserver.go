package main

import (
	"log"
	"github.com/tappleby/slack-auth-proxy/slack"
	"net/http"
	"fmt"
	"net/http/httputil"
	"strings"
)

const signInPath = "/oauth2/sign_in"
const oauthStartPath = "/oauth2/start"
const oauthCallbackPath = "/oauth2/callback"


type OAuthServer struct {
	slackOauth *slack.OAuthClient
	serveMux	*http.ServeMux

	upstreamsConfig UpstreamConfigurationMap
}

func NewOauthServer(slackOauth *slack.OAuthClient, upstreams []*UpstreamConfiguration) *OAuthServer {
	serveMux := http.NewServeMux()

	upstreamsPathMap := make(UpstreamConfigurationMap)

	for _, upstream := range upstreams {
		u := upstream.HostURL
		path := u.Path
		u.Path = ""

		if path == "" {
			path = "/"
		}

		log.Printf("mapping %s => %s", path, u)
		serveMux.Handle(path, httputil.NewSingleHostReverseProxy(u))

		upstreamsPathMap[path] = upstream
	}

	return &OAuthServer{
		serveMux: serveMux,
		slackOauth: slackOauth,
		upstreamsConfig: upstreamsPathMap,
	}
}

func (s *OAuthServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// check if this is a redirect back at the end of oauth
	remoteIP := req.Header.Get("X-Real-IP")
	if remoteIP == "" {
		remoteIP = req.RemoteAddr
	}
	log.Printf("%s %s %s", remoteIP, req.Method, req.URL.Path)

	if req.URL.Path == signInPath {
		s.handleSignIn(rw, req)
		return
	} else if req.URL.Path == oauthStartPath {
		s.handleOAuthStart(rw, req)
		return
	} else if (req.URL.Path == oauthCallbackPath) {
		s.handleOAuthCallback(rw, req)
		return
	}

	handler, pattern := s.serveMux.Handler(req)

 	upstreamConfig := s.upstreamsConfig[pattern]

	if upstreamConfig == nil {
		pattern = strings.TrimPrefix(pattern, "/")
		upstreamConfig = s.upstreamsConfig[pattern]
	}

	if upstreamConfig != nil {
		log.Println(upstreamConfig)
	} else {
		handler = http.NotFoundHandler()
	}

	handler.ServeHTTP(rw, req)
}

func (s *OAuthServer) handleSignIn(rw http.ResponseWriter, req *http.Request) {

}

func (s *OAuthServer) handleOAuthStart(rw http.ResponseWriter, req *http.Request) {
	http.Redirect(rw, req, s.slackOauth.LoginUrl("").String(), 302)
}

func (s *OAuthServer) handleOAuthCallback(rw http.ResponseWriter, req *http.Request) {
	// finish the oauth cycle
	err := req.ParseForm()
	if err != nil {
		s.ErrorPage(rw, 500, "Internal Error", err.Error())
		return
	}
	errorString := req.Form.Get("error")
	if errorString != "" {
		s.ErrorPage(rw, 403, "Permission Denied", errorString)
		return
	}

	access, err := s.slackOauth.RedeemCode(req.Form.Get("code"))

	if err != nil {
		log.Printf("error redeeming code %s", err.Error())
		s.ErrorPage(rw, 500, "Internal Error", err.Error())
		return
	}

	cl := slack.NewClient(access.Token)
	auth, err := cl.Auth.Test()

	if err != nil {
		log.Printf("error redeeming code %s", err.Error())
		s.ErrorPage(rw, 500, "Internal Error", err.Error())
		return
	}

	log.Println(auth)

	fmt.Fprintln(rw, "OK")
}

func (p *OAuthServer) ErrorPage(rw http.ResponseWriter, code int, title string, message string) {
	log.Printf("ErrorPage %d %s %s", code, title, message)
	rw.WriteHeader(code)
	fmt.Fprintln(rw, message)
}
