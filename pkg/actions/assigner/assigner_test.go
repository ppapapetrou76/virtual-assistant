package assigner

import (
	"errors"
	"testing"

	"github.com/ppapapetrou76/virtual-assistant/pkg/config"
	"github.com/ppapapetrou76/virtual-assistant/pkg/github"
	"github.com/ppapapetrou76/virtual-assistant/pkg/testutil"
)

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

func TestAssigner_HandleEvent(t *testing.T) {
	type fields struct {
		repo github.Repo
	}
	type args struct {
		labels    []string
		payload   []byte
		eventName string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		expectedError error
	}{
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
						github.MockGetIssueResponse(),
						github.MockListProjectColumnsResponse(),
						github.MockGenericSuccessResponse(),
					}),
				},
			},
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
			expectedError: errors.New("cannot get issue with number 1. error message : GET https://api.github.com/repos/ppapapetrou76/virtual-assistant/issues/1: 401 Bad credentials []"),
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
			assigner := Assigner{
				AssignerConfig: &config.AssignerConfig{
					IssuesAssignerConfig: config.IssuesAssignerConfig{
						IssuesAssignerProjectConfig: config.IssuesAssignerProjectConfig{
							ProjectID: 1,
							Column:    "To Do",
						},
						Actions: []string{"opened"},
					},
				},
				Repo: tt.fields.repo,
			}
			err := assigner.HandleEvent(tt.args.eventName, &tt.args.payload)
			testutil.AssertError(t, tt.wantErr, tt.expectedError, err)
		})
	}
}
