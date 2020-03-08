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

func TestAssigner_HandleEvent(t *testing.T) {
	type args struct {
		labels    []string
		payload   []byte
		eventName string
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		expectedError error
		assigner      Assigner
	}{
		{
			name: "should handle an issue event",
			args: args{
				labels:    []string{"bug", "enhancement"},
				payload:   []byte(webhookIssuePayload),
				eventName: "issues",
			},
			assigner: issueAssigner(github.Repo{
				GHClient: github.MockGithubClient([]github.MockResponse{
					github.MockGetIssueResponse(),
					github.MockListRepositoryProjectsResponse(),
					github.MockListProjectColumnsResponse(),
					github.MockGenericSuccessResponse(),
				}),
			}),
		},
		{
			name: "should return error if event is issue and fetching current label fails",
			args: args{
				labels:    []string{"bug", "enhancement"},
				payload:   []byte(webhookIssuePayload),
				eventName: "issues",
			},
			assigner: issueAssigner(github.Repo{
				GHClient: github.MockGithubClient([]github.MockResponse{
					github.UnAuthorizedMockResponse(),
				}),
				Owner: "ppapapetrou76",
				Name:  "virtual-assistant",
			}),
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
			assigner: issueAssigner(github.Repo{
				GHClient: github.MockGithubClient([]github.MockResponse{github.MockGenericSuccessResponse()}),
				Owner:    "ppapapetrou76",
				Name:     "virtual-assistant",
			}),
			wantErr:       true,
			expectedError: errors.New("invalid character 'r' looking for beginning of value"),
		},
		{
			name: "should handle a PR event",
			args: args{
				payload:   []byte(webhookPayload),
				eventName: "pull_request",
			},
			assigner: prAssigner(github.Repo{
				GHClient: github.MockGithubClient([]github.MockResponse{
					github.MockGetIssueResponse(),
					github.MockGetIssueResponse(),
					github.MockGenericSuccessResponse(),
				}),
			}),
		},
		{
			name: "should return error if event is PR and fetching github issue fails",
			args: args{
				payload:   []byte(webhookPayload),
				eventName: "pull_request",
			},
			assigner: prAssigner(github.Repo{
				GHClient: github.MockGithubClient([]github.MockResponse{
					github.UnAuthorizedMockResponse(),
				}),
			}),
			wantErr:       true,
			expectedError: errors.New("cannot get issue with number 2. error message : GET https://api.github.com/repos///issues/2: 401 Bad credentials []"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.assigner.HandleEvent(tt.args.eventName, &tt.args.payload)
			testutil.AssertError(t, tt.wantErr, tt.expectedError, err)
		})
	}
}

func issueAssigner(repo github.Repo) Assigner {
	return Assigner{
		AssignerConfig: &config.AssignerConfig{
			IssuesAssignerConfig: config.IssuesAssignerConfig{
				IssuesAssignerProjectConfig: config.IssuesAssignerProjectConfig{
					ProjectURL: "https://github.com/ppapapetrou76/virtual-assistant/projects/1",
					Column:     "To Do",
				},
				Actions: []string{"opened"},
			},
		},
		Repo: repo,
	}
}

func prAssigner(repo github.Repo) Assigner {
	return Assigner{
		AssignerConfig: &config.AssignerConfig{
			PullRequestsAssignerConfig: config.PullRequestsAssignerConfig{
				Assignee: config.PullRequestsAutoAssigneeConfig{
					Auto: true,
				},
				Actions: []string{"opened"},
			},
		},
		Repo: repo,
	}
}
