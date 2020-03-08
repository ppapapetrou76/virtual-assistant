package github

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

const getIssueResponse = `{
  "id": 1,
  "node_id": "MDU6SXNzdWUx",
  "number": 1347,
  "state": "open",
  "title": "Found a bug",
  "body": "I'm having a problem with this.",
  "user": {
    "login": "ppapapetrou76",
    "id": 1
  }
}
`
const listProjectColumnsResponse = `[
  {
    "url": "https://api.github.com/projects/columns/367",
    "project_url": "https://api.github.com/projects/120",
    "cards_url": "https://api.github.com/projects/columns/367/cards",
    "id": 367,
    "node_id": "MDEzOlByb2plY3RDb2x1bW4zNjc=",
    "name": "To Do",
    "created_at": "2016-09-05T14:18:44Z",
    "updated_at": "2016-09-05T14:22:28Z"
  }
]`

const listIssueLabelsResponse = `[
  {
    "id": 208045946,
    "node_id": "MDU6TGFiZWwyMDgwNDU5NDY=",
    "name": "bug",
    "description": "Something isn't working",
    "color": "f29513",
    "default": true
  },
  {
    "id": 208045947,
    "node_id": "MDU6TGFiZWwyMDgwNDU5NDc=",
    "name": "enhancement",
    "description": "New feature or request",
    "color": "a2eeef",
    "default": false
  }
]`

const listRepositoryProjectsResponse = `[
  {
    "owner_url": "https://api.github.com/repos/api-playground/projects-test",
    "html_url": "https://github.com/ppapapetrou76/virtual-assistant/projects/1",
    "columns_url": "https://api.github.com/projects/1002604/columns",
    "id": 1002604,
    "node_id": "MDc6UHJvamVjdDEwMDI2MDQ=",
    "name": "Projects Documentation",
    "body": "Developer documentation project for the developer site.",
    "number": 1,
    "state": "open"
  }
]`

const listOrganizationProjectsResponse = `[
  {
    "owner_url": "https://api.github.com/repos/api-playground/projects-test",
    "html_url": "https://github.com/orgs/myorg/projects/1",
    "columns_url": "https://api.github.com/projects/9999/columns",
    "id": 9999,
    "node_id": "MDc6UHJvamVjdDEwMDI2MDQ=",
    "name": "Projects Documentation",
    "body": "Developer documentation project for the developer site.",
    "number": 1,
    "state": "open"
  }
]`
const listEmptyProjectsResponse = `[]`

// MockResponse mocks an http response
type MockResponse struct {
	StatusCode int
	Response   string
}

// MockRoundTripper mocks a RoundTripper
type MockRoundTripper struct {
	Responses         []MockResponse
	nextResponseIndex int
}

// RoundTrip implements the RoundTripper interface
func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r := m.nextResponse()
	return &http.Response{
		StatusCode: r.StatusCode,
		Body:       ioutil.NopCloser(bytes.NewBufferString(r.Response)),
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
func MockGithubClient(responses []MockResponse) ClientWrapper {
	return Client(NewTestClient(&MockRoundTripper{
		Responses: responses,
	}))
}

func (m *MockRoundTripper) nextResponse() MockResponse {
	if m.nextResponseIndex >= len(m.Responses) {
		panic("no more responses mocked. please add more and re-run the test")
	}

	r := m.Responses[m.nextResponseIndex]
	m.nextResponseIndex++
	return r
}

// UnAuthorizedMockResponse returns a mock response with a 401 error code and message
func UnAuthorizedMockResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusUnauthorized,
		Response: `{
					  "message": "Bad credentials",
  		 			  "documentation_url": "https://developer.github.com/v3"
				   }`,
	}
}

// MockListIssueLabelsResponse returns a mock response for the list issue labels call
func MockListIssueLabelsResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusOK,
		Response:   listIssueLabelsResponse,
	}
}

// MockGetIssueResponse returns a mock response for the get issue call
func MockGetIssueResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusOK,
		Response:   getIssueResponse,
	}
}

// MockListProjectColumnsResponse returns a mock response for the list project columns call
func MockListProjectColumnsResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusOK,
		Response:   listProjectColumnsResponse,
	}
}

// MockListRepositoryProjectsResponse returns a mock response for the list repository projects call
func MockListRepositoryProjectsResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusOK,
		Response:   listRepositoryProjectsResponse,
	}
}

// MockListOrganizationProjectsResponse returns a mock response for the list organization projects call
func MockListOrganizationProjectsResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusOK,
		Response:   listOrganizationProjectsResponse,
	}
}

// MockListEmptyProjectsResponse returns a mock response with an empty projects list
func MockListEmptyProjectsResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusOK,
		Response:   listEmptyProjectsResponse,
	}
}

// MockGenericSuccessResponse returns a generic success mock response
func MockGenericSuccessResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusOK,
	}
}
