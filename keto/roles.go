package keto

import (
	"fmt"
	"net/http"
)

func (c *Client) GetORYAccessControlPolicyRole(flavor, id string) (*ORYAccessControlPolicyRoleJSON, bool, error) {

	var jsonClient *ORYAccessControlPolicyRoleJSON

	req, err := c.newRequest(http.MethodGet, makePath(flavor, resourceRoles, id), nil)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.do(req, &jsonClient)
	if err != nil {
		return nil, false, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return jsonClient, true, nil
	case http.StatusNotFound:
		return nil, false, nil
	default:
		return nil, false, fmt.Errorf("%s %s http request returned unexpected status code %s", req.Method, req.URL.String(), resp.Status)
	}
}

func (c *Client) ListORYAccessControlPolicyRole(flavor string) ([]*ORYAccessControlPolicyRoleJSON, error) {

	var jsonClientList []*ORYAccessControlPolicyRoleJSON

	req, err := c.newRequest(http.MethodGet, makePath(flavor, resourceRoles, ""), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req, &jsonClientList)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return jsonClientList, nil
	default:
		return nil, fmt.Errorf("%s %s http request returned unexpected status code %s", req.Method, req.URL.String(), resp.Status)
	}
}

func (c *Client) PutORYAccessControlPolicyRole(flavor string, o *ORYAccessControlPolicyRoleJSON) (*ORYAccessControlPolicyRoleJSON, error) {

	var jsonClient *ORYAccessControlPolicyRoleJSON

	req, err := c.newRequest(http.MethodPut, makePath(flavor, resourceRoles, ""), o)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req, &jsonClient)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s %s http request returned unexpected status code: %s", req.Method, req.URL, resp.Status)
	}

	return jsonClient, nil
}

func (c *Client) DeleteORYAccessControlPolicyRole(flavor, id string) error {

	req, err := c.newRequest(http.MethodDelete, makePath(flavor, resourceRoles, id), nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		fmt.Printf("client with id %s does not exist", id)
		return nil
	default:
		return fmt.Errorf("%s %s http request returned unexpected status code %s", req.Method, req.URL.String(), resp.Status)
	}
}
