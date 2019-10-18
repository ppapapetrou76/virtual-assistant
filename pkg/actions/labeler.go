package labeler

import (
	"log"

	gh "github.com/google/go-github/v27/github"
	"github.com/hashicorp/go-multierror"

	"github.com/ppapapetrou76/virtual-assistant/pkg/config"
	"github.com/ppapapetrou76/virtual-assistant/pkg/github"
	"github.com/ppapapetrou76/virtual-assistant/pkg/util/slices"
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

	var actions slices.StringSlice

	switch event := event.(type) {
	case *gh.PullRequestEvent:
		actions = l.Config.PullRequestsLabelerConfig.Actions
		if actions.IsEmpty() {
			actions.Add("opened")
		}
		if !actions.HasString(*event.Action) {
			log.Printf("Pull request event is `%s` - eligible actions are `%v`. Skipping issues labeler", *event.Action, actions)
			return nil
		}
		err = l.runOn(event.PullRequest)
	case *gh.IssuesEvent:
		actions = l.Config.IssuesLabelerConfig.Actions
		if actions.IsEmpty() {
			actions.Add("opened")
		}
		if !actions.HasString(*event.Action) {
			log.Printf("Issues event is `%s` - eligible actions are `%v`. Skipping issues labeler", *event.Action, actions)
			return nil
		}
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

	desiredLabels := append(l.PullRequestsLabelerConfig.Labels, currLabels...)
	log.Printf("Desired labels: %s", desiredLabels)
	return pullRequest.ReplaceLabels(desiredLabels)
}

func (l *Labeler) runOnIssue(i *gh.Issue) error {
	issue := github.NewIssue(l.Repo, *i.Number)
	currLabels, err := issue.CurrentLabels()

	if err != nil {
		return err
	}

	desiredLabels := append(l.IssuesLabelerConfig.Labels, currLabels...)
	log.Printf("Desired labels: %s", desiredLabels)

	merr := new(multierror.Error)
	merr = multierror.Append(merr, issue.ReplaceLabels(desiredLabels))
	merr = multierror.Append(merr, issue.AtLeastOne(l.IssuesLabelerConfig.PossibleLabels, l.IssuesLabelerConfig.Default))

	return merr.ErrorOrNil()
}

// New creates a new labeler object
func New(c *config.Config, repo github.Repo) *Labeler {
	return &Labeler{
		Config: c,
		Repo:   repo,
	}
}
