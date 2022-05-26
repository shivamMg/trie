package trie_test

import (
	"strings"
	"testing"

	"github.com/shivamMg/ppds/tree"
	"github.com/shivamMg/trie"
	"github.com/stretchr/testify/assert"
)

func TestTrie_Put(t *testing.T) {
	tri := trie.New()
	tri.Put([]string{"the"}, 1)
	tri.Put([]string{"an", "umbrella"}, 2)
	tri.Put([]string{"the", "tree"}, 3)
	tri.Put([]string{"the", "swimmer"}, 4)

	expected := `*
├─ an
│  └─ umbrella ($)
└─ the ($)
   ├─ swimmer ($)
   └─ tree ($)
`

	actual := tree.SprintHrn(tri.Root())
	assert.Equal(t, expected, actual)

	root := tri.Root()
	assert.Equal(t, root.KeyPart(), trie.RootKeyPart)
	assert.False(t, root.IsTerminal())
	assert.Nil(t, root.Value())

	rootChildren := root.ChildNodes()
	an, the := rootChildren[0], rootChildren[1]
	assert.Equal(t, an.KeyPart(), "an")
	assert.False(t, an.IsTerminal())
	assert.Nil(t, an.Value())

	assert.Equal(t, the.KeyPart(), "the")
	assert.True(t, the.IsTerminal())
	assert.Equal(t, the.Value(), 1)

	umbrella := an.ChildNodes()[0]
	assert.Equal(t, umbrella.KeyPart(), "umbrella")
	assert.True(t, umbrella.IsTerminal())
	assert.Equal(t, umbrella.Value(), 2)

	theChildren := the.ChildNodes()
	swimmer, tree_ := theChildren[0], theChildren[1]
	assert.Equal(t, swimmer.KeyPart(), "swimmer")
	assert.True(t, swimmer.IsTerminal())
	assert.Equal(t, swimmer.Value(), 4)

	assert.Equal(t, tree_.KeyPart(), "tree")
	assert.True(t, tree_.IsTerminal())
	assert.Equal(t, tree_.Value(), 3)
}

func TestTrie_Search(t *testing.T) {
	tri := trie.New()
	tri.Put([]string{"the"}, 1)
	tri.Put([]string{"the", "quick", "brown", "fox"}, 2)
	tri.Put([]string{"the", "quick", "swimmer"}, 3)
	tri.Put([]string{"the", "green", "tree"}, 4)
	tri.Put([]string{"an", "apple", "tree"}, 5)
	tri.Put([]string{"an", "umbrella"}, 6)

	testCases := []struct {
		inputKey        []string
		inputOptions    []func(*trie.SearchOptions)
		expectedResults *trie.SearchResults
	}{
		{
			inputKey: []string{"the"},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1},
					{Key: []string{"the", "green", "tree"}, Value: 4},
					{Key: []string{"the", "quick", "brown", "fox"}, Value: 2},
					{Key: []string{"the", "quick", "swimmer"}, Value: 3},
				},
			},
		},
		{
			inputKey: []string{"the", "quick"},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the", "quick", "brown", "fox"}, Value: 2},
					{Key: []string{"the", "quick", "swimmer"}, Value: 3},
				},
			},
		},
		{
			inputKey:     []string{"the"},
			inputOptions: []func(*trie.SearchOptions){trie.WithExactKey()},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1},
				},
			},
		},
		{
			inputKey:     []string{"the", "quick", "swimmer"},
			inputOptions: []func(*trie.SearchOptions){trie.WithExactKey()},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the", "quick", "swimmer"}, Value: 3},
				},
			},
		},
		{
			inputKey:     []string{"the", "tree"},
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(1)},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1},
					{Key: []string{"the", "green", "tree"}, Value: 4},
				},
			},
		},
		{
			inputKey:     []string{"the", "tree"},
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(1), trie.WithEditOps()},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNone, KeyPart: "the"},
						{Type: trie.EditOpTypeInsert, KeyPart: "tree"},
					}},
					{Key: []string{"the", "green", "tree"}, Value: 4, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNone, KeyPart: "the"},
						{Type: trie.EditOpTypeDelete, KeyPart: "green"},
						{Type: trie.EditOpTypeNone, KeyPart: "tree"},
					}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(strings.Join(tc.inputKey, "-"), func(t *testing.T) {
			actual := tri.Search(tc.inputKey, tc.inputOptions...)
			assert.Equal(t, tc.expectedResults, actual)
		})
	}
}
