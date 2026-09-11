package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// apiClient makes authenticated requests to the crevisto platform.
type apiClient struct {
	base  string
	token string
	http  *http.Client
}

func newAPIClient() *apiClient {
	c := loadConfig()
	return &apiClient{
		base:  c.APIBase,
		token: c.APIToken,
		http:  &http.Client{Timeout: 60 * time.Second},
	}
}

// do performs an HTTP request with the bearer token.
func (c *apiClient) do(method, path string, body interface{}) (*http.Response, []byte, error) {
	var reqBody io.Reader
	if body != nil {
		j, _ := json.Marshal(body)
		reqBody = bytes.NewReader(j)
	}

	req, err := http.NewRequest(method, c.base+path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	return resp, data, nil
}

// get performs a GET request.
func (c *apiClient) get(path string) (*http.Response, []byte, error) {
	return c.do("GET", path, nil)
}

// post performs a POST request.
func (c *apiClient) post(path string, body interface{}) (*http.Response, []byte, error) {
	return c.do("POST", path, body)
}

// del performs a DELETE request.
func (c *apiClient) del(path string) (*http.Response, []byte, error) {
	return c.do("DELETE", path, nil)
}

// requireToken fails if no token is set.
func requireToken() string {
	t := getToken()
	if t == "" {
		fail(ExitResource, "not_authenticated", "no API token set",
			"crevisto auth <token>  — set your bearer token (invent any 16+ char string)")
	}
	return t
}
