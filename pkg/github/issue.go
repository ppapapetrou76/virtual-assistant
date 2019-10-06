package github

import (
	"context"
	"log"

	"github.com/google/go-github/v27/github"

	"github.com/ppapapetrou76/virtual-assistant/pkg/util/slices"
)

// Issue is the struct to represent a github pull request
type Issue struct {
	Repo
	Number int
}

// ReplaceLabels replace the labels of the issue/pull request with the ones passed as method argument
func (pr Issue) ReplaceLabels(labels []string) error {
	log.Printf("Setting labels to %s/%s#%d: %s", pr.Owner, pr.Name, pr.Number, labels)
	_, _, err := pr.GHClient.Issues.ReplaceLabelsForIssue(
		context.Background(), pr.Owner, pr.Name, pr.Number, labels)
	return err
}

// AtLeastOne replace the labels of the issue/pull request with the ones passed as method argument
func (pr Issue) AtLeastOne(labels slices.StringSlice, defaultLabel string) error {
	if labels.IsEmpty() || defaultLabel == "" {
		return nil
	}

	currentLabels, err := pr.CurrentLabels()
	if err != nil {
		return err
	}
	for _, label := range currentLabels {
		if labels.HasString(label) {
			return nil
		}
	}
	desiredLabels := append(currentLabels, defaultLabel)
	log.Printf("Setting labels to %s/%s#%d: %s", pr.Owner, pr.Name, pr.Number, desiredLabels)
	_, _, err = pr.GHClient.Issues.ReplaceLabelsForIssue(
		context.Background(), pr.Owner, pr.Name, pr.Number, desiredLabels)
	return err
}

// CurrentLabels returns the current labels of an issue/pull request
func (pr Issue) CurrentLabels() (slices.StringSlice, error) {
	opts := github.ListOptions{}
	currLabels, _, err := pr.GHClient.Issues.ListLabelsByIssue(
		context.Background(), pr.Owner, pr.Name, pr.Number, &opts)

	labels := make([]string, 0, len(currLabels))
	for _, label := range currLabels {
		labels = append(labels, *label.Name)
	}
	return labels, err
}

// NewIssue returns a new Issue struct
func NewIssue(r Repo, number int) Issue {
	return Issue{
		Repo:   r,
		Number: number,
	}
}
