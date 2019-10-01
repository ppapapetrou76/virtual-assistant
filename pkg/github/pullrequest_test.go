package github

import (
	"errors"
	"reflect"
	"testing"

	testutil "github.com/ppapapetrou76/virtual-assistant/pkg/util"
)

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

func TestPullRequest_CurrentLabels(t *testing.T) {
	type fields struct {
		ghClient ClientWrapper
	}
	tests := []struct {
		name           string
		fields         fields
		wantErr        bool
		expectedError  error
		expectedLabels []string
	}{
		{
			name: "should return the issue labels",
			fields: fields{
				ghClient: MockGithubClient(200, listIssueLabelsResponse),
			},
			expectedLabels: []string{"bug", "enhancement"},
		},
		{
			name: "should error if labels cannot be loaded",
			fields: fields{
				ghClient: MockGithubClient(200, "ok"),
			},
			expectedError: errors.New("invalid character 'o' looking for beginning of value"),
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := Repo{
				GHClient: tt.fields.ghClient,
				Owner:    "ppapapetrou76",
				Name:     "virtual-assistant",
			}

			pr := PullRequest{
				Repo:   repo,
				Number: 0,
			}
			actualLabels, err := pr.CurrentLabels()
			testutil.AssertError(t, tt.wantErr, tt.expectedError, err)

			if !tt.wantErr && !reflect.DeepEqual(actualLabels, tt.expectedLabels) {
				t.Errorf("Expect: \n%+v Got: \n%+v", tt.expectedLabels, actualLabels)
			}
		})
	}
}

func TestNewPullRequest(t *testing.T) {
	repo := Repo{
		Owner: "ppapapetrou76",
		Name:  "virtual-assistant",
	}
	type args struct {
		repo   Repo
		number int
	}
	tests := []struct {
		name     string
		args     args
		expected PullRequest
	}{
		{
			name: "should return a new repo",
			args: args{
				repo:   repo,
				number: 123,
			},
			expected: PullRequest{
				Repo:   repo,
				Number: 123,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewPullRequest(tt.args.repo, tt.args.number)

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("Expect: \n%+v Got: \n%+v", tt.expected, actual)
			}
		})
	}
}

func TestPullRequest_ReplaceLabels(t *testing.T) {
	type fields struct {
		ghClient ClientWrapper
	}
	type args struct {
		labels []string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		expectedError  error
		expectedLabels []string
	}{
		{
			name: "should replace the issue labels",
			args: args{labels: []string{"bug", "enhancement"}},
			fields: fields{
				ghClient: MockGithubClient(200, listIssueLabelsResponse),
			},
		},
		{
			name: "should error if labels cannot be replaced",
			fields: fields{
				ghClient: MockGithubClient(401, `{
				  "message": "Bad credentials",
  				  "documentation_url": "https://developer.github.com/v3"
				}`),
			},
			expectedError: errors.New("PUT https://api.github.com/repos/ppapapetrou76/virtual-assistant/issues/0/labels: 401 Bad credentials []"),
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := Repo{
				GHClient: tt.fields.ghClient,
				Owner:    "ppapapetrou76",
				Name:     "virtual-assistant",
			}

			pr := PullRequest{
				Repo:   repo,
				Number: 0,
			}
			err := pr.ReplaceLabels(tt.args.labels)
			testutil.AssertError(t, tt.wantErr, tt.expectedError, err)
		})
	}
}
