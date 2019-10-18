package github

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/ppapapetrou76/virtual-assistant/pkg/util/slices"

	"github.com/ppapapetrou76/virtual-assistant/pkg/testutil"
)

func TestIssue_CurrentLabels(t *testing.T) {
	type fields struct {
		ghClient ClientWrapper
	}
	tests := []struct {
		name           string
		fields         fields
		wantErr        bool
		expectedError  error
		expectedLabels slices.StringSlice
	}{
		{
			name: "should return the issue labels",
			fields: fields{
				ghClient: MockGithubClient([]MockResponse{
					MockListIssueLabelsResponse(),
				}),
			},
			expectedLabels: []string{"bug", "enhancement"},
		},
		{
			name: "should error if labels cannot be loaded",
			fields: fields{
				ghClient: MockGithubClient([]MockResponse{
					{
						StatusCode: http.StatusOK,
						Response:   "ok",
					},
				}),
			},
			expectedError: errors.New("invalid character 'o' looking for beginning of value"),
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

			pr := Issue{
				Repo:   repo,
				Number: 0,
			}
			actualLabels, err := pr.CurrentLabels()
			testutil.AssertError(t, tt.wantErr, tt.expectedError, err)

			if !tt.wantErr && !reflect.DeepEqual(actualLabels, tt.expectedLabels) {
				t.Errorf("Expect: \n%+v Got: \n%+v", tt.expectedLabels, actualLabels)
			}
		})
	}
}

func TestNewIssue(t *testing.T) {
	repo := Repo{
		Owner: "ppapapetrou76",
		Name:  "virtual-assistant",
	}
	type args struct {
		repo   Repo
		number int
	}
	tests := []struct {
		name     string
		args     args
		expected Issue
	}{
		{
			name: "should return a new repo",
			args: args{
				repo:   repo,
				number: 123,
			},
			expected: Issue{
				Repo:   repo,
				Number: 123,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewIssue(tt.args.repo, tt.args.number)

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("Expect: \n%+v Got: \n%+v", tt.expected, actual)
			}
		})
	}
}

func TestIssue_ReplaceLabels(t *testing.T) {
	type fields struct {
		ghClient ClientWrapper
	}
	type args struct {
		labels []string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		expectedError  error
		expectedLabels []string
	}{
		{
			name: "should replace the issue labels",
			args: args{labels: []string{"bug", "enhancement"}},
			fields: fields{
				ghClient: MockGithubClient([]MockResponse{
					MockListIssueLabelsResponse(),
				}),
			},
		},
		{
			name: "should error if labels cannot be replaced",
			fields: fields{
				ghClient: MockGithubClient([]MockResponse{
					UnAuthorizedMockResponse(),
				}),
			},
			expectedError: errors.New("PUT https://api.github.com/repos/ppapapetrou76/virtual-assistant/issues/0/labels: 401 Bad credentials []"),
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

			pr := Issue{
				Repo:   repo,
				Number: 0,
			}
			err := pr.ReplaceLabels(tt.args.labels)
			testutil.AssertError(t, tt.wantErr, tt.expectedError, err)
		})
	}
}

func TestIssue_AtLeastOne(t *testing.T) {
	type fields struct {
		ghClient ClientWrapper
	}
	type args struct {
		labels       slices.StringSlice
		defaultLabel string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		expectedError  error
		expectedLabels []string
	}{
		{
			name: "should do nothing if labels group is empty",
			args: args{
				labels:       []string{},
				defaultLabel: "label1",
			},
			fields: fields{
				ghClient: MockGithubClient([]MockResponse{
					MockListIssueLabelsResponse(),
				}),
			},
		},
		{
			name: "should do nothing if default label is missing",
			args: args{
				labels: []string{"label1"},
			},
			fields: fields{
				ghClient: MockGithubClient([]MockResponse{
					MockListIssueLabelsResponse(),
				}),
			},
		},
		{
			name: "should error if current labels cannot be retrieved",
			args: args{
				labels:       []string{"label1", "label2"},
				defaultLabel: "label1",
			},
			fields: fields{
				ghClient: MockGithubClient([]MockResponse{
					UnAuthorizedMockResponse(),
				}),
			},
			expectedError: errors.New("GET https://api.github.com/repos/ppapapetrou76/virtual-assistant/issues/0/labels: 401 Bad credentials []"),
			wantErr:       true,
		},
		{
			name: "should do nothing if one of the labels group is already assigned to the github issue",
			args: args{
				labels:       []string{"bug", "label2"},
				defaultLabel: "label1",
			},
			fields: fields{
				ghClient: MockGithubClient([]MockResponse{
					MockListIssueLabelsResponse(),
				}),
			},
		},
		{
			name: "should add default label none of the labels group is assigned to the github issue",
			args: args{
				labels:       []string{"label1", "label2"},
				defaultLabel: "label1",
			},
			fields: fields{
				ghClient: MockGithubClient([]MockResponse{
					MockListIssueLabelsResponse(),
					MockGenericSuccessResponse(),
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := Repo{
				GHClient: tt.fields.ghClient,
				Owner:    "ppapapetrou76",
				Name:     "virtual-assistant",
			}

			pr := Issue{
				Repo:   repo,
				Number: 0,
			}
			err := pr.AtLeastOne(tt.args.labels, tt.args.defaultLabel)
			testutil.AssertError(t, tt.wantErr, tt.expectedError, err)
		})
	}
}

func TestIssue_AddToProject(t *testing.T) {
	type fields struct {
		ghClient ClientWrapper
	}
	type args struct {
		ProjectURL string
		column     string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		expectedError error
	}{
		{
			name: "should fail to add project if get issue fails",
			args: args{
				ProjectURL: "https://github.com/ppapapetrou76/virtual-assistant/projects/1",
				column:     "To Do",
			},
			fields: fields{
				ghClient: MockGithubClient([]MockResponse{
					UnAuthorizedMockResponse(),
				}),
			},
			wantErr:       true,
			expectedError: errors.New("cannot get issue with number 0. error message : GET https://api.github.com/repos/ppapapetrou76/virtual-assistant/issues/0: 401 Bad credentials []"),
		},
		{
			name: "should fail to add project if list of projects fails",
			args: args{
				ProjectURL: "https://github.com/ppapapetrou76/virtual-assistant/projects/1",
				column:     "To Do",
			},
			fields: fields{
				ghClient: MockGithubClient([]MockResponse{
					MockGetIssueResponse(),
					UnAuthorizedMockResponse(),
				}),
			},
			wantErr:       true,
			expectedError: errors.New("cannot get repository (ppapapetrou76/virtual-assistant) projects. error message : GET https://api.github.com/repos/ppapapetrou76/virtual-assistant/projects: 401 Bad credentials []"),
		},
		{
			name: "should fail to add project if list of project cards fails",
			args: args{
				ProjectURL: "https://github.com/ppapapetrou76/virtual-assistant/projects/1",
				column:     "To Do",
			},
			fields: fields{
				ghClient: MockGithubClient([]MockResponse{
					MockGetIssueResponse(),
					MockLisRepositoryProjectsResponse(),
					UnAuthorizedMockResponse(),
				}),
			},
			wantErr:       true,
			expectedError: errors.New("cannot get project (1002604) columns. error message : GET https://api.github.com/projects/1002604/columns: 401 Bad credentials []"),
		},
		{
			name: "should fail to add project if create project card fails",
			args: args{
				ProjectURL: "https://github.com/ppapapetrou76/virtual-assistant/projects/1",
				column:     "To Do",
			},
			fields: fields{
				ghClient: MockGithubClient([]MockResponse{
					MockGetIssueResponse(),
					MockLisRepositoryProjectsResponse(),
					MockListProjectColumnsResponse(),
					UnAuthorizedMockResponse(),
				}),
			},
			wantErr:       true,
			expectedError: errors.New("cannot add issue (0) to project (1002604). error message : POST https://api.github.com/projects/columns/367/cards: 401 Bad credentials []"),
		},
		{
			name: "should succeed to add project",
			args: args{
				ProjectURL: "https://github.com/ppapapetrou76/virtual-assistant/projects/1",
				column:     "To Do",
			},
			fields: fields{
				ghClient: MockGithubClient([]MockResponse{
					MockGetIssueResponse(),
					MockLisRepositoryProjectsResponse(),
					MockListProjectColumnsResponse(),
					MockGenericSuccessResponse(),
				}),
			},
		},
		{
			name: "should fail to add project if the given column doesn't exist in the project columns list",
			args: args{
				ProjectURL: "https://github.com/ppapapetrou76/virtual-assistant/projects/1",
				column:     "Invalid Column",
			},
			fields: fields{
				ghClient: MockGithubClient([]MockResponse{
					MockGetIssueResponse(),
					MockLisRepositoryProjectsResponse(),
					MockListProjectColumnsResponse(),
				}),
			},
			wantErr:       true,
			expectedError: errors.New("cannot add issue (0) to project (1002604). error message : no project columm found with name Invalid Column"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := Repo{
				GHClient: tt.fields.ghClient,
				Owner:    "ppapapetrou76",
				Name:     "virtual-assistant",
			}

			pr := Issue{
				Repo:   repo,
				Number: 0,
			}
			err := pr.AddToProject(tt.args.ProjectURL, tt.args.column)
			testutil.AssertError(t, tt.wantErr, tt.expectedError, err)
		})
	}
}
