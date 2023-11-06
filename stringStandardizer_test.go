package ryde

import (
	"reflect"
	"testing"
)

func TestStandardizeSpaces(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"  hello  world  ", "hello world"},
		{"\tfoo\tbar\t", "foo bar"},
		{"\tfoo\rbar\t", "foo bar"},
		{"\n\n\nbaz\n\n", "baz"},
		{"", ""},
	}

	for _, test := range tests {
		if got := StandardizeString(test.input); got != test.want {
			t.Errorf("StandardizeSpaces(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}
func TestStandardizeStringSlice(t *testing.T) {
	var tests = []struct {
		input []string
		want  []string
	}{
		{[]string{"  hello  world  ", "\tfoo\tbar\t", "\tfoo\rbar\t", "\n\n\nbaz\n\n", ""}, []string{"hello world", "foo bar", "foo bar", "baz", ""}},
		{[]string{"  hello  world  ", "  foo  bar  ", "  baz  "}, []string{"hello world", "foo bar", "baz"}},
		{[]string{"", "", ""}, []string{"", "", ""}},
	}

	for _, test := range tests {
		got := StandardizeStringSlice(test.input)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("StandardizeStringSlice(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}
