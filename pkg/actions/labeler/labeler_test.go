package labeler

import (
	"errors"
	"testing"

	"github.com/ppapapetrou76/virtual-assistant/pkg/config"
	"github.com/ppapapetrou76/virtual-assistant/pkg/github"
	"github.com/ppapapetrou76/virtual-assistant/pkg/testutil"
)

const webhookPayload = `{
  "action": "opened",
  "number": 2,
  "pull_request": {
    "id": 279147437,
    "node_id": "MDExOlB1bGxSZXF1ZXN0Mjc5MTQ3NDM3",
    "number": 2,
    "state": "open",
    "locked": false,
    "title": "Update the README with new information.",
    "body": "This is a pretty simple change that we need to pull into master."
   }
}`

const webhookIssuePayload = `
{
  "action": "opened",
  "issue": {
    "id": 444500041,
    "node_id": "MDU6SXNzdWU0NDQ1MDAwNDE=",
    "number": 1,
    "title": "Some random issue"
  }
}`

func TestLabeler_HandleEvent(t *testing.T) {
	type fields struct {
		repo github.Repo
	}
	type args struct {
		labels    []string
		payload   []byte
		eventName string
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
			name: "should handle a pr event",
			args: args{
				labels:    []string{"bug", "enhancement"},
				payload:   []byte(webhookPayload),
				eventName: "pull_request",
			},
			fields: fields{
				repo: github.Repo{
					GHClient: github.MockGithubClient([]github.MockResponse{
						github.MockGenericSuccessResponse(),
						github.MockGenericSuccessResponse(),
						github.MockGenericSuccessResponse(),
					}),
					Owner: "ppapapetrou76",
					Name:  "virtual-assistant",
				},
			},
		},
		{
			name: "should handle an issue event",
			args: args{
				labels:    []string{"bug", "enhancement"},
				payload:   []byte(webhookIssuePayload),
				eventName: "issues",
			},
			fields: fields{
				repo: github.Repo{
					GHClient: github.MockGithubClient([]github.MockResponse{
						github.MockGenericSuccessResponse(),
						github.MockGenericSuccessResponse(),
						github.MockGenericSuccessResponse(),
					}),
					Owner: "ppapapetrou76",
					Name:  "virtual-assistant",
				},
			},
		},
		{
			name: "should return error if event is pull request and fetching current label fails",
			args: args{
				labels:    []string{"bug", "enhancement"},
				payload:   []byte(webhookPayload),
				eventName: "pull_request",
			},
			fields: fields{
				repo: github.Repo{
					GHClient: github.MockGithubClient([]github.MockResponse{
						github.UnAuthorizedMockResponse(),
					}),
					Owner: "ppapapetrou76",
					Name:  "virtual-assistant",
				},
			},
			wantErr:       true,
			expectedError: errors.New("GET https://api.github.com/repos/ppapapetrou76/virtual-assistant/issues/2/labels: 401 Bad credentials []"),
		},
		{
			name: "should return error if event is issue and fetching current label fails",
			args: args{
				labels:    []string{"bug", "enhancement"},
				payload:   []byte(webhookIssuePayload),
				eventName: "issues",
			},
			fields: fields{
				repo: github.Repo{
					GHClient: github.MockGithubClient([]github.MockResponse{
						github.UnAuthorizedMockResponse(),
					}),
					Owner: "ppapapetrou76",
					Name:  "virtual-assistant",
				},
			},
			wantErr:       true,
			expectedError: errors.New("GET https://api.github.com/repos/ppapapetrou76/virtual-assistant/issues/1/labels: 401 Bad credentials []"),
		},
		{
			name: "should return error parsing webhook",
			args: args{
				labels:    []string{"bug", "enhancement"},
				payload:   []byte("random payload"),
				eventName: "pull_request",
			},
			fields: fields{
				repo: github.Repo{
					GHClient: github.MockGithubClient([]github.MockResponse{github.MockGenericSuccessResponse()}),
					Owner:    "ppapapetrou76",
					Name:     "virtual-assistant",
				},
			},
			wantErr:       true,
			expectedError: errors.New("invalid character 'r' looking for beginning of value"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labeler := Labeler{
				LabelerConfig: &config.LabelerConfig{
					PullRequestsLabelerConfig: config.PullRequestsLabelerConfig{
						Labels: []string{"bug"},
					},
					IssuesLabelerConfig: config.IssuesLabelerConfig{
						Labels: []string{"feature"},
					},
				},
				Repo: tt.fields.repo,
			}
			err := labeler.HandleEvent(tt.args.eventName, &tt.args.payload)
			testutil.AssertError(t, tt.wantErr, tt.expectedError, err)
		})
	}
}
