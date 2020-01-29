package keto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	KetoURL    url.URL
	HTTPClient *http.Client
}

const FlavorExact = "exact"
const FlavorGlob = "glob"
const FlavorRegex = "regex"

const resourcePolicies = "policies"
const resourceRoles = "roles"

func (c *Client) newRequest(method, relativePath string, body interface{}) (*http.Request, error) {

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	u := c.KetoURL
	u.Path = path.Join(u.Path, relativePath)

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	return req, nil

}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if v != nil && resp.StatusCode < 300 {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}

func makePath(flavor, resource, path string) string {
	if path != "" {
		if path[0] != '/' {
			path = "/" + path
		}
	}

	return fmt.Sprintf("/acp/ory/%s/%s%s", flavor, resource, path)
}
