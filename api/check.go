package api

import "github.com/jrasell/nomad-sherlock/sdk"

type Check struct {
	client *Client
}

func (c *Client) Check() *Check {
	return &Check{client: c}
}

func (c *Check) List(q *QueryOptions) ([]*sdk.CheckInfo, error) {
	var resp []*sdk.CheckInfo
	err := c.client.query("/v1/checks", &resp, q)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Check) Info(id string, q *QueryOptions) (*sdk.CheckInfo, error) {
	var resp *sdk.CheckInfo
	err := c.client.query("/v1/check/"+id, &resp, q)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
