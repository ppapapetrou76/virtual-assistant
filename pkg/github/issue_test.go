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
		owner    string
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
				owner: "ppapapetrou76",
				ghClient: MockGithubClient([]MockResponse{
					UnAuthorizedMockResponse(),
				}),
			},
			wantErr:       true,
			expectedError: errors.New("cannot get issue with number 0. error message : GET https://api.github.com/repos/ppapapetrou76/virtual-assistant/issues/0: 401 Bad credentials []"),
		},
		{
			name: "should fail to add project if list of repository projects fails",
			args: args{
				ProjectURL: "https://github.com/ppapapetrou76/virtual-assistant/projects/1",
				column:     "To Do",
			},
			fields: fields{
				owner: "ppapapetrou76",
				ghClient: MockGithubClient([]MockResponse{
					MockGetIssueResponse(),
					UnAuthorizedMockResponse(),
				}),
			},
			wantErr:       true,
			expectedError: errors.New("cannot get repository (ppapapetrou76/virtual-assistant) projects. error message : GET https://api.github.com/repos/ppapapetrou76/virtual-assistant/projects: 401 Bad credentials []"),
		},
		{
			name: "should fail to add project if list of organization projects fails",
			args: args{
				ProjectURL: "https://github.com/orgs/myorg/projects/1",
				column:     "To Do",
			},
			fields: fields{
				owner: "ppapapetrou76",
				ghClient: MockGithubClient([]MockResponse{
					MockGetIssueResponse(),
					MockListRepositoryProjectsResponse(),
					UnAuthorizedMockResponse(),
				}),
			},
			wantErr:       true,
			expectedError: errors.New("cannot get organization (ppapapetrou76) projects. error message : GET https://api.github.com/orgs/ppapapetrou76/projects: 401 Bad credentials []"),
		},
		{
			name: "should fail to add project if list of project cards fails",
			args: args{
				ProjectURL: "https://github.com/ppapapetrou76/virtual-assistant/projects/1",
				column:     "To Do",
			},
			fields: fields{
				owner: "ppapapetrou76",
				ghClient: MockGithubClient([]MockResponse{
					MockGetIssueResponse(),
					MockListRepositoryProjectsResponse(),
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
				owner: "ppapapetrou76",
				ghClient: MockGithubClient([]MockResponse{
					MockGetIssueResponse(),
					MockListRepositoryProjectsResponse(),
					MockListProjectColumnsResponse(),
					UnAuthorizedMockResponse(),
				}),
			},
			wantErr:       true,
			expectedError: errors.New("cannot add issue (0) to project (1002604). error message : POST https://api.github.com/projects/columns/367/cards: 401 Bad credentials []"),
		},
		{
			name: "should succeed to add personal project",
			args: args{
				ProjectURL: "https://github.com/ppapapetrou76/virtual-assistant/projects/1",
				column:     "To Do",
			},
			fields: fields{
				owner: "ppapapetrou76",
				ghClient: MockGithubClient([]MockResponse{
					MockGetIssueResponse(),
					MockListRepositoryProjectsResponse(),
					MockListProjectColumnsResponse(),
					MockGenericSuccessResponse(),
				}),
			},
		},
		{
			name: "should succeed to add organizational project",
			args: args{
				ProjectURL: "https://github.com/orgs/myorg/projects/1",
				column:     "To Do",
			},
			fields: fields{
				owner: "ppapapetrou76",
				ghClient: MockGithubClient([]MockResponse{
					MockGetIssueResponse(),
					MockListRepositoryProjectsResponse(),
					MockListOrganizationProjectsResponse(),
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
				owner: "ppapapetrou76",
				ghClient: MockGithubClient([]MockResponse{
					MockGetIssueResponse(),
					MockListRepositoryProjectsResponse(),
					MockListProjectColumnsResponse(),
				}),
			},
			wantErr:       true,
			expectedError: errors.New("cannot add issue (0) to project (1002604). error message : no project column found with name Invalid Column"),
		},
		{
			name: "should fail if no projects are found",
			args: args{
				ProjectURL: "https://github.com/orgs/myorg/projects/1",
				column:     "To Do",
			},
			fields: fields{
				owner: "myorg",
				ghClient: MockGithubClient([]MockResponse{
					MockGetIssueResponse(),
					MockListEmptyProjectsResponse(),
					MockListEmptyProjectsResponse(),
					MockListProjectColumnsResponse(),
					MockGenericSuccessResponse(),
				}),
			},
			wantErr:       true,
			expectedError: errors.New("no repository/organization (ppapapetrou76/virtual-assistant) projects found from the given url (https://github.com/orgs/myorg/projects/1)"),
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
