package trie

import (
	"container/heap"
	"errors"
	"math"
)

type EditOpType int

const (
	EditOpTypeNoEdit EditOpType = iota
	EditOpTypeInsert
	EditOpTypeDelete
	EditOpTypeReplace
)

// EditOp represents an Edit Operation.
type EditOp struct {
	Type EditOpType
	// KeyPart:
	// - In case of NoEdit, KeyPart is to be retained.
	// - In case of Insert, KeyPart is to be inserted in the key.
	// - In case of Delete/Replace, KeyPart is the part of the key on which delete/replace is performed.
	KeyPart string
	// ReplaceWith is set for Type=EditOpTypeReplace
	ReplaceWith string
}

type SearchResults struct {
	Results         []*SearchResult
	heap            *searchResultMaxHeap
	tiebreakerCount int
}

type SearchResult struct {
	// Key is the key that was Put() into the Trie.
	Key []string
	// Value is the value that was Put() into the Trie.
	Value interface{}
	// EditDistance is the number of edits (insert/delete/replace) needed to convert Key into the Search()-ed key.
	EditDistance int
	// EditOps is the list of edit operations (see EditOpType) needed to convert Key into the Search()-ed key.
	EditOps []*EditOp

	tiebreaker int
}

type SearchOptions struct {
	// - WithExactKey
	// - WithMaxResults
	// - WithMaxEditDistance
	//   - WithEditOps
	//   - WithTopKLeastEdited
	exactKey        bool
	maxResults      bool
	maxResultsCount int
	editDistance    bool
	maxEditDistance int
	editOps         bool
	topKLeastEdited bool
}

// WithExactKey can be passed to Search(). When passed, Search() returns just the result with
// Key=Search()-ed key. If the key does not exist, result list will be empty.
func WithExactKey() func(*SearchOptions) {
	return func(so *SearchOptions) {
		so.exactKey = true
	}
}

// WithMaxResults can be passed to Search(). When passed, Search() will return at most maxResults
// number of results.
func WithMaxResults(maxResults int) func(*SearchOptions) {
	if maxResults <= 0 {
		panic(errors.New("invalid usage: maxResults must be greater than zero"))
	}
	return func(so *SearchOptions) {
		so.maxResults = true
		so.maxResultsCount = maxResults
	}
}

// WithMaxEditDistance can be passed to Search(). When passed, Search() changes its default behaviour from
// Prefix search to Edit distance search. It can be used to return "Approximate" results instead of strict
// Prefix search results.
//
// maxDistance is the maximum number of edits allowed on Trie keys to consider them as a SearchResult.
// Higher the maxDistance, more lenient and slower the search becomes.
//
// e.g. If a Trie stores English words, then searching for "wheat" with maxDistance=1 might return similar
// looking words like "wheat", "cheat", "heat", "what", etc. With maxDistance=2 it might also return words like
// "beat", "ahead", etc.
//
// Read about Edit distance: https://en.wikipedia.org/wiki/Edit_distance
func WithMaxEditDistance(maxDistance int) func(*SearchOptions) {
	if maxDistance <= 0 {
		panic(errors.New("invalid usage: maxDistance must be greater than zero"))
	}
	return func(so *SearchOptions) {
		so.editDistance = true
		so.maxEditDistance = maxDistance
	}
}

// WithEditOps can be passed to Search() alongside WithMaxEditDistance(). When passed, Search() also returns EditOps
// for each SearchResult. EditOps can be used to determine the minimum number of edit operations needed to convert
// a result Key into the Search()-ed key.
//
// e.g. Searching for "wheat" in a Trie that stores English words might return "beat". EditOps for this result might be:
// 1. insert "w" 2. replace "b" with "h".
//
// There might be multiple ways to edit a key into another. EditOps represents only one.
//
// Computing EditOps makes Search() slower.
func WithEditOps() func(*SearchOptions) {
	return func(so *SearchOptions) {
		so.editOps = true
	}
}

// WithTopKLeastEdited can be passed to Search() alongside WithMaxEditDistance() and WithMaxResults(). When passed,
// Search() returns maxResults number of results that have the lowest EditDistances. Results are sorted on EditDistance
// (lowest to highest).
//
// e.g. In a Trie that stores English words searching for "wheat" might return "wheat" (EditDistance=0), "cheat" (EditDistance=1),
// "beat" (EditDistance=2) - in that order.
func WithTopKLeastEdited() func(*SearchOptions) {
	return func(so *SearchOptions) {
		so.topKLeastEdited = true
	}
}

// Search() takes a key and some options to return results (see SearchResult) from the Trie.
// Without any options, it does a Prefix search i.e. result Keys have the same prefix as key.
// Order of the results is deterministic and will follow the order in which Put() was called for the keys.
// See "With*" functions for options accepted by Search().
func (t *Trie) Search(key []string, options ...func(*SearchOptions)) *SearchResults {
	opts := &SearchOptions{}
	for _, f := range options {
		f(opts)
	}
	if opts.editOps && !opts.editDistance {
		panic(errors.New("invalid usage: WithEditOps() must not be passed without WithMaxEditDistance()"))
	}
	if opts.topKLeastEdited && !opts.editDistance {
		panic(errors.New("invalid usage: WithTopKLeastEdited() must not be passed without WithMaxEditDistance()"))
	}
	if opts.exactKey && opts.editDistance {
		panic(errors.New("invalid usage: WithExactKey() must not be passed with WithMaxEditDistance()"))
	}
	if opts.exactKey && opts.maxResults {
		panic(errors.New("invalid usage: WithExactKey() must not be passed with WithMaxResults()"))
	}
	if opts.topKLeastEdited && !opts.maxResults {
		panic(errors.New("invalid usage: WithTopKLeastEdited() must not be passed without WithMaxResults()"))
	}

	if opts.editDistance {
		return t.searchWithEditDistance(key, opts)
	}
	return t.search(key, opts)
}

func (t *Trie) searchWithEditDistance(key []string, opts *SearchOptions) *SearchResults {
	// https://en.wikipedia.org/wiki/Levenshtein_distance#Iterative_with_full_matrix
	// http://stevehanov.ca/blog/?id=114
	columns := len(key) + 1
	newRow := make([]int, columns)
	for i := 0; i < columns; i++ {
		newRow[i] = i
	}
	m := len(key)
	if m == 0 {
		m = 1
	}
	rows := make([][]int, 1, m)
	rows[0] = newRow
	results := &SearchResults{}
	if opts.topKLeastEdited {
		results.heap = &searchResultMaxHeap{}
	}

	keyColumn := make([]string, 1, m)
	stop := false
	// prioritize Node that has the same keyPart as key. this results in better results
	// e.g. if key=national, build with Node(keyPart=n) first so that keys like notional, nation, nationally, etc. are prioritized
	// same logic is used inside the recursive buildWithEditDistance() method
	var prioritizedNode *Node
	if len(key) > 0 {
		if prioritizedNode = t.root.children[key[0]]; prioritizedNode != nil {
			keyColumn[0] = prioritizedNode.keyPart
			t.buildWithEditDistance(&stop, results, prioritizedNode, &keyColumn, &rows, key, opts)
		}
	}
	for dllNode := t.root.childrenDLL.head; dllNode != nil; dllNode = dllNode.next {
		node := dllNode.trieNode
		if node == prioritizedNode {
			continue
		}
		keyColumn[0] = node.keyPart
		t.buildWithEditDistance(&stop, results, node, &keyColumn, &rows, key, opts)
	}
	if opts.topKLeastEdited {
		n := results.heap.Len()
		results.Results = make([]*SearchResult, n)
		for n != 0 {
			result := heap.Pop(results.heap).(*SearchResult)
			result.tiebreaker = 0
			results.Results[n-1] = result
			n--
		}
		results.heap = nil
		results.tiebreakerCount = 0
	}
	return results
}

func (t *Trie) buildWithEditDistance(stop *bool, results *SearchResults, node *Node, keyColumn *[]string, rows *[][]int, key []string, opts *SearchOptions) {
	if *stop {
		return
	}
	prevRow := (*rows)[len(*rows)-1]
	columns := len(key) + 1
	newRow := make([]int, columns)
	newRow[0] = prevRow[0] + 1
	for i := 1; i < columns; i++ {
		replaceCost := 1
		if key[i-1] == (*keyColumn)[len(*keyColumn)-1] {
			replaceCost = 0
		}
		newRow[i] = min(
			newRow[i-1]+1,            // insertion
			prevRow[i]+1,             // deletion
			prevRow[i-1]+replaceCost, // substitution
		)
	}
	*rows = append(*rows, newRow)

	if newRow[columns-1] <= opts.maxEditDistance && node.isTerminal {
		editDistance := newRow[columns-1]
		lazyCreate := func() *SearchResult { // optimization for the case where topKLeastEdited=true and the result should not be pushed to heap
			resultKey := make([]string, len(*keyColumn))
			copy(resultKey, *keyColumn)
			result := &SearchResult{Key: resultKey, Value: node.value, EditDistance: editDistance}
			if opts.editOps {
				result.EditOps = t.getEditOps(rows, keyColumn, key)
			}
			return result
		}
		if opts.topKLeastEdited {
			results.tiebreakerCount++
			if results.heap.Len() < opts.maxResultsCount {
				result := lazyCreate()
				result.tiebreaker = results.tiebreakerCount
				heap.Push(results.heap, result)
			} else if (*results.heap)[0].EditDistance > editDistance {
				result := lazyCreate()
				result.tiebreaker = results.tiebreakerCount
				heap.Pop(results.heap)
				heap.Push(results.heap, result)
			}
		} else {
			result := lazyCreate()
			results.Results = append(results.Results, result)
			if opts.maxResults && len(results.Results) == opts.maxResultsCount {
				*stop = true
				return
			}
		}
	}

	if min(newRow...) <= opts.maxEditDistance {
		var prioritizedNode *Node
		m := len(*keyColumn)
		if m < len(key) {
			if prioritizedNode = node.children[key[m]]; prioritizedNode != nil {
				*keyColumn = append(*keyColumn, prioritizedNode.keyPart)
				t.buildWithEditDistance(stop, results, prioritizedNode, keyColumn, rows, key, opts)
				*keyColumn = (*keyColumn)[:len(*keyColumn)-1]
			}
		}
		for dllNode := node.childrenDLL.head; dllNode != nil; dllNode = dllNode.next {
			child := dllNode.trieNode
			if child == prioritizedNode {
				continue
			}
			*keyColumn = append(*keyColumn, child.keyPart)
			t.buildWithEditDistance(stop, results, child, keyColumn, rows, key, opts)
			*keyColumn = (*keyColumn)[:len(*keyColumn)-1]
		}
	}

	*rows = (*rows)[:len(*rows)-1]
}

func (t *Trie) getEditOps(rows *[][]int, keyColumn *[]string, key []string) []*EditOp {
	// https://gist.github.com/jlherren/d97839b1276b9bd7faa930f74711a4b6
	ops := make([]*EditOp, 0, len(key))
	r, c := len(*rows)-1, len((*rows)[0])-1
	for r > 0 || c > 0 {
		insertionCost, deletionCost, substitutionCost := math.MaxInt, math.MaxInt, math.MaxInt
		if c > 0 {
			insertionCost = (*rows)[r][c-1]
		}
		if r > 0 {
			deletionCost = (*rows)[r-1][c]
		}
		if r > 0 && c > 0 {
			substitutionCost = (*rows)[r-1][c-1]
		}
		minCost := min(insertionCost, deletionCost, substitutionCost)
		if minCost == substitutionCost {
			if (*rows)[r][c] > (*rows)[r-1][c-1] {
				ops = append(ops, &EditOp{Type: EditOpTypeReplace, KeyPart: (*keyColumn)[r-1], ReplaceWith: key[c-1]})
			} else {
				ops = append(ops, &EditOp{Type: EditOpTypeNoEdit, KeyPart: (*keyColumn)[r-1]})
			}
			r -= 1
			c -= 1
		} else if minCost == deletionCost {
			ops = append(ops, &EditOp{Type: EditOpTypeDelete, KeyPart: (*keyColumn)[r-1]})
			r -= 1
		} else if minCost == insertionCost {
			ops = append(ops, &EditOp{Type: EditOpTypeInsert, KeyPart: key[c-1]})
			c -= 1
		}
	}
	for i, j := 0, len(ops)-1; i < j; i, j = i+1, j-1 {
		ops[i], ops[j] = ops[j], ops[i]
	}
	return ops
}

func (t *Trie) search(prefixKey []string, opts *SearchOptions) *SearchResults {
	results := &SearchResults{}
	node := t.root
	for _, keyPart := range prefixKey {
		child, ok := node.children[keyPart]
		if !ok {
			return results
		}
		node = child
	}
	if opts.exactKey {
		if node.isTerminal {
			result := &SearchResult{Key: prefixKey, Value: node.value}
			results.Results = append(results.Results, result)
		}
		return results
	}
	t.build(results, node, &prefixKey, opts)
	return results
}

// SelectOnValue: find item based on value.
// apply eval func to each Value of the tree
// build a result with key, value for each itme eval == true
func (t *Trie) SelectOnValue(eval func(interface{}) bool) *SearchResults {
	results := &SearchResults{}
	node := t.root
	t.selectOnValue(results, node, []string{}, eval)
	return results
}

// selectOnValue: recursively apply eval an item tree
func (t *Trie) selectOnValue(results *SearchResults, node *Node, prefixKey []string, eval func(interface{}) bool) *SearchResults {

	for key, child := range node.children {
		pfx := append(prefixKey, key)
		if child.isTerminal {
			if eval(child.value) {
				t.scanBuild(results, child, &pfx)
			}
			t.selectOnValue(results, child, pfx, eval)

		} else {
			t.selectOnValue(results, child, pfx, eval)
		}
	}
	return results
}

// scanBuild: build a result off scanned items
func (t *Trie) scanBuild(results *SearchResults, node *Node, prefixKey *[]string) (stop bool) {
	if node.isTerminal {
		key := make([]string, len(*prefixKey))
		copy(key, *prefixKey)
		result := &SearchResult{Key: key, Value: node.value}
		results.Results = append(results.Results, result)
	}

	return false
}

func (t *Trie) build(results *SearchResults, node *Node, prefixKey *[]string, opts *SearchOptions) (stop bool) {
	if node.isTerminal {
		key := make([]string, len(*prefixKey))
		copy(key, *prefixKey)
		result := &SearchResult{Key: key, Value: node.value}
		results.Results = append(results.Results, result)
		if opts.maxResults && len(results.Results) == opts.maxResultsCount {
			return true
		}
	}

	for dllNode := node.childrenDLL.head; dllNode != nil; dllNode = dllNode.next {
		child := dllNode.trieNode
		*prefixKey = append(*prefixKey, child.keyPart)
		stop := t.build(results, child, prefixKey, opts)
		*prefixKey = (*prefixKey)[:len(*prefixKey)-1]
		if stop {
			return true
		}
	}
	return false
}

func min(s ...int) int {
	m := s[0]
	for _, a := range s[1:] {
		if a < m {
			m = a
		}
	}
	return m
}
