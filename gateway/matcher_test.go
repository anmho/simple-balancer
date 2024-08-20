package gateway

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"simple-balancer/balancer"
	"testing"
)

func Test_MakeSegments(t *testing.T) {
	tests := []struct {
		description string
		path        string

		expectedSegments []string
	}{
		{
			description:      "happy case: root path",
			path:             "/",
			expectedSegments: []string{""},
		},
		{
			description:      "happy case: single level",
			path:             "/posts",
			expectedSegments: []string{"", "posts"},
		},
		{
			description:      "happy case: two levels",
			path:             "/posts/123",
			expectedSegments: []string{"", "posts", "123"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			segments := MakeSegments(tc.path)
			assert.Equal(t, tc.expectedSegments, segments)
		})
	}
}

func TestURLMatcher_Match(t *testing.T) {
	matcher := NewURLMatcher()

	matcher.Add("/", balancer.New("/"))
	matcher.Add("/books", balancer.New("/books"))

	tests := []struct {
		description string
		path        string

		expectedShouldMatch bool
		expectedPath        string
	}{
		{
			description:         "happy path: root path when root present",
			path:                "/",
			expectedShouldMatch: true,
			expectedPath:        "/",
		},
		{
			description:         "happy path: exact path match",
			path:                "/books",
			expectedShouldMatch: true,
			expectedPath:        "/books",
		},
		{
			description:         "happy path: path longest prefix match",
			path:                "/books/123",
			expectedShouldMatch: true,
			expectedPath:        "/books",
		},
		{
			description:         "happy path: longest prefix matches root when root present",
			path:                "/authors",
			expectedShouldMatch: true,
			expectedPath:        "/",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description+" - "+tc.path, func(t *testing.T) {
			lb, ok := matcher.Match(tc.path)
			assert.Equal(t, tc.expectedShouldMatch, ok)
			require.NotNil(t, lb)
			assert.Equal(t, tc.expectedPath, lb.Path)
		})
	}
}
