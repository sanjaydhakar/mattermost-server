package oauthgoogle


import (
	"encoding/json"
	"io"
	"strconv"
	"strings"

	"github.com/mattermost/mattermost-server/v5/einterfaces"
	"github.com/mattermost/mattermost-server/v5/model"
)

type GoogleProvider struct {
}

type GoogleUser struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

func init() {
	provider := &GoogleProvider{}
	einterfaces.RegisterOauthProvider(model.USER_AUTH_SERVICE_GOOGLE, provider)
}

func userFromGoogleUser(glu *GoogleUser) *model.User {
	user := &model.User{}
	username := glu.Username
	if username == "" {
		username = glu.Login
	}
	user.Username = model.CleanUsername(username)
	splitName := strings.Split(glu.Name, " ")
	if len(splitName) == 2 {
		user.FirstName = splitName[0]
		user.LastName = splitName[1]
	} else if len(splitName) >= 2 {
		user.FirstName = splitName[0]
		user.LastName = strings.Join(splitName[1:], " ")
	} else {
		user.FirstName = glu.Name
	}
	user.Email = glu.Email
	user.Email = strings.ToLower(user.Email)
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
	if glu.Id == 0 {
		return false
	}

	if len(glu.Email) == 0 {
		return false
	}

	return true
}


func (glu *GoogleUser) getAuthData() string {
	return strconv.FormatInt(glu.Id, 10)
}

func (m *GoogleProvider) GetUserFromJson(data io.Reader) *model.User {
	glu := googleUserFromJson(data)
	if glu.IsValid() {
		return userFromGoogleUser(glu)
	}

	return &model.User{}
}
