#!/bin/bash

mkdir $1
touch $1/client.go
echo 'package '$1 > $1/client.go
echo '
import (
    "context"
	"net/http"
)

type RequestEditorFn func(ctx context.Context, req *http.Request) error

type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientOption func(*Client) error

func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

type Client struct {
	Server         string
	Client         HttpRequestDoer
	RequestEditors []RequestEditorFn
}

func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// ensure the server URL always has a trailing slash
	client := &Client{
		Server: server,
	}

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return client, nil

}' >> $1/client.go