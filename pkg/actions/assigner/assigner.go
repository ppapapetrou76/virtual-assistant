package assigner

import (
	gh "github.com/google/go-github/v27/github"

	"github.com/ppapapetrou76/virtual-assistant/pkg/actions"
	"github.com/ppapapetrou76/virtual-assistant/pkg/config"
	"github.com/ppapapetrou76/virtual-assistant/pkg/github"
)

// Assigner is the struct to handle auto-assigning issues and pull-requests
type Assigner struct {
	*config.AssignerConfig
	github.Repo
}

// HandleEvent takes a GitHub Event and its raw payload (see link below)
// to trigger an update to the issue / PR's labels.
//
// https://developer.github.com/v3/activity/events/types/
func (l *Assigner) HandleEvent(eventName string, payload *[]byte) error {
	event, err := gh.ParseWebHook(eventName, *payload)
	if err != nil {
		return err
	}
	switch event := event.(type) {
	case *gh.PullRequestEvent:
		if actions.ShouldRunOnPullRequest(event, l.IssuesAssignerConfig.Actions) {
			err = l.runOnPR(event.PullRequest)
		}
	case *gh.IssuesEvent:
		if actions.ShouldRunOnIssue(event, l.IssuesAssignerConfig.Actions) {
			err = l.runOnIssue(event.Issue)
		}
	}
	return err
}

func (l *Assigner) runOnPR(i *gh.PullRequest) error {
	issue := github.NewIssue(l.Repo, *i.Number)
	if !l.Assignee.Auto {
		return nil
	}
	return issue.AddAssignee()
}

func (l *Assigner) runOnIssue(i *gh.Issue) error {
	issue := github.NewIssue(l.Repo, *i.Number)
	return issue.AddToProject(l.AssignerConfig.ProjectURL, l.AssignerConfig.Column)
}

// New creates a new labeler object
func New(c *config.Config, repo github.Repo) *Assigner {
	return &Assigner{
		AssignerConfig: &c.AssignerConfig,
		Repo:           repo,
	}
}
