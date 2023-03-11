package trie_test

import (
	"bufio"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/shivamMg/trie"
	"github.com/stretchr/testify/assert"
)

var (
	benchmarkResults *trie.SearchResults // https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go
	wordsTrie        *trie.Trie
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
			name:     "prefix-empty",
			inputKey: []string{},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1},
					{Key: []string{"the", "quick", "brown", "fox"}, Value: 2},
					{Key: []string{"the", "quick", "swimmer"}, Value: 3},
					{Key: []string{"the", "green", "tree"}, Value: 4},
					{Key: []string{"an", "apple", "tree"}, Value: 5},
					{Key: []string{"an", "umbrella"}, Value: 6},
				},
			},
		},
		{
			name:     "prefix-nil",
			inputKey: nil,
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1},
					{Key: []string{"the", "quick", "brown", "fox"}, Value: 2},
					{Key: []string{"the", "quick", "swimmer"}, Value: 3},
					{Key: []string{"the", "green", "tree"}, Value: 4},
					{Key: []string{"an", "apple", "tree"}, Value: 5},
					{Key: []string{"an", "umbrella"}, Value: 6},
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
					{Key: []string{"the"}, Value: 1, EditDistance: 1},
					{Key: []string{"the", "green", "tree"}, Value: 4, EditDistance: 1},
				},
			},
		},
		{
			name:         "edit-distance-one-edit-with-edit-opts",
			inputKey:     []string{"the", "tree"},
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(1), trie.WithEditOps()},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1, EditDistance: 1, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeInsert, KeyPart: "tree"},
					}},
					{Key: []string{"the", "green", "tree"}, Value: 4, EditDistance: 1, EditOps: []*trie.EditOp{
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
					{Key: []string{"the"}, Value: 1, EditDistance: 1, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeInsert, KeyPart: "tree"},
					}},
					{Key: []string{"the", "quick", "swimmer"}, Value: 3, EditDistance: 2, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeDelete, KeyPart: "quick"},
						{Type: trie.EditOpTypeReplace, KeyPart: "swimmer", ReplaceWith: "tree"},
					}},
					{Key: []string{"the", "green", "tree"}, Value: 4, EditDistance: 1, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeDelete, KeyPart: "green"},
						{Type: trie.EditOpTypeNoEdit, KeyPart: "tree"},
					}},
					{Key: []string{"an", "apple", "tree"}, Value: 5, EditDistance: 2, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeDelete, KeyPart: "an"},
						{Type: trie.EditOpTypeReplace, KeyPart: "apple", ReplaceWith: "the"},
						{Type: trie.EditOpTypeNoEdit, KeyPart: "tree"},
					}},
					{Key: []string{"an", "umbrella"}, Value: 6, EditDistance: 2, EditOps: []*trie.EditOp{
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
					{Key: []string{"the"}, Value: 1, EditDistance: 1, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeInsert, KeyPart: "tree"},
					}},
					{Key: []string{"the", "quick", "swimmer"}, Value: 3, EditDistance: 2, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeDelete, KeyPart: "quick"},
						{Type: trie.EditOpTypeReplace, KeyPart: "swimmer", ReplaceWith: "tree"},
					}},
					{Key: []string{"the", "green", "tree"}, Value: 4, EditDistance: 1, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeDelete, KeyPart: "green"},
						{Type: trie.EditOpTypeNoEdit, KeyPart: "tree"},
					}},
					{Key: []string{"an", "apple", "tree"}, Value: 5, EditDistance: 2, EditOps: []*trie.EditOp{
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
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(1), trie.WithTopKLeastEdited(), trie.WithMaxResults(1)},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1, EditDistance: 1},
				},
			},
		},
		{
			name:         "edit-distance-two-edits-with-edit-opts-with-topk",
			inputKey:     []string{"the", "tree"},
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(2), trie.WithEditOps(), trie.WithTopKLeastEdited(), trie.WithMaxResults(4)},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1, EditDistance: 1, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeInsert, KeyPart: "tree"},
					}},
					{Key: []string{"the", "green", "tree"}, Value: 4, EditDistance: 1, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeDelete, KeyPart: "green"},
						{Type: trie.EditOpTypeNoEdit, KeyPart: "tree"},
					}},
					{Key: []string{"the", "quick", "swimmer"}, Value: 3, EditDistance: 2, EditOps: []*trie.EditOp{
						{Type: trie.EditOpTypeNoEdit, KeyPart: "the"},
						{Type: trie.EditOpTypeDelete, KeyPart: "quick"},
						{Type: trie.EditOpTypeReplace, KeyPart: "swimmer", ReplaceWith: "tree"},
					}},
					{Key: []string{"an", "apple", "tree"}, Value: 5, EditDistance: 2, EditOps: []*trie.EditOp{
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
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(2), trie.WithTopKLeastEdited(), trie.WithMaxResults(2)},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1, EditDistance: 1},
					{Key: []string{"the", "green", "tree"}, Value: 4, EditDistance: 1},
				},
			},
		},
		{
			name:         "edit-distance-empty",
			inputKey:     []string{},
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(2), trie.WithTopKLeastEdited(), trie.WithMaxResults(5)},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1, EditDistance: 1},
					{Key: []string{"an", "umbrella"}, Value: 6, EditDistance: 2},
				},
			},
		},
		{
			name:         "edit-distance-nil",
			inputKey:     nil,
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(2), trie.WithTopKLeastEdited(), trie.WithMaxResults(5)},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: []string{"the"}, Value: 1, EditDistance: 1},
					{Key: []string{"an", "umbrella"}, Value: 6, EditDistance: 2},
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

func TestTrie_Search_WordsTrie(t *testing.T) {
	tri := getWordsTrie()
	testCases := []struct {
		name            string
		inputKey        []string
		inputOptions    []func(*trie.SearchOptions)
		expectedResults *trie.SearchResults
	}{
		{
			name:     "prefix",
			inputKey: strings.Split("aband", ""),
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: strings.Split("abandon", "")},
					{Key: strings.Split("abandoned", "")},
					{Key: strings.Split("abandoning", "")},
					{Key: strings.Split("abandonment", "")},
				},
			},
		},
		{
			name:         "edit-distance",
			inputKey:     strings.Split("wheat", ""),
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(1)},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: strings.Split("wheat", ""), EditDistance: 0},
					{Key: strings.Split("wheal", ""), EditDistance: 1},
					{Key: strings.Split("whet", ""), EditDistance: 1},
					{Key: strings.Split("what", ""), EditDistance: 1},
					{Key: strings.Split("cheat", ""), EditDistance: 1},
					{Key: strings.Split("heat", ""), EditDistance: 1},
				},
			},
		},
		{
			name:         "edit-distance-with-max-results",
			inputKey:     strings.Split("national", ""),
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(3), trie.WithMaxResults(13)},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: strings.Split("nation", ""), EditDistance: 2},
					{Key: strings.Split("national", ""), EditDistance: 0},
					{Key: strings.Split("nationalism", ""), EditDistance: 3},
					{Key: strings.Split("nationalist", ""), EditDistance: 3},
					{Key: strings.Split("nationality", ""), EditDistance: 3},
					{Key: strings.Split("nationalize", ""), EditDistance: 3},
					{Key: strings.Split("nationally", ""), EditDistance: 2},
					{Key: strings.Split("natal", ""), EditDistance: 3},
					{Key: strings.Split("natural", ""), EditDistance: 3},
					{Key: strings.Split("nautical", ""), EditDistance: 3},
					{Key: strings.Split("notion", ""), EditDistance: 3},
					{Key: strings.Split("notional", ""), EditDistance: 1},
					{Key: strings.Split("notionally", ""), EditDistance: 3},
				},
			},
		},
		{
			name:         "edit-distance-with-topk",
			inputKey:     strings.Split("national", ""),
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(3), trie.WithMaxResults(13), trie.WithTopKLeastEdited()},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: strings.Split("national", ""), EditDistance: 0},
					{Key: strings.Split("notional", ""), EditDistance: 1},
					{Key: strings.Split("rational", ""), EditDistance: 1},
					{Key: strings.Split("nation", ""), EditDistance: 2},
					{Key: strings.Split("nationally", ""), EditDistance: 2},
					{Key: strings.Split("atonal", ""), EditDistance: 2},
					{Key: strings.Split("factional", ""), EditDistance: 2},
					{Key: strings.Split("optional", ""), EditDistance: 2},
					{Key: strings.Split("rationale", ""), EditDistance: 2},
					{Key: strings.Split("nationalism", ""), EditDistance: 3},
					{Key: strings.Split("nationalist", ""), EditDistance: 3},
					{Key: strings.Split("nationality", ""), EditDistance: 3},
					{Key: strings.Split("nationalize", ""), EditDistance: 3},
				},
			},
		},
		{
			name:         "edit-distance-with-topk-stop-after-prioritized",
			inputKey:     strings.Split("national", ""),
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(3), trie.WithMaxResults(2), trie.WithTopKLeastEdited()},
			expectedResults: &trie.SearchResults{
				Results: []*trie.SearchResult{
					{Key: strings.Split("national", ""), EditDistance: 0},
					{Key: strings.Split("notional", ""), EditDistance: 1},
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

func TestTrie_Search_InvalidUsage_EditOpsWithoutMaxEditDistance(t *testing.T) {
	tri := trie.New()

	assert.PanicsWithError(t, "invalid usage: WithEditOps() must not be passed without WithMaxEditDistance()", func() {
		tri.Search(nil, trie.WithEditOps())
	})
}

func TestTrie_Search_InvalidUsage_TopKWithoutMaxEditDistance(t *testing.T) {
	tri := trie.New()

	assert.PanicsWithError(t, "invalid usage: WithTopKLeastEdited() must not be passed without WithMaxEditDistance()", func() {
		tri.Search(nil, trie.WithTopKLeastEdited())
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

func TestTrie_Search_InvalidUsage_TopKLeastEditedWithoutMaxResults(t *testing.T) {
	tri := trie.New()

	assert.PanicsWithError(t, "invalid usage: WithTopKLeastEdited() must not be passed without WithMaxResults()", func() {
		tri.Search(nil, trie.WithMaxEditDistance(1), trie.WithTopKLeastEdited())
	})
}

func BenchmarkTrie_Search_WordsTrie(b *testing.B) {
	tri := getWordsTrie()
	benchmarks := []struct {
		name         string
		inputKey     []string
		inputOptions []func(*trie.SearchOptions)
	}{
		{
			name:     "prefix",
			inputKey: strings.Split("ab", ""),
		},
		{
			name:         "prefix-with-max-results",
			inputKey:     strings.Split("ab", ""),
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxResults(20)},
		},
		{
			name:         "edit-distance",
			inputKey:     strings.Split("someday", ""),
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(5)},
		},
		{
			name:         "edit-distance-with-edit-ops",
			inputKey:     strings.Split("someday", ""),
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(5), trie.WithEditOps()},
		},
		{
			name:         "edit-distance-with-edit-ops-with-max-results",
			inputKey:     strings.Split("someday", ""),
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(5), trie.WithEditOps(), trie.WithMaxResults(20)},
		},
		{
			name:     "edit-distance-with-edit-ops-with-max-results-with-top-k",
			inputKey: strings.Split("someday", ""),
			inputOptions: []func(*trie.SearchOptions){trie.WithMaxEditDistance(5), trie.WithEditOps(), trie.WithMaxResults(20),
				trie.WithTopKLeastEdited()},
		},
	}
	var results *trie.SearchResults
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				results = tri.Search(bm.inputKey, bm.inputOptions...)
			}
			benchmarkResults = results
		})
	}
}

func getWordsTrie() *trie.Trie {
	if wordsTrie != nil {
		return wordsTrie
	}
	f, err := os.Open("./demo/wasm/words.txt")
	if err != nil {
		panic(err)
	}
	tri := trie.New()
	r := bufio.NewReader(f)
	for {
		word, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		word = strings.TrimRight(word, "\n")
		word = strings.TrimRight(word, "\r") // windows
		key := strings.Split(word, "")
		tri.Put(key, nil)
	}
	wordsTrie = tri
	return tri
}
