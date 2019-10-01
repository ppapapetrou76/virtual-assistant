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
