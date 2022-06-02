package trie_test

import (
	"testing"

	"github.com/shivamMg/trie"
	"github.com/stretchr/testify/assert"
)

func TestTrie_Search(t *testing.T) {
	tri := trie.New()
	tri.Put([]string{"the"}, 1)
	tri.Put([]string{"the", "quick", "brown", "fox"}, 2)
	tri.Put([]string{"the", "quick", "swimmer"}, 3)
	tri.Put([]string{"the", "green", "tree"}, 4)
	tri.Put([]string{"an", "apple", "tree"}, 5)
	tri.Put([]string{"an", "umbrella"}, 6)

	testCases := []struct {
		name            string
		inputKey        []string
		inputOptions    []func(*trie.SearchOptions)
		expectedResults *trie.SearchResults
	}{
		{
			name:     "prefix-one-word",
			inputKey: []string{"the"},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1},
					{Key: []string{"the", "quick", "brown", "fox"}, Value: 2},
					{Key: []string{"the", "quick", "swimmer"}, Value: 3},
					{Key: []string{"the", "green", "tree"}, Value: 4},
				},
			},
		},
		{
			name:     "prefix-multiple-words",
			inputKey: []string{"the", "quick"},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the", "quick", "brown", "fox"}, Value: 2},
					{Key: []string{"the", "quick", "swimmer"}, Value: 3},
				},
			},
		},
		{
			name:         "prefix-one-word-with-exact-key",
			inputKey:     []string{"the"},
			inputOptions: []func(*trie.SearchOptions){trie.WithExactKey()},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1},
				},
			},
		},
		{
			name:         "prefix-multiple-words-with-exact-key",
			inputKey:     []string{"the", "quick", "swimmer"},
			inputOptions: []func(*trie.SearchOptions){trie.WithExactKey()},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the", "quick", "swimmer"}, Value: 3},
				},
			},
		},
		{
			name:         "edit-distance-one-edit",
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
			name:         "edit-distance-one-edit-with-edit-opts",
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
			name:         "edit-distance-two-edits-with-edit-opts",
			inputKey:     []string{"the", "tree"},
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(2), trie.WithEditOps()},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNone, KeyPart: "the"},
						{Type: trie.EditOpTypeInsert, KeyPart: "tree"},
					}},
					{Key: []string{"the", "quick", "swimmer"}, Value: 3, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNone, KeyPart: "the"},
						{Type: trie.EditOpTypeDelete, KeyPart: "quick"},
						{Type: trie.EditOpTypeReplace, KeyPart: "swimmer", ReplaceWith: "tree"},
					}},
					{Key: []string{"the", "green", "tree"}, Value: 4, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNone, KeyPart: "the"},
						{Type: trie.EditOpTypeDelete, KeyPart: "green"},
						{Type: trie.EditOpTypeNone, KeyPart: "tree"},
					}},
					{Key: []string{"an", "apple", "tree"}, Value: 5, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeDelete, KeyPart: "an"},
						{Type: trie.EditOpTypeReplace, KeyPart: "apple", ReplaceWith: "the"},
						{Type: trie.EditOpTypeNone, KeyPart: "tree"},
					}},
					{Key: []string{"an", "umbrella"}, Value: 6, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeReplace, KeyPart: "an", ReplaceWith: "the"},
						{Type: trie.EditOpTypeReplace, KeyPart: "umbrella", ReplaceWith: "tree"},
					}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tri.Search(tc.inputKey, tc.inputOptions...)
			assert.Len(t, actual.Results, len(tc.expectedResults.Results))
			assert.Equal(t, tc.expectedResults, actual)
		})
	}
}

func TestTrie_Search_InvalidUsage_NegativeDistance(t *testing.T) {
	assert.PanicsWithError(t, "invalid usage: maxDistance must be greater than zero", func() {
		trie.WithMaxEditDistance(-1)
	})
}

func TestTrie_Search_InvalidUsage_EditOpsWithoutMaxEditDistance(t *testing.T) {
	tri := trie.New()

	assert.PanicsWithError(t, "invalid usage: WithEditOps() must not be passed without WithMaxEditDistance()", func() {
		tri.Search(nil, trie.WithEditOps())
	})
}

func TestTrie_Search_InvalidUsage_ExactKeyWithMaxEditDistance(t *testing.T) {
	tri := trie.New()

	assert.PanicsWithError(t, "invalid usage: WithExactKey() must not be passed with WithMaxEditDistance()", func() {
		tri.Search(nil, trie.WithExactKey(), trie.WithMaxEditDistance(1))
	})
}
