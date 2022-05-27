package trie_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/shivamMg/ppds/tree"
	"github.com/shivamMg/trie"
	"github.com/stretchr/testify/assert"
)

func TestTrie_Put(t *testing.T) {
	tri := trie.New()
	existed := tri.Put([]string{"the"}, 1)
	assert.False(t, existed)
	existed = tri.Put([]string{"an", "umbrella"}, 2)
	assert.False(t, existed)
	existed = tri.Put([]string{"the", "tree"}, 3)
	assert.False(t, existed)
	existed = tri.Put([]string{"the", "swimmer"}, 4)
	assert.False(t, existed)

	// validate full tree
	expected := `^
├─ an
│  └─ umbrella ($)
└─ the ($)
   ├─ swimmer ($)
   └─ tree ($)
`
	actual := tree.SprintHrn(tri.Root())
	assert.Equal(t, expected, actual)

	// validate attributes for each node
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

	// validate update
	existed = tri.Put([]string{"an", "umbrella"}, 5)
	assert.True(t, existed)
	assert.Equal(t, umbrella.Value(), 5)
}

func TestTrie_Delete(t *testing.T) {
	tri := trie.New()
	tri.Put([]string{"the"}, 1)
	tri.Put([]string{"the", "quick", "brown", "fox"}, 2)
	tri.Put([]string{"the", "quick", "swimmer"}, 3)
	tri.Put([]string{"the", "green", "tree"}, 4)
	tri.Put([]string{"an", "apple", "tree"}, 5)
	tri.Put([]string{"an", "umbrella"}, 6)

	value, existed := tri.Delete([]string{"the", "quick", "brown", "fox"})
	assert.True(t, existed)
	assert.Equal(t, value, 2)
	expected := `^
├─ an
│  ├─ apple
│  │  └─ tree ($)
│  └─ umbrella ($)
└─ the ($)
   ├─ green
   │  └─ tree ($)
   └─ quick
      └─ swimmer ($)
`
	assert.Equal(t, expected, tri.Root().Sprint())

	value, existed = tri.Delete([]string{"the", "quick", "swimmer"})
	assert.True(t, existed)
	assert.Equal(t, value, 3)
	expected = `^
├─ an
│  ├─ apple
│  │  └─ tree ($)
│  └─ umbrella ($)
└─ the ($)
   └─ green
      └─ tree ($)
`
	assert.Equal(t, expected, tri.Root().Sprint())

	value, existed = tri.Delete([]string{"the"})
	assert.True(t, existed)
	assert.Equal(t, value, 1)
	expected = `^
├─ an
│  ├─ apple
│  │  └─ tree ($)
│  └─ umbrella ($)
└─ the
   └─ green
      └─ tree ($)
`
	assert.Equal(t, expected, tri.Root().Sprint())

	value, existed = tri.Delete([]string{"non", "existing"})
	assert.False(t, existed)
	assert.Nil(t, value)
	expected = `^
├─ an
│  ├─ apple
│  │  └─ tree ($)
│  └─ umbrella ($)
└─ the
   └─ green
      └─ tree ($)
`
	assert.Equal(t, expected, tri.Root().Sprint())

	value, existed = tri.Delete([]string{"an"})
	assert.False(t, existed)
	assert.Nil(t, value)
	expected = `^
├─ an
│  ├─ apple
│  │  └─ tree ($)
│  └─ umbrella ($)
└─ the
   └─ green
      └─ tree ($)
`
	assert.Equal(t, expected, tri.Root().Sprint())
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
		{
			inputKey:     []string{"the", "tree"},
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(2), trie.WithEditOps()},
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
					{Key: []string{"the", "quick", "swimmer"}, Value: 3, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNone, KeyPart: "the"},
						{Type: trie.EditOpTypeDelete, KeyPart: "quick"},
						{Type: trie.EditOpTypeReplace, KeyPart: "swimmer", ReplacedWith: "tree"},
					}},
					{Key: []string{"an", "apple", "tree"}, Value: 5, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeReplace, KeyPart: "an", ReplacedWith: "the"},
						{Type: trie.EditOpTypeDelete, KeyPart: "apple"},
						{Type: trie.EditOpTypeNone, KeyPart: "tree"},
					}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(strings.Join(tc.inputKey, "-"), func(t *testing.T) {
			actual := tri.Search(tc.inputKey, tc.inputOptions...)
			assert.Len(t, actual.Results, len(tc.expectedResults.Results))
			// sort before comparing
			sortResults(tc.expectedResults)
			sortResults(actual)
			assert.Equal(t, tc.expectedResults, actual)
		})
	}
}

func TestTrie_Search_InvalidUsage_NegativeDistance(t *testing.T) {
	assert.PanicsWithError(t, "invalid usage: maxDistance must be > 0", func() {
		trie.WithMaxEditDistance(-1)
	})
}

func TestTrie_Search_InvalidUsage_EditOpsWithoutMaxEditDistance(t *testing.T) {
	tri := trie.New()

	assert.PanicsWithError(t, "invalid usage: WithEditOps() must be passed with WithMaxEditDistance()", func() {
		tri.Search(nil, trie.WithEditOps())
	})
}

func TestTrie_Search_InvalidUsage_ExactKeyWithMaxEditDistance(t *testing.T) {
	tri := trie.New()

	assert.PanicsWithError(t, "invalid usage: WithExactKey() cannot be passed with WithMaxEditDistance()", func() {
		tri.Search(nil, trie.WithExactKey(), trie.WithMaxEditDistance(1))
	})
}

func sortResults(results *trie.SearchResults) {
	sort.Slice(results.Results, func(i, j int) bool {
		return strings.Join(results.Results[i].Key, "") < strings.Join(results.Results[j].Key, "")
	})
}
