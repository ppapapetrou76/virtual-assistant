package testutil

import "testing"

// AssertError asserts the expected error message of a given test
func AssertError(t *testing.T, wantError bool, expected, actual error) {
	if (actual != nil) != wantError {
		t.Errorf("%s error = %v, wantErr %v", t.Name(), actual, wantError)
		return
	}

	if wantError && actual == nil {
		t.Errorf("%s expected errors = '%v' but no errors returned", t.Name(), expected)
	}

	if wantError && actual.Error() != expected.Error() {
		t.Errorf("%s expected errors = \n%v but got \n%v", t.Name(), expected, actual)
	}

	if !wantError && actual != nil {
		t.Errorf("%s expected no errors but got %v", t.Name(), actual)
	}
}
