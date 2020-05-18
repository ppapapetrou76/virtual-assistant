package github

import (
	"context"
	"fmt"
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
func (i Issue) ReplaceLabels(labels []string) error {
	log.Printf("Setting labels to %s/%s#%d: %s", i.Owner, i.Name, i.Number, labels)
	_, _, err := i.GHClient.Issues.ReplaceLabelsForIssue(
		context.Background(), i.Owner, i.Name, i.Number, labels)
	return err
}

// AtLeastOne replace the labels of the issue/pull request with the ones passed as method argument
func (i Issue) AtLeastOne(labels slices.StringSlice, defaultLabel string) error {
	if labels.IsEmpty() || defaultLabel == "" {
		return nil
	}

	currentLabels, err := i.CurrentLabels()
	if err != nil {
		return err
	}
	for _, label := range currentLabels {
		if labels.HasString(label) {
			return nil
		}
	}
	desiredLabels := append(currentLabels, defaultLabel)
	log.Printf("Setting labels to %s/%s#%d: %s", i.Owner, i.Name, i.Number, desiredLabels)
	_, _, err = i.GHClient.Issues.ReplaceLabelsForIssue(
		context.Background(), i.Owner, i.Name, i.Number, desiredLabels)
	return err
}

// CurrentLabels returns the current labels of an issue/pull request
func (i Issue) CurrentLabels() (slices.StringSlice, error) {
	opts := github.ListOptions{}
	currLabels, _, err := i.GHClient.Issues.ListLabelsByIssue(
		context.Background(), i.Owner, i.Name, i.Number, &opts)

	labels := make([]string, 0, len(currLabels))
	for _, label := range currLabels {
		labels = append(labels, *label.Name)
	}
	return labels, err
}

// AddAssignee adds the user who created the issue/PR as assignee
func (i Issue) AddAssignee() error {
	log.Printf("Assigning the PR/Issue to the user who created it")
	issue, _, err := i.GHClient.Issues.Get(context.Background(), i.Owner, i.Name, i.Number)
	if err != nil {
		return fmt.Errorf("cannot get issue with number %d. error message : %s", i.Number, err.Error())
	}
	_, _, err = i.GHClient.Issues.AddAssignees(context.Background(), i.Owner, i.Name, i.Number, []string{*issue.User.Login})
	return err
}

// AddToProject adds the issue to the given project. If the project doesn't exist it returns an error
func (i Issue) AddToProject(projectURL, column string) error {
	log.Printf("Adding to project %s in column %s", projectURL, column)
	issue, _, err := i.GHClient.Issues.Get(context.Background(), i.Owner, i.Name, i.Number)
	if err != nil {
		return fmt.Errorf("cannot get issue with number %d. error message : %s", i.Number, err.Error())
	}

	projectID, err := i.Repo.GetProjectID(projectURL)
	if err != nil {
		return err
	}
	opts := &github.ProjectCardOptions{
		ContentType: "Issue",
		ContentID:   *issue.ID,
	}
	columns, _, err := i.GHClient.Projects.ListProjectColumns(context.Background(), projectID, &github.ListOptions{})
	if err != nil {
		return fmt.Errorf("cannot get project (%d) columns. error message : %s", projectID, err.Error())
	}

	for _, c := range columns {
		if *c.Name == column {
			_, _, err := i.GHClient.Projects.CreateProjectCard(context.Background(), *c.ID, opts)
			if err != nil {
				return fmt.Errorf("cannot add issue (%d) to project (%d). error message : %s",
					i.Number, projectID, err.Error())
			}
			return nil
		}
	}

	return fmt.Errorf("cannot add issue (%d) to project (%d). error message : no project column found with name %s",
		i.Number, projectID, column)
}

// NewIssue returns a new Issue struct
func NewIssue(r Repo, number int) Issue {
	return Issue{
		Repo:   r,
		Number: number,
	}
}
