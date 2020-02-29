package freebox

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const AuthorizationHeader = "X-Fbx-App-Auth"

type Client struct {
	discoveryUrl string
	client       *http.Client
	apiVersion   *APIVersion
	baseUrl      string

	appId        string
	appToken     string
	appVersion   string
	sessionToken string
}

type APIVersion struct {
	BoxModelName   string `json:"box_model_name"`
	APIBaseURL     string `json:"api_base_url"`
	HTTPSPort      int    `json:"https_port"`
	DeviceName     string `json:"device_name"`
	HTTPSAvailable bool   `json:"https_available"`
	BoxModel       string `json:"box_model"`
	APIDomain      string `json:"api_domain"`
	UID            string `json:"uid"`
	APIVersion     string `json:"api_version"`
	DeviceType     string `json:"device_type"`
}

type APIResponse struct {
	Success   bool    `json:"success"`
	Message   *string `json:"msg"`
	ErrorCode *string `json:"error_code"`
}

type Option func(client *Client) error

func NewClient(options ...Option) (*Client, error) {
	client := &Client{
		discoveryUrl: "http://mafreebox.freebox.fr",
		client:       http.DefaultClient,
	}

	for _, option := range options {
		err := option(client)
		if err != nil {
			return nil, err
		}
	}

	// Discover remote configuration
	if err := client.discover(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) discover() error {
	// Fetch api version data
	resp, err := c.client.Get(strings.TrimRight(c.discoveryUrl, "/") + "/api_version")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var output *APIVersion
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return err
	}

	parts := strings.Split(output.APIVersion, ".")
	if len(parts) == 0 {
		return errors.New("invalid version")
	}
	if parts[0] != "6" {
		return errors.New("unsupported version")
	}

	c.apiVersion = output
	// TODO: support external api (api_domain:https_port)
	c.baseUrl = fmt.Sprintf("http://mafreebox.freebox.fr%sv%s/", c.apiVersion.APIBaseURL, parts[0])

	return nil
}

func (c *Client) request(method, resource string, authenticated bool, body interface{}, output interface{}) error {
	var payload io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		payload = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, c.baseUrl+resource, payload)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if authenticated {
		req.Header.Set(AuthorizationHeader, c.sessionToken)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return err
	}

	// TODO: extract challenge
	// TODO: check success
	// TODO: check 403 return code for auth_required

	return nil
}

func (c *Client) Get(resource string, authenticated bool, output interface{}) error {
	return c.request("GET", resource, authenticated, nil, output)
}

func (c *Client) Post(resource string, authenticated bool, body interface{}, output interface{}) error {
	return c.request("POST", resource, authenticated, body, output)
}

func (c *Client) Put(resource string, authenticated bool, body interface{}, output interface{}) error {
	return c.request("PUT", resource, authenticated, body, output)
}

func (c *Client) APIVersion() *APIVersion {
	return c.apiVersion
}

func (c *Client) SetApp(appId, appToken, appVersion string) {
	c.appId = appId
	c.appToken = appToken
	c.appVersion = appVersion
}

func (c *Client) SetSessionToken(sessionToken string) {
	c.sessionToken = sessionToken
}

func WithDiscoveryURL(url string) Option {
	return func(client *Client) error {
		client.discoveryUrl = url
		return nil
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(client *Client) error {
		client.client = httpClient
		return nil
	}
}

func WithApp(appId, appToken, appVersion string) Option {
	return func(client *Client) error {
		client.appId = appId
		client.appToken = appToken
		client.appVersion = appVersion
		return nil
	}
}

func WithSessionToken(sessionToken string) Option {
	return func(client *Client) error {
		client.sessionToken = sessionToken
		return nil
	}
}
