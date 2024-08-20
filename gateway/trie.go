package gateway

func NewTrieNode[K comparable, V any](segment K) *TrieNode[K, V] {
	return &TrieNode[K, V]{
		id:       segment,
		children: make(map[K]*TrieNode[K, V]),
		hasValue: false,
	}
}

func (n *TrieNode[K, V]) SetValue(value V) {
	n.Value = value
	n.hasValue = true
}

func (n *TrieNode[K, V]) RemoveValue() {
	n.Value = *new(V)
	n.hasValue = false
}

type TrieNode[K comparable, V any] struct {
	children map[K]*TrieNode[K, V]
	id       K
	Value    V
	hasValue bool
}

func (n *TrieNode[K, V]) GetChild(segment K) (*TrieNode[K, V], bool) {
	child, ok := n.children[segment]
	return child, ok
}

func (n *TrieNode[K, V]) HasValue() bool {
	return n.hasValue
}

type Trie[K comparable, V any] struct {
	root *TrieNode[K, V]
}

func NewTrie[K comparable, V any]() *Trie[K, V] {
	return &Trie[K, V]{
		root: NewTrieNode[K, V](*new(K)),
	}
}

func (t *Trie[K, V]) Add(key []K, value V) {
	cur := t.root

	for _, k := range key {
		_, hasChild := cur.GetChild(k)
		if !hasChild {
			cur.children[k] = NewTrieNode[K, V](k)
		}
		node, _ := cur.GetChild(k)
		cur = node
	}

	cur.SetValue(value)
}

func (t *Trie[K, V]) Remove(key []K) bool {
	return false
}

func (t *Trie[K, V]) MatchLongestPrefix(key []K) (V, bool) {
	var res *TrieNode[K, V]
	if t.root.HasValue() {
		res = t.root
	}
	cur := t.root

	for i := range len(key) {
		node, ok := cur.GetChild(key[i])
		if !ok {
			break
		}
		if node.HasValue() {
			res = node
		}
		cur = node
	}
	if res == nil {
		return *new(V), false
	}
	return res.Value, res.HasValue()
}

func (t *Trie[K, V]) Get(key []K) (V, bool) {
	cur := t.root
	for i := range len(key) {
		node, ok := cur.GetChild(key[i])
		if !ok {
			return *new(V), false
		}
		cur = node
	}

	return cur.Value, cur.HasValue()
}
