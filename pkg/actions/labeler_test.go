package labeler

import (
	"errors"
	"testing"

	"github.com/ppapapetrou76/virtual-assistant/pkg/config"
	"github.com/ppapapetrou76/virtual-assistant/pkg/github"
	testutil "github.com/ppapapetrou76/virtual-assistant/pkg/util"
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

func TestLabeler_HandleEvent(t *testing.T) {
	type fields struct {
		repo github.Repo
	}
	type args struct {
		labels  []string
		payload []byte
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
				labels:  []string{"bug", "enhancement"},
				payload: []byte(webhookPayload),
			},
			fields: fields{
				repo: github.Repo{
					GHClient: github.MockGithubClient(200, ""),
					Owner:    "ppapapetrou76",
					Name:     "virtual-assistant",
				},
			},
		},
		{
			name: "should return error if action run fails",
			args: args{
				labels:  []string{"bug", "enhancement"},
				payload: []byte(webhookPayload),
			},
			fields: fields{
				repo: github.Repo{
					GHClient: github.MockGithubClient(401, `{
						  "message": "Bad credentials",
						  "documentation_url": "https://developer.github.com/v3"
						}`),
					Owner: "ppapapetrou76",
					Name:  "virtual-assistant",
				},
			},
			wantErr:       true,
			expectedError: errors.New("GET https://api.github.com/repos/ppapapetrou76/virtual-assistant/issues/2/labels: 401 Bad credentials []"),
		},
		{
			name: "should return error parsing webhook",
			args: args{
				labels:  []string{"bug", "enhancement"},
				payload: []byte("random payload"),
			},
			fields: fields{
				repo: github.Repo{
					GHClient: github.MockGithubClient(200, ""),
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
				Config: &config.Config{Labels: []string{"bug"}},
				Repo:   tt.fields.repo,
			}
			err := labeler.HandleEvent("pull_request", &tt.args.payload)
			testutil.AssertError(t, tt.wantErr, tt.expectedError, err)
		})
	}
}
