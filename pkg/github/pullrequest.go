package github

import (
	"context"
	"log"

	"github.com/google/go-github/v27/github"
)

// PullRequest is the struct to represent a github pull request
type PullRequest struct {
	Repo
	Number int
}

// ReplaceLabels replace the labels of the pull request with the ones passed as method argument
func (pr PullRequest) ReplaceLabels(labels []string) error {
	log.Printf("Setting labels to %s/%s#%d: %s", pr.Owner, pr.Name, pr.Number, labels)
	_, _, err := pr.GHClient.Issues.ReplaceLabelsForIssue(
		context.Background(), pr.Owner, pr.Name, pr.Number, labels)
	return err
}

// CurrentLabels returns the current labels of a pull request
func (pr PullRequest) CurrentLabels() ([]string, error) {
	opts := github.ListOptions{}
	currLabels, _, err := pr.GHClient.Issues.ListLabelsByIssue(
		context.Background(), pr.Owner, pr.Name, pr.Number, &opts)

	labels := make([]string, 0, len(currLabels))
	for _, label := range currLabels {
		labels = append(labels, *label.Name)
	}
	return labels, err
}

// NewPullRequest returns a new PullRequest struct
func NewPullRequest(r Repo, number int) PullRequest {
	return PullRequest{
		Repo:   r,
		Number: number,
	}
}
