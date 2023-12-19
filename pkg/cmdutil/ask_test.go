package cmdutil

import "testing"

func TestAsteriskStr(t *testing.T) {
	testcases := []struct {
		str  string
		want string
	}{
		{"", ""},
		{"t", "*"},
		{"tt", "t*"},
		{"ttt", "t*t"},
		{"tttt", "t**t"},
		{"tabt", "t**t"},
		{"tabct", "ta*ct"},
		{"abcdefg", "ab***fg"},
	}

	for _, tc := range testcases {
		if tc.want != AsteriskStr(tc.str) {
			t.Errorf("AsteriskStr(%q) = %q; want %q", tc.str, AsteriskStr(tc.str), tc.want)
		}
	}
}
