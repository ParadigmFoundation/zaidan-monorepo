package tests

import (
	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	rpc *rpc.Client
}

func NewClient(addr string) (*Client, error) {
	c, err := rpc.Dial(addr)
	if err != nil {
		return nil, err
	}

	return &Client{rpc: c}, nil
}

func (c *Client) Close() { c.rpc.Close() }

func (c *Client) Call(resp interface{}, meth string, args ...interface{}) error {
	return c.rpc.Call(resp, meth, args...)
}
