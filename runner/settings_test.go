package runner

import "testing"

func TestLogColor(t *testing.T) {
	test := []struct {
		color    string
		expected string
	}{
		{color: "main", expected: "36"},
		{color: "build", expected: "33"},
		{color: "runner", expected: "32"},
		{color: "watcher", expected: "35"},
		{color: "app", expected: ""},
	}

	for _, v := range test {
		actual := logColor(v.color)
		if actual != v.expected {
			t.Errorf("Expected %v, got %v (for %q)", v.expected, actual, v.color)
		}
	}
}
