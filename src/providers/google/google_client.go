package google

import (
	"encoding/json"

	"github.com/scoville/scvl/src/domain"
	"github.com/scoville/scvl/src/engine"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleClient struct {
	config *oauth2.Config
}

// NewClient returns the googleClient
func NewClient(clientID, clientSecret, redirectURL string) engine.GoogleClient {
	return &googleClient{
		&oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/userinfo.email",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (c *googleClient) FetchUserInfo(code string) (user domain.User, err error) {
	tok, err := c.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return
	}
	client := c.config.Client(oauth2.NoContext, tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&user)
	return
}

func (c *googleClient) AuthCodeURL(state string) string {
	return c.config.AuthCodeURL(state)
}
