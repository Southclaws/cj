package commands

import "testing"

func TestNormalizeSearchQuery(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{name: "plain", query: "exact phrase", want: "exact phrase"},
		{name: "double quoted", query: `"exact phrase"`, want: "exact phrase"},
		{name: "single quoted", query: `'exact phrase'`, want: "exact phrase"},
		{name: "whitespace", query: "  message text  ", want: "message text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeSearchQuery(tt.query); got != tt.want {
				t.Errorf("normalizeSearchQuery() = %q, want %q", got, tt.want)
			}
		})
	}
}
