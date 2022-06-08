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
			name:         "prefix-one-word-with-max-three-results",
			inputKey:     []string{"the"},
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxResults(3)},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1},
					{Key: []string{"the", "quick", "brown", "fox"}, Value: 2},
					{Key: []string{"the", "quick", "swimmer"}, Value: 3},
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
			name:            "prefix-non-existing",
			inputKey:        []string{"non-existing"},
			expectedResults: &trie.SearchResults{},
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
					{Key: []string{"the"}, Value: 1, EditCount: 1},
					{Key: []string{"the", "green", "tree"}, Value: 4, EditCount: 1},
				},
			},
		},
		{
			name:         "edit-distance-one-edit-with-edit-opts",
			inputKey:     []string{"the", "tree"},
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(1), trie.WithEditOps()},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1, EditCount: 1, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeInsert, KeyPart: "tree"},
					}},
					{Key: []string{"the", "green", "tree"}, Value: 4, EditCount: 1, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeDelete, KeyPart: "green"},
						{Type: trie.EditOpTypeNoEdit, KeyPart: "tree"},
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
					{Key: []string{"the"}, Value: 1, EditCount: 1, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeInsert, KeyPart: "tree"},
					}},
					{Key: []string{"the", "quick", "swimmer"}, Value: 3, EditCount: 2, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeDelete, KeyPart: "quick"},
						{Type: trie.EditOpTypeReplace, KeyPart: "swimmer", ReplaceWith: "tree"},
					}},
					{Key: []string{"the", "green", "tree"}, Value: 4, EditCount: 1, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeDelete, KeyPart: "green"},
						{Type: trie.EditOpTypeNoEdit, KeyPart: "tree"},
					}},
					{Key: []string{"an", "apple", "tree"}, Value: 5, EditCount: 2, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeDelete, KeyPart: "an"},
						{Type: trie.EditOpTypeReplace, KeyPart: "apple", ReplaceWith: "the"},
						{Type: trie.EditOpTypeNoEdit, KeyPart: "tree"},
					}},
					{Key: []string{"an", "umbrella"}, Value: 6, EditCount: 2, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeReplace, KeyPart: "an", ReplaceWith: "the"},
						{Type: trie.EditOpTypeReplace, KeyPart: "umbrella", ReplaceWith: "tree"},
					}},
				},
			},
		},
		{
			name:         "edit-distance-two-edits-with-edit-opts-with-max-four-results",
			inputKey:     []string{"the", "tree"},
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(2), trie.WithEditOps(), trie.WithMaxResults(4)},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1, EditCount: 1, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeInsert, KeyPart: "tree"},
					}},
					{Key: []string{"the", "quick", "swimmer"}, Value: 3, EditCount: 2, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeDelete, KeyPart: "quick"},
						{Type: trie.EditOpTypeReplace, KeyPart: "swimmer", ReplaceWith: "tree"},
					}},
					{Key: []string{"the", "green", "tree"}, Value: 4, EditCount: 1, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeDelete, KeyPart: "green"},
						{Type: trie.EditOpTypeNoEdit, KeyPart: "tree"},
					}},
					{Key: []string{"an", "apple", "tree"}, Value: 5, EditCount: 2, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeDelete, KeyPart: "an"},
						{Type: trie.EditOpTypeReplace, KeyPart: "apple", ReplaceWith: "the"},
						{Type: trie.EditOpTypeNoEdit, KeyPart: "tree"},
					}},
				},
			},
		},
		{
			name:         "edit-distance-one-edit-with-topk",
			inputKey:     []string{"the", "tree"},
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(1), trie.WithTopKLeastEdited(1)},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1, EditCount: 1},
				},
			},
		},
		{
			name:         "edit-distance-two-edits-with-edit-opts-with-topk",
			inputKey:     []string{"the", "tree"},
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(2), trie.WithEditOps(), trie.WithTopKLeastEdited(4)},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1, EditCount: 1, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeInsert, KeyPart: "tree"},
					}},
					{Key: []string{"the", "green", "tree"}, Value: 4, EditCount: 1, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeDelete, KeyPart: "green"},
						{Type: trie.EditOpTypeNoEdit, KeyPart: "tree"},
					}},
					{Key: []string{"the", "quick", "swimmer"}, Value: 3, EditCount: 2, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeDelete, KeyPart: "quick"},
						{Type: trie.EditOpTypeReplace, KeyPart: "swimmer", ReplaceWith: "tree"},
					}},
					{Key: []string{"an", "apple", "tree"}, Value: 5, EditCount: 2, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeDelete, KeyPart: "an"},
						{Type: trie.EditOpTypeReplace, KeyPart: "apple", ReplaceWith: "the"},
						{Type: trie.EditOpTypeNoEdit, KeyPart: "tree"},
					}},
				},
			},
		},
		{
			name:         "edit-distance-two-edits-with-two-topk",
			inputKey:     []string{"the", "tree"},
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(2), trie.WithTopKLeastEdited(2)},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1, EditCount: 1},
					{Key: []string{"the", "green", "tree"}, Value: 4, EditCount: 1},
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

func TestTrie_Search_InvalidUsage_EditDistance_LessThanZeroDistance(t *testing.T) {
	assert.PanicsWithError(t, "invalid usage: maxDistance must be greater than zero", func() {
		trie.WithMaxEditDistance(0)
	})
	assert.PanicsWithError(t, "invalid usage: maxDistance must be greater than zero", func() {
		trie.WithMaxEditDistance(-1)
	})
}

func TestTrie_Search_InvalidUsage_MaxResults_LessThanZero(t *testing.T) {
	assert.PanicsWithError(t, "invalid usage: maxResults must be greater than zero", func() {
		trie.WithMaxResults(0)
	})
	assert.PanicsWithError(t, "invalid usage: maxResults must be greater than zero", func() {
		trie.WithMaxResults(-1)
	})
}

func TestTrie_Search_InvalidUsage_TopK_LessThanZeroK(t *testing.T) {
	assert.PanicsWithError(t, "invalid usage: k must be greater than zero", func() {
		trie.WithTopKLeastEdited(0)
	})
	assert.PanicsWithError(t, "invalid usage: k must be greater than zero", func() {
		trie.WithTopKLeastEdited(-1)
	})
}

func TestTrie_Search_InvalidUsage_EditOpsWithoutMaxEditDistance(t *testing.T) {
	tri := trie.New()

	assert.PanicsWithError(t, "invalid usage: WithEditOps() must not be passed without WithMaxEditDistance()", func() {
		tri.Search(nil, trie.WithEditOps())
	})
}

func TestTrie_Search_InvalidUsage_TopKWithoutMaxEditDistance(t *testing.T) {
	tri := trie.New()

	assert.PanicsWithError(t, "invalid usage: WithTopKLeastEdited() must not be passed without WithMaxEditDistance()", func() {
		tri.Search(nil, trie.WithTopKLeastEdited(1))
	})
}

func TestTrie_Search_InvalidUsage_ExactKeyWithMaxEditDistance(t *testing.T) {
	tri := trie.New()

	assert.PanicsWithError(t, "invalid usage: WithExactKey() must not be passed with WithMaxEditDistance()", func() {
		tri.Search(nil, trie.WithExactKey(), trie.WithMaxEditDistance(1))
	})
}

func TestTrie_Search_InvalidUsage_ExactKeyWithMaxResults(t *testing.T) {
	tri := trie.New()

	assert.PanicsWithError(t, "invalid usage: WithExactKey() must not be passed with WithMaxResults()", func() {
		tri.Search(nil, trie.WithExactKey(), trie.WithMaxResults(1))
	})
}

func TestTrie_Search_InvalidUsage_MaxResultsWithTopKLeastEdited(t *testing.T) {
	tri := trie.New()

	assert.PanicsWithError(t, "invalid usage: WithMaxResults() must not be passed with WithTopKLeastEdited()", func() {
		tri.Search(nil, trie.WithMaxResults(1), trie.WithMaxEditDistance(1), trie.WithTopKLeastEdited(1))
	})
}
