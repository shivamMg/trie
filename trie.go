package trie

import (
	"sort"

	"github.com/shivamMg/ppds/tree"
)

const RootKeyPart = "*"

type EditOpType int

const (
	EditOpTypeNone EditOpType = iota
	EditOpTypeInsert
	EditOpTypeDelete
	EditOpTypeReplace
)

type EditOp struct {
	Type         EditOpType
	KeyPart      string
	ReplacedWith string
}

type Trie struct {
	root *Node
}

type Node struct {
	keyPart    string
	isTerminal bool
	value      any
	children   map[string]*Node
}

type SearchResults struct {
	Results []*SearchResult
}

type SearchResult struct {
	Key     []string
	Value   any
	EditOps []*EditOp
}

type SearchOptions struct {
	exactKey        bool
	editDistance    bool
	maxEditDistance int
	editOps         bool
}

func WithExactKey() func(*SearchOptions) {
	return func(so *SearchOptions) {
		so.exactKey = true
	}
}

func WithMaxEditDistance(maxDistance int) func(*SearchOptions) {
	if maxDistance < 0 {
		panic("maxDistance cannot be negative")
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

func (n *Node) KeyPart() string {
	return n.keyPart
}

func (n *Node) IsTerminal() bool {
	return n.isTerminal
}

func (n *Node) Value() any {
	return n.value
}

func (n *Node) ChildNodes() []*Node {
	return n.childNodes()
}

func (n *Node) Data() interface{} {
	data := n.keyPart
	if n.isTerminal {
		data += " ($)"
	}
	return data
}

func (n *Node) Children() []tree.Node {
	children := n.childNodes()
	result := make([]tree.Node, len(children))
	for i, child := range children {
		result[i] = tree.Node(child)
	}
	return result
}

func (n *Node) childNodes() []*Node {
	children := make([]*Node, 0, len(n.children))
	for _, child := range n.children {
		children = append(children, child)
	}
	sort.Slice(children, func(i, j int) bool {
		return children[i].keyPart < children[j].keyPart
	})
	return children
}

func New() *Trie {
	return &Trie{root: &Node{keyPart: RootKeyPart}}
}

func (t *Trie) Root() *Node {
	return t.root
}

func (t *Trie) Put(key []string, value any) (existed bool) {
	node := t.root
	for i, part := range key {
		if node.children == nil {
			node.children = make(map[string]*Node)
		}
		child, ok := node.children[part]
		if !ok {
			child = &Node{keyPart: part}
			node.children[part] = child
		}
		if i == len(key)-1 {
			existed = child.isTerminal
			child.isTerminal = true
			child.value = value
		}
		node = child
	}
	return existed
}

// WithExactKey()
// WithMaxEditDistance(maxDistance int)
// WithEditOps()
// WithEqual(func(a, b string) bool)
func (t *Trie) Search(key []string, options ...func(*SearchOptions)) *SearchResults {
	opts := &SearchOptions{}
	for _, f := range options {
		f(opts)
	}

	if opts.maxEditDistance == 0 {
		return t.prefixSearch(key)
	}
	return t.levenshteinSearch(key, opts)
}

func (t *Trie) levenshteinSearch(prefixTokens []string, opts *SearchOptions) *SearchResults {
	curRow := make([]int, len(prefixTokens)+1)
	for i := 0; i < len(prefixTokens)+1; i++ {
		curRow[i] = i
	}
	results := &SearchResults{}
	for token, node := range t.root.children {
		t.levenshteinBuild(node, []string{token}, curRow, prefixTokens, results, opts)
	}
	return results
}

func (t *Trie) levenshteinBuild(node *Node, prevTokens []string, prevRow []int, prefixTokens []string, results *SearchResults, opts *SearchOptions) {
	columns := len(prefixTokens) + 1
	curRow := make([]int, len(prefixTokens)+1)
	curRow[0] = prevRow[0] + 1

	for i := 1; i < columns; i++ {
		replaceCost := 0
		if prefixTokens[i-1] != prevTokens[len(prevTokens)-1] {
			replaceCost = 1
		}
		minCost := min(
			curRow[i-1]+1,            // insertion
			prevRow[i]+1,             // deletion
			prevRow[i-1]+replaceCost, // replace
		)
		curRow[i] = minCost
		if replaceCost == 0 && minCost == prevRow[i-1] {
			// edited=true
		}
	}

	if curRow[columns-1] <= opts.maxEditDistance && node.isTerminal {
		result := &SearchResult{Value: node.value}
		results.Results = append(results.Results, result)
	}
}

func min(a, b, c int) int {
	if b < a {
		a = b
	}
	if c < a {
		a = c
	}
	return a
}

func (t *Trie) prefixSearch(prefixTokens []string) *SearchResults {
	results := &SearchResults{}
	node := t.root
	for _, token := range prefixTokens {
		if node.children == nil {
			return results
		}
		child, ok := node.children[token]
		if !ok {
			return results
		}
		node = child
	}
	t.populate(results, node, prefixTokens)
	return results
}

func (t *Trie) populate(results *SearchResults, node *Node, prefixTokens []string) {
	if node.isTerminal {
		// tokenResults := make([]*TokenResult, len(prefixTokens))
		// for i, token := range prefixTokens {
		// 	tokenResults[i].Token = token
		// }
		// result := &SearchResult{Value: node.value, TokenResults: tokenResults}
		// results.Results = append(results.Results, result)
	}
	for token, child := range node.children {
		prefixTokens = append(prefixTokens, token)
		t.populate(results, child, prefixTokens)
		prefixTokens = prefixTokens[:len(prefixTokens)-1]
	}
}
