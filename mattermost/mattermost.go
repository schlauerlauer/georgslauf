package mattermost

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"net/http"
)

type Client struct {
	url string
}

func NewClient(baseUrl string) *Client {
	return &Client{
		url: baseUrl,
	}
}

type getUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Firstname     string `json:"first_name"`
	Lastname      string `json:"last_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	PictureUpdate int64  `json:"last_picture_update"`
}

var (
	ErrInvalidStatusCode = errors.New("Response code is not ok")
)

func (client *Client) GetUser(accessToken string) (*getUser, error) {
	req, err := http.NewRequest(http.MethodGet, client.url+"/api/v4/users/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, ErrInvalidStatusCode
	}

	var resp getUser
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (client *Client) GetUserPicture(accessToken string, userId string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v4/users/%s/image", client.url, userId), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, ErrInvalidStatusCode
	}

	return io.ReadAll(res.Body)
}
