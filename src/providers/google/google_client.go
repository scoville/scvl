package google

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/Iwark/spreadsheet"
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
				"https://www.googleapis.com/auth/drive.metadata.readonly",
				spreadsheet.Scope,
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
	if err != nil {
		return
	}
	token, err := json.Marshal(tok)
	if err != nil {
		return
	}
	user.GoogleToken = string(token)
	return
}

func (c *googleClient) GetDriveFileTitle(user *domain.User, id string) (title string, err error) {
	tok := &oauth2.Token{}
	err = json.Unmarshal([]byte(user.GoogleToken), tok)
	if err != nil {
		return
	}
	client := c.config.Client(oauth2.NoContext, tok)
	resp, err := client.Get("https://www.googleapis.com/drive/v3/files/" + id)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var m map[string]string
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &m)
	if err != nil {
		log.Println(string(b))
	}
	title = m["name"]
	return
}

func (c *googleClient) FetchSpreadsheet(user *domain.User, id string) (ss spreadsheet.Spreadsheet, err error) {
	tok := &oauth2.Token{}
	err = json.Unmarshal([]byte(user.GoogleToken), tok)
	if err != nil {
		return
	}
	client := c.config.Client(oauth2.NoContext, tok)

	service := spreadsheet.NewServiceWithClient(client)
	return service.FetchSpreadsheet(id)
}

func (c *googleClient) AuthCodeURL(state string) string {
	return c.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}
