package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/go-yaml/yaml"

	"github.com/ppapapetrou76/virtual-assistant/pkg/testutil"
)

func TestLoad(t *testing.T) {
	type fields struct {
		fileName string
	}
	type args struct {
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		expected    *Config
		wantErr     bool
		expectedErr error
	}{
		{
			name: "should load config",
			fields: fields{
				fileName: "../../test_data/valid-config.yml",
			},
			expected: &Config{
				AssignerConfig: AssignerConfig{
					IssuesAssignerConfig: IssuesAssignerConfig{
						Actions: []string{
							"opened",
							"milestoned",
						},
						IssuesAssignerProjectConfig: IssuesAssignerProjectConfig{
							ProjectURL: "https://github.com/ppapapetrou76/virtual-assistant/projects/1",
							Column:     "To Do",
						},
					},
				},
				LabelerConfig: LabelerConfig{
					IssuesLabelerConfig: IssuesLabelerConfig{
						Labels: []string{
							"label1",
							"label2",
							"area:label3",
						},
						Actions: []string{
							"opened",
							"milestoned",
						},
						OneOfaKind: OneOfaKind{
							PossibleLabels: []string{
								"priority:1",
								"priority:2",
								"priority:3",
							},
							Default: "priority:2",
						},
					},
					PullRequestsLabelerConfig: PullRequestsLabelerConfig{
						Labels: []string{
							"label1",
							"label2",
						},
						Actions: []string{
							"opened",
							"synchronize",
						},
					},
				},
			},
		},
		{
			name: "should error if byte array is empty",
			fields: fields{
				fileName: "../../test_data/phantom.yml",
			},
			wantErr:     true,
			expectedErr: errors.New("load config : unable to un-marshall empty byte array"),
			expected:    &Config{},
		},
		{
			name: "should error if byte array contains invalid data",
			fields: fields{
				fileName: "../../test_data/invalid-config.yml",
			},
			wantErr: true,
			expectedErr: fmt.Errorf(
				"load config : unable to un-marshall config [%v], %w", string(*getContents("../../test_data/invalid-config.yml")),
				&yaml.TypeError{Errors: []string{"line 1: cannot unmarshal !!str `labels ...` into config.Config"}}),
			expected: &Config{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, actualErr := Load(getContents(tt.fields.fileName))
			testutil.AssertError(t, tt.wantErr, tt.expectedErr, actualErr)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("Expect: \n%+v Got: \n%+v", tt.expected, actual)
			}
		})
	}
}

func getContents(filename string) *[]byte {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}

	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}

	return &contents
}
