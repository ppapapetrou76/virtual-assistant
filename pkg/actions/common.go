package actions

import (
	"log"

	"github.com/google/go-github/v27/github"

	"github.com/ppapapetrou76/virtual-assistant/pkg/util/slices"
)

func ShouldRunOnIssue(event *github.IssuesEvent, configuredActions slices.StringSlice) bool {
	if addDefaultIfEmpty(configuredActions).HasString(*event.Action) {
		return true
	}
	log.Printf("Issues event is `%s` - eligible actions are `%v`. Skipping issues labeler", *event.Action, configuredActions)
	return false
}

func ShouldRunOnPullRequest(event *github.PullRequestEvent, configuredActions slices.StringSlice) bool {
	if addDefaultIfEmpty(configuredActions).HasString(*event.Action) {
		return true
	}
	log.Printf("Pull request event is `%s` - eligible actions are `%v`. Skipping issues labeler", *event.Action, configuredActions)
	return false
}

func addDefaultIfEmpty(actions slices.StringSlice) slices.StringSlice {
	if actions.IsEmpty() {
		return actions.Add("opened")
	}

	return actions
}
