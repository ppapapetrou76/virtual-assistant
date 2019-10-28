package github

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v27/github"
)

// Repo is the struct to represent a github repository
type Repo struct {
	Owner, Name string
	GHClient    ClientWrapper
}

// NewRepo returns a new and properly initialized Repo struct
func NewRepo() Repo {
	t := strings.Split(os.Getenv(RepoEnvVar), "/")
	return Repo{
		Owner:    t[0],
		Name:     t[1],
		GHClient: DefaultClient(),
	}
}

// LoadFile loads a repo file and returns it in raw format (pointer of byte array)
func (r Repo) LoadFile(file, sha string) (*[]byte, error) {
	// ignore directory content and response as we don't need them here
	fileContent, _, _, err := r.GHClient.Repositories.GetContents(
		context.Background(),
		r.Owner,
		r.Name,
		file,
		&github.RepositoryContentGetOptions{Ref: sha})

	var content string
	if err == nil {
		content, err = fileContent.GetContent()
	}

	if err != nil {
		return nil, fmt.Errorf("load file : unable to load file from %s/%s@%s/%s: %w",
			r.Owner, r.Name, sha, file, err)
	}

	raw := []byte(content)
	return &raw, nil
}

// GetProjectID returns the id of a project given its url
func (r Repo) GetProjectID(projectURL string) (int64, error) {
	projects, _, err := r.GHClient.Repositories.ListProjects(context.Background(), r.Owner, r.Name, &github.ProjectListOptions{})
	if err != nil {
		return 0, fmt.Errorf("cannot get repository (%s/%s) projects. error message : %s", r.Owner, r.Name, err.Error())
	}

	var orgProjects []*github.Project
	if strings.Contains(projectURL, "orgs") {
		orgProjects, _, err = r.GHClient.Organizations.ListProjects(context.Background(), r.Owner, &github.ProjectListOptions{})
		if err != nil {
			return 0, fmt.Errorf("cannot get organization (%s) projects. error message : %s", r.Owner, err.Error())
		}
	}

	projects = append(projects, orgProjects...)

	var projectID int64
	for _, p := range projects {
		if *p.HTMLURL == projectURL {
			projectID = *p.ID
		}
	}

	if projectID == 0 {
		return 0, fmt.Errorf("no repository/organization (%s/%s) projects found from the given url (%s)", r.Owner, r.Name, projectURL)
	}

	return projectID, nil
}
