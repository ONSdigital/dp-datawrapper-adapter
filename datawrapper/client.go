package datawrapper

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/v2/log"
)

// Client is a Datawrapper client which can be used to make requests to the server
type Client struct {
	APIURL     string
	APIToken   string
	HTTPClient *http.Client
}

// DefaultTransport is the default configuration of Transport
var DefaultTransport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout: 5 * time.Second,
	}).DialContext,
	TLSHandshakeTimeout: 5 * time.Second,
	MaxIdleConns:        10,
	IdleConnTimeout:     30 * time.Second,
}

// NewClient creates a new instance of Client with a given url and token
func NewClient(APIURL string, APIToken string) *Client {
	return &Client{
		APIURL:   APIURL,
		APIToken: APIToken,
		HTTPClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: DefaultTransport,
		},
	}
}

// Checker calls an app health endpoint and returns a check object to the caller
func (c *Client) Checker(ctx context.Context, state *health.CheckState) error {

	code, err := c.get(ctx, "/v3/me", c.APIToken)
	if err != nil {
		log.Error(ctx, "failed to request datawrapper health", err)
	}

	log.Info(ctx, fmt.Sprintf("datawrapper health response code: %v", code))

	switch {
	case code == 0: // When there is a problem with the client return error in message
		return state.Update(health.StatusCritical, err.Error(), 0)
	case code == 200:
		return state.Update(health.StatusOK, "datawrapper is ok", code)
	case code >= 400 && code < 500:
		return state.Update(health.StatusWarning, "datawrapper is degraded, but at least partially functioning", code)
	default:
		return state.Update(health.StatusCritical, "datawrapper functionality is unavailable or non-functioning", code)
	}
}

func (c *Client) get(ctx context.Context, path string, token string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.APIURL+path, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer closeResponseBody(ctx, resp)

	return resp.StatusCode, nil
}

// closeResponseBody closes the response body and logs an error if unsuccessful
func closeResponseBody(ctx context.Context, resp *http.Response) {
	if resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			log.Error(ctx, "error closing http response body", err)
		}
	}
}
