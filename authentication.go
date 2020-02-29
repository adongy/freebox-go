package freebox

import (
	"errors"
	"fmt"
	"time"
)

type TokenRequestPayload struct {
	AppID      string `json:"app_id"`
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
	DeviceName string `json:"device_name"`
}

type AuthorizationResponse struct {
	APIResponse
	Result struct {
		AppToken string `json:"app_token"`
		TrackID  int    `json:"track_id"`
	} `json:"result"`
}

type AuthorizationProgressResponse struct {
	APIResponse
	Result struct {
		Status    string `json:"status"`
		Challenge string `json:"challenge"`
	} `json:"result"`
}

type OpenSessionResponse struct {
	APIResponse
	Result struct {
		SessionToken string          `json:"session_token"`
		Permissions  map[string]bool `json:"permissions"`
		Challenge    string          `json:"challenge"`
	} `json:"result"`
}

type LoginResponse struct {
	APIResponse
	Result struct {
		LoggedIn  bool   `json:"logged_in"`
		Challenge string `json:"challenge"`
	} `json:"result"`
}

func (c *Client) RequestAuthorization(payload *TokenRequestPayload) (*AuthorizationResponse, error) {
	var output *AuthorizationResponse
	if err := c.Post("login/authorize/", false, payload, &output); err != nil {
		return nil, err
	}

	return output, nil
}

func (c *Client) TrackAuthorizationProgress(trackId int) (*AuthorizationProgressResponse, error) {
	var output *AuthorizationProgressResponse
	if err := c.Get(fmt.Sprintf("login/authorize/%d", trackId), false, &output); err != nil {
		return nil, err
	}

	return output, nil
}

func (c *Client) OpenSession(challenge string) (*OpenSessionResponse, error) {
	if c.appId == "" || c.appToken == "" {
		return nil, errors.New("no app configured")
	}

	password, err := generatePassword(c.appToken, challenge)
	if err != nil {
		return nil, err
	}

	var output *OpenSessionResponse
	if err := c.Post("login/session/", false, map[string]interface{}{
		"app_id":   c.appId,
		"password": password,
	}, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// Wrapper to request a new app, and wait for status changes
// and configure it on our client
func (c *Client) Authorize(payload *TokenRequestPayload) (string, error) {
	resp, err := c.RequestAuthorization(payload)
	if err != nil {
		return "", err
	}

	var (
		status   string
		trackId  = resp.Result.TrackID
		appToken = resp.Result.AppToken
	)
	for {
		resp, err := c.TrackAuthorizationProgress(trackId)
		if err != nil {
			return "", err
		}

		status = resp.Result.Status
		if status != "pending" {
			break
		}

		time.Sleep(2 * time.Second)
	}

	if status != "granted" {
		return "", fmt.Errorf("invalid authorization status: %s", status)
	}

	c.SetApp(payload.AppID, appToken, payload.AppVersion)

	return appToken, nil
}

// Wrapper to login using the configured app on our client
func (c *Client) Login() error {
	if c.appId == "" || c.appToken == "" {
		return errors.New("no app configured")
	}

	var output *LoginResponse
	if err := c.Get("login/", true, &output); err != nil {
		return err
	}

	if output.Result.LoggedIn {
		return nil
	}

	resp, err := c.OpenSession(output.Result.Challenge)
	if err != nil {
		return err
	}

	c.SetSessionToken(resp.Result.SessionToken)
	return nil
}
