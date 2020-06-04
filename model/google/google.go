package oauthgoogle

import (
	"encoding/json"
	"io"
	"strings"
	"github.com/mattermost/mattermost-server/v5/einterfaces"
	"github.com/mattermost/mattermost-server/v5/model"
)

type GoogleProvider struct {
}


type GoogleUser struct {
	ResourceName string `json:"resourceName"`
	Etag         string `json:"etag"`
	Names        []struct {
		NamesMetadata struct {
			Primary bool `json:"primary"`
			Source  struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			} `json:"source"`
		} `json:"metadata,omitempty"`
		DisplayName          string `json:"displayName"`
		FamilyName           string `json:"familyName"`
		GivenName            string `json:"givenName"`
		DisplayNameLastFirst string `json:"displayNameLastFirst"`
		Metadata             struct {
			Source struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			} `json:"source"`
		} `json:"metadata,omitempty"`
	} `json:"names"`
	EmailAddresses []struct {
		EmailMetadata struct {
			Primary  bool `json:"primary"`
			Verified bool `json:"verified"`
			Source   struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			} `json:"source"`
		} `json:"metadata"`
		Value string `json:"value"`
	} `json:"emailAddresses"`
}


func init() {
	provider := &GoogleProvider{}
	einterfaces.RegisterOauthProvider(model.USER_AUTH_SERVICE_GOOGLE, provider)
}

func userFromGoogleUser(glu *GoogleUser) *model.User {
	user := &model.User{}
	username := glu.EmailAddresses[0].Value

	splitUserName := strings.Split(username, "@")
	if ((len(splitUserName) == 2) && (splitUserName[1] == "flipkart.com")) {
		user.Username = model.CleanUsername(splitUserName[0])
	} else {
		user.Username = model.CleanUsername(username)
	}

	splitName := strings.Split(glu.Names[0].DisplayName, " ")
	if len(splitName) == 2 {
		user.FirstName = splitName[0]
		user.LastName = splitName[1]
	} else if len(splitName) >= 2 {
		user.FirstName = splitName[0]
		user.LastName = strings.Join(splitName[1:], " ")
	} else {
		user.FirstName = glu.Names[0].GivenName
	}
	user.Email = glu.EmailAddresses[0].Value

	userId := glu.getAuthData()
	user.AuthData = &userId
	user.AuthService = model.USER_AUTH_SERVICE_GOOGLE

	return user
}

func googleUserFromJson(data io.Reader) *GoogleUser {
	decoder := json.NewDecoder(data)
	var glu GoogleUser
	err := decoder.Decode(&glu)
	if err == nil {
		return &glu
	} else {
		return nil
	}
}

func (glu *GoogleUser) ToJson() string {
	b, err := json.Marshal(glu)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func (glu *GoogleUser) IsValid() bool {
	if glu.EmailAddresses[0].EmailMetadata.Source.ID == "0" {
		return false
	}

	if glu.EmailAddresses[0].Value == "" {
		return false
	}

	return true
}


func (glu *GoogleUser) getAuthData() string {
	return glu.EmailAddresses[0].EmailMetadata.Source.ID
}

func (m *GoogleProvider) GetUserFromJson(data io.Reader) *model.User {
	glu := googleUserFromJson(data)
	if glu.IsValid() {
		return userFromGoogleUser(glu)
	}

	return &model.User{}
}
