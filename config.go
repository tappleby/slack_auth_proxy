package main

import (
	yaml "gopkg.in/yaml.v1"
	"io/ioutil"
	"path/filepath"
	"fmt"
)

const (
	defaultServerAddr = "127.0.0.1:4180"
)



type Configuration struct {
	// Server settings
	ServerAddr		string 						`yaml:"server_addr,omitempty"`
	Upstreams		[]*UpstreamConfiguration 	`yaml:"upstreams,omitempty"`
	RedirectUri		string						`yaml:"redirect_uri,omitempty"`

	// Slack Settings
	ClientId 		string 						`yaml:"client_id"`
	ClientSecret 	string 						`yaml:"client_secret"`
	SlackTeam 		string 					 	`yaml:"team_id,omitempty"`
	AuthToken		string 					 	`yaml:"auth_token,omitempty"`

}

func LoadConfiguration() (config *Configuration, err error) {

	configPath, _ := filepath.Abs("./config.yaml")
	configBuf, err := ioutil.ReadFile(configPath)

	if err != nil {
		err = fmt.Errorf("Failed to read configuration %s: %v", configPath, err)
		return
	}

	config = new(Configuration)


	if err = yaml.Unmarshal(configBuf, &config); err != nil {
		return
	}

	if config.ServerAddr == "" {
		config.ServerAddr = defaultServerAddr
	}

	if config.ClientId == "" {
		err = fmt.Errorf("Client id must be set in configuration")
		return
	}

	if config.ClientSecret == "" {
		err = fmt.Errorf("Client secret must be set in configuration")
		return
	}

	if config.RedirectUri == "" {
		config.RedirectUri = fmt.Sprintf("http://%s%s", config.ServerAddr, oauthCallbackPath)
	}

	for _, upstream := range config.Upstreams {
		if err = upstream.Parse(); err != nil {
			err = fmt.Errorf("Error parsing upstream %s: %v", upstream.Host, err)
			return
		}
	}


	return
}
