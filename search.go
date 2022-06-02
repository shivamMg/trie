package trie

import (
	"errors"
	"math"
)

type EditOpType int

const (
	EditOpTypeNone EditOpType = iota
	EditOpTypeInsert
	EditOpTypeDelete
	EditOpTypeReplace
)

type EditOp struct {
	Type        EditOpType
	KeyPart     string
	ReplaceWith string
}

type SearchResults struct {
	Results []*SearchResult
}

type SearchResult struct {
	Key       []string
	Value     interface{}
	EditCount int
	EditOps   []*EditOp
}

type SearchOptions struct {
	exactKey        bool
	editDistance    bool
	maxEditDistance int
	editOps         bool
	topKLeastEdited bool
}

func WithExactKey() func(*SearchOptions) {
	return func(so *SearchOptions) {
		so.exactKey = true
	}
}

func WithMaxEditDistance(maxDistance int) func(*SearchOptions) {
	if maxDistance <= 0 {
		panic(errors.New("invalid usage: maxDistance must be greater than zero"))
	}
	return func(so *SearchOptions) {
		so.editDistance = true
		so.maxEditDistance = maxDistance
	}
}

func WithEditOps() func(*SearchOptions) {
	return func(so *SearchOptions) {
		so.editOps = true
	}
}

func WithTopKLeastEdited() func(*SearchOptions) {
	return func(so *SearchOptions) {
		so.topKLeastEdited = true
	}
}

func (t *Trie) Search(key []string, options ...func(*SearchOptions)) *SearchResults {
	opts := &SearchOptions{}
	for _, f := range options {
		f(opts)
	}
	if opts.editOps && !opts.editDistance {
		panic(errors.New("invalid usage: WithEditOps() must not be passed without WithMaxEditDistance()"))
	}
	if opts.exactKey && opts.editDistance {
		panic(errors.New("invalid usage: WithExactKey() must not be passed with WithMaxEditDistance()"))
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
	rows := [][]int{newRow}
	results := &SearchResults{}
	dllNode := t.root.childrenDLL.head
	for dllNode != nil {
		node := dllNode.trieNode
		t.buildWithEditDistance(results, node, []string{node.keyPart}, rows, key, opts)
		dllNode = dllNode.next
	}
	return results
}

func (t *Trie) buildWithEditDistance(results *SearchResults, node *Node, keyColumn []string, rows [][]int, key []string, opts *SearchOptions) {
	prevRow := rows[len(rows)-1]
	columns := len(key) + 1
	newRow := make([]int, columns)
	newRow[0] = prevRow[0] + 1
	for i := 1; i < columns; i++ {
		replaceCost := 1
		if key[i-1] == keyColumn[len(keyColumn)-1] {
			replaceCost = 0
		}
		newRow[i] = min(
			newRow[i-1]+1,            // insertion
			prevRow[i]+1,             // deletion
			prevRow[i-1]+replaceCost, // substitution
		)
	}
	rows = append(rows, newRow)

	if newRow[columns-1] <= opts.maxEditDistance && node.isTerminal {
		result := &SearchResult{Key: keyColumn, Value: node.value, EditCount: newRow[columns-1]}
		if opts.editOps {
			result.EditOps = t.getEditOps(rows, keyColumn, key)
		}
		results.Results = append(results.Results, result)
	}

	if min(newRow...) <= opts.maxEditDistance {
		dllNode := node.childrenDLL.head
		for dllNode != nil {
			child := dllNode.trieNode
			t.buildWithEditDistance(results, child, append(keyColumn, child.keyPart), rows, key, opts)
			dllNode = dllNode.next
		}
	}
}

func (t *Trie) getEditOps(rows [][]int, keyColumn []string, key []string) []*EditOp {
	// https://gist.github.com/jlherren/d97839b1276b9bd7faa930f74711a4b6
	ops := make([]*EditOp, 0, len(key))
	r, c := len(rows)-1, len(rows[0])-1
	for r > 0 || c > 0 {
		insertionCost, deletionCost, substitutionCost := math.MaxInt, math.MaxInt, math.MaxInt
		if c > 0 {
			insertionCost = rows[r][c-1]
		}
		if r > 0 {
			deletionCost = rows[r-1][c]
		}
		if r > 0 && c > 0 {
			substitutionCost = rows[r-1][c-1]
		}
		minCost := min(insertionCost, deletionCost, substitutionCost)
		if minCost == substitutionCost {
			if rows[r][c] > rows[r-1][c-1] {
				ops = append(ops, &EditOp{Type: EditOpTypeReplace, KeyPart: keyColumn[r-1], ReplaceWith: key[c-1]})
			} else {
				ops = append(ops, &EditOp{Type: EditOpTypeNone, KeyPart: keyColumn[r-1]})
			}
			r -= 1
			c -= 1
		} else if minCost == deletionCost {
			ops = append(ops, &EditOp{Type: EditOpTypeDelete, KeyPart: keyColumn[r-1]})
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
	t.build(results, node, prefixKey)
	return results
}

func (t *Trie) build(results *SearchResults, node *Node, prefixKey []string) {
	if node.isTerminal {
		result := &SearchResult{Key: prefixKey, Value: node.value}
		results.Results = append(results.Results, result)
	}
	dllNode := node.childrenDLL.head
	for dllNode != nil {
		child := dllNode.trieNode
		t.build(results, child, append(prefixKey, child.keyPart))
		dllNode = dllNode.next
	}
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
