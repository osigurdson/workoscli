package workoscli

import (
	"fmt"

	"golang.org/x/oauth2"
)

type WorkOSConf struct {
	AuthKitURI   string
	RedirectPort int
	RedirectIP   string
	ClientID     string
	Scopes       []string
}

func (c *WorkOSConf) RedirectURL() string {
	return fmt.Sprintf("http://%s:%d/callback", c.RedirectIP, c.RedirectPort)
}

func (c *WorkOSConf) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.RedirectIP, c.RedirectPort)
}

func (c *WorkOSConf) toOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:    c.ClientID,
		RedirectURL: c.RedirectURL(),
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.AuthKitURI + "/oauth2/authorize",
			TokenURL: c.AuthKitURI + "/oauth2/token",
		},
		Scopes: c.Scopes,
	}
}
