package labeler

import (
	"log"

	gh "github.com/google/go-github/v27/github"
	"github.com/hashicorp/go-multierror"

	"github.com/ppapapetrou76/virtual-assistant/pkg/config"
	"github.com/ppapapetrou76/virtual-assistant/pkg/github"
)

// Labeler is the struct to handle auto-labeling of issues, PRs etc.
type Labeler struct {
	*config.Config
	github.Repo
}

// HandleEvent takes a GitHub Event and its raw payload (see link below)
// to trigger an update to the issue / PR's labels.
//
// https://developer.github.com/v3/activity/events/types/
func (l *Labeler) HandleEvent(eventName string, payload *[]byte) error {
	event, err := gh.ParseWebHook(eventName, *payload)
	if err != nil {
		return err
	}
	switch event := event.(type) {
	case *gh.PullRequestEvent:
		err = l.runOn(event.PullRequest)
	case *gh.IssuesEvent:
		err = l.runOnIssue(event.Issue)
	}
	return err
}

func (l *Labeler) runOn(pr *gh.PullRequest) error {
	pullRequest := github.NewIssue(l.Repo, *pr.Number)
	currLabels, err := pullRequest.CurrentLabels()

	if err != nil {
		return err
	}

	desiredLabels := append(l.PullRequestsConfig.Labels, currLabels...)
	log.Printf("Desired labels: %s", desiredLabels)
	return pullRequest.ReplaceLabels(desiredLabels)
}

func (l *Labeler) runOnIssue(i *gh.Issue) error {
	issue := github.NewIssue(l.Repo, *i.Number)
	currLabels, err := issue.CurrentLabels()

	if err != nil {
		return err
	}

	desiredLabels := append(l.IssuesConfig.Labels, currLabels...)
	log.Printf("Desired labels: %s", desiredLabels)

	merr := new(multierror.Error)
	merr = multierror.Append(merr, issue.ReplaceLabels(desiredLabels))
	merr = multierror.Append(merr, issue.AtLeastOne(l.IssuesConfig.PossibleLabels, l.IssuesConfig.Default))

	return merr.ErrorOrNil()
}

// New creates a new labeler object
func New(c *config.Config, repo github.Repo) *Labeler {
	return &Labeler{
		Config: c,
		Repo:   repo,
	}
}
