package gateway

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrie_Get(t *testing.T) {
	trie := NewTrie[string, int]()

	key := []string{"1", "2", "3"}
	trie.Add(key, 10)

	value, ok := trie.Get(key)
	assert.True(t, ok)
	assert.Equal(t, 10, value)

	value, ok = trie.Get([]string{"1", "2"})
	assert.False(t, ok)
}

func TestTrie_MatchLongestPrefix(t *testing.T) {

	tests := []struct {
		description string
		key         []string
		value       int

		queryKey       []string
		expectedHasKey bool
		expectedValue  int
	}{
		{
			description: "happy case: key is a sub-path of the input key",
			key:         []string{"1", "2", "3"},
			value:       10,

			queryKey:       []string{"1", "2", "3", "4"},
			expectedHasKey: true,
			expectedValue:  10,
		},
		{
			description: "happy case: empty key should match any key query",
			key:         []string{},
			value:       10,

			queryKey:       []string{"1", "2", "3", "4"},
			expectedHasKey: true,
			expectedValue:  10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			trie := NewTrie[string, int]()
			trie.Add(tc.key, tc.value)

			value, ok := trie.MatchLongestPrefix(tc.queryKey)
			assert.Equal(t, tc.expectedValue, value)
			assert.Equal(t, tc.expectedHasKey, ok)

		})
	}
}
