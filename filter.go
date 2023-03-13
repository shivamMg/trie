package trie

type FilterFunc func(interface{}) bool

// Filter: find item based on value.
// apply eval func to each Value of the tree
// build a result with key, value for each itme eval == true
func (t *Trie) Filter(filter FilterFunc) *SearchResults {
	results := &SearchResults{}
	node := t.root
	t.filter(results, node, &[]string{}, filter)
	return results
}

// filter: recursively apply eval an item tree
func (t *Trie) filter(results *SearchResults, node *Node, prefixKey *[]string, filter FilterFunc) *SearchResults {

	if node.isTerminal && filter(node.value) {
		key := make([]string, len(*prefixKey))
		copy(key, *prefixKey)
		result := &SearchResult{Key: key, Value: node.value}
		results.Results = append(results.Results, result)
	}

	for dllNode := node.childrenDLL.head; dllNode != nil; dllNode = dllNode.next {
		child := dllNode.trieNode
		key := child.keyPart
		pfx := append(*prefixKey, key)

		t.filter(results, child, &pfx, filter)

	}
	return results
}
