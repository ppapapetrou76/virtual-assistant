package actions

import (
	"reflect"
	"testing"

	"github.com/google/go-github/v27/github"

	"github.com/ppapapetrou76/virtual-assistant/pkg/util/slices"
)

func TestShouldRunOnIssue(t *testing.T) {
	opened := "opened"
	type args struct {
		event             *github.IssuesEvent
		configuredActions slices.StringSlice
	}
	tests := []struct {
		name     string
		expected bool
		args     args
	}{
		{
			name:     "should return true",
			expected: true,
			args: args{
				event: &github.IssuesEvent{
					Action: &opened,
				},
				configuredActions: []string{"opened", "closed"},
			},
		},
		{
			name:     "should use the default values and return true",
			expected: true,
			args: args{
				event: &github.IssuesEvent{
					Action: &opened,
				},
			},
		},
		{
			name: "should return false",
			args: args{
				event: &github.IssuesEvent{
					Action: &opened,
				},
				configuredActions: []string{"closed"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ShouldRunOnIssue(tt.args.event, tt.args.configuredActions)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("Expect: \n%+v Got: \n%+v", tt.expected, actual)
			}
		})
	}
}

func TestShouldRunOnPullRequest(t *testing.T) {
	opened := "opened"
	type args struct {
		event             *github.PullRequestEvent
		configuredActions slices.StringSlice
	}
	tests := []struct {
		name     string
		expected bool
		args     args
	}{
		{
			name:     "should return true",
			expected: true,
			args: args{
				event: &github.PullRequestEvent{
					Action: &opened,
				},
				configuredActions: []string{"opened", "closed"},
			},
		},
		{
			name:     "should use the default values and return true",
			expected: true,
			args: args{
				event: &github.PullRequestEvent{
					Action: &opened,
				},
			},
		},
		{
			name: "should return false",
			args: args{
				event: &github.PullRequestEvent{
					Action: &opened,
				},
				configuredActions: []string{"closed"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ShouldRunOnPullRequest(tt.args.event, tt.args.configuredActions)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("Expect: \n%+v Got: \n%+v", tt.expected, actual)
			}
		})
	}
}
