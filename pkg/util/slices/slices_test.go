package slices

import (
	"reflect"
	"testing"
)

func TestStringSlice_IsEmpty(t *testing.T) {
	type fields struct {
		slice StringSlice
	}
	tests := []struct {
		name     string
		expected bool
		fields   fields
	}{
		{
			name:     "should return true if slice is empty",
			fields:   fields{slice: StringSlice{}},
			expected: true,
		},
		{
			name:   "should return false if slice is not empty",
			fields: fields{slice: StringSlice{"some value"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.fields.slice.IsEmpty()
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("Expect: \n%+v Got: \n%+v", tt.expected, actual)
			}
		})
	}
}

func TestStringSlice_HasString(t *testing.T) {
	type fields struct {
		slice StringSlice
	}
	type args struct {
		value string
	}
	tests := []struct {
		name     string
		expected bool
		fields   fields
		args     args
	}{
		{
			name:     "should return true if slice has string",
			fields:   fields{slice: StringSlice{"value", "anothervalue"}},
			expected: true,
			args:     args{value: "value"},
		},
		{
			name:   "should return false if slice doesn't have string",
			fields: fields{slice: StringSlice{"value", "anothervalue"}},
			args:   args{value: "random"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.fields.slice.HasString(tt.args.value)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("Expect: \n%+v Got: \n%+v", tt.expected, actual)
			}
		})
	}
}
