package gateway

import (
	"simple-balancer/balancer"
	"strings"
)

// URLMatcher takes a URL
type URLMatcher struct {
	trie *Trie[string, *balancer.LoadBalancer]
}

func NewURLMatcher() URLMatcher {
	return URLMatcher{
		trie: NewTrie[string, *balancer.LoadBalancer](),
	}
}

// Add splits the
func (m *URLMatcher) Add(path string, loadBalancer *balancer.LoadBalancer) {
	segments := MakeSegments(path)
	m.trie.Add(segments, loadBalancer)
}
func MakeSegments(path string) []string {
	if path == "/" {
		return []string{""}
	}
	segments := strings.Split(path, "/")
	if len(segments) == 0 {
		return []string{""}
	}
	return segments[:]
}

func (m *URLMatcher) Match(path string) (*balancer.LoadBalancer, bool) {
	segments := MakeSegments(path)
	lb, ok := m.trie.MatchLongestPrefix(segments)
	hasURL := lb != nil && ok
	return lb, hasURL
}
