package github

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// MockRoundTripper mocks a RoundTripper
type MockRoundTripper struct {
	StatusCode int
	Response   string
}

// RoundTrip implements the RoundTripper interface
func (m MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: m.StatusCode,
		Body:       ioutil.NopCloser(bytes.NewBufferString(m.Response)),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn http.RoundTripper) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

// MockGithubClient returns a mocked Github client for testing purposes
func MockGithubClient(statusCode int, response string) ClientWrapper {
	return Client(NewTestClient(MockRoundTripper{
		StatusCode: statusCode,
		Response:   response,
	}))
}
