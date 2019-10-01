package github

import (
	"errors"
	"os"
	"reflect"
	"testing"

	testutil "github.com/ppapapetrou76/virtual-assistant/pkg/util"
)

func TestNewRepo(t *testing.T) {
	os.Setenv(TokenEnvVar, "some-token")
	ghClient := DefaultClient()

	type fields struct {
		repo  string
		token string
	}
	type args struct {
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected Repo
	}{
		{
			name: "should return a new repo",
			fields: fields{
				repo:  "ppapapetrou76/virtual-assistant",
				token: "some-token",
			},
			expected: Repo{
				Owner:    "ppapapetrou76",
				Name:     "virtual-assistant",
				GHClient: ghClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(RepoEnvVar, tt.fields.repo)

			actualRepo := NewRepo()
			if !reflect.DeepEqual(actualRepo, tt.expected) {
				t.Errorf("Expect: \n%+v Got: \n%+v", tt.expected, actualRepo)
			}
		})
	}
}

const getContentResponse = `{
  "type": "file",
  "encoding": "base64",
  "size": 642,
  "name": "README.md",
  "path": "README.md",
  "content": "SSdtIGEgc2VjcmV0"
}`

func TestRepo_LoadFile(t *testing.T) {
	type fields struct {
		ghClient ClientWrapper
	}
	tests := []struct {
		name             string
		fields           fields
		wantErr          bool
		expectedError    error
		expectedContents []byte
	}{
		{
			name: "should return the file contents",
			fields: fields{
				ghClient: MockGithubClient(200, getContentResponse),
			},
			expectedContents: []byte("I'm a secret"),
		},
		{
			name: "should error if file cannot be loaded",
			fields: fields{
				ghClient: MockGithubClient(200, "ok"),
			},
			expectedError: errors.New("load file : unable to load file from ppapapetrou76/virtual-assistant@/some-file: invalid character 'o' looking for beginning of value"),
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := Repo{
				GHClient: tt.fields.ghClient,
				Owner:    "ppapapetrou76",
				Name:     "virtual-assistant",
			}
			actualContent, err := repo.LoadFile("some-file", "")
			testutil.AssertError(t, tt.wantErr, tt.expectedError, err)

			if !tt.wantErr && !reflect.DeepEqual(*actualContent, tt.expectedContents) {
				t.Errorf("Expect: \n%+v Got: \n%+v", tt.expectedContents, *actualContent)
			}
		})
	}
}
