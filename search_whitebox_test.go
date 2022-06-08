package trie

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrie_getEditOps(t *testing.T) {
	testCases := []struct {
		fromKeyColumn []string
		toKey         []string
		rows          [][]int
		expectedOps   []*EditOp
	}{
		{
			fromKeyColumn: strings.Split("sitting", ""),
			toKey:         strings.Split("kitten", ""),
			rows: [][]int{
				{0, 1, 2, 3, 4, 5, 6},
				{1, 1, 2, 3, 4, 5, 6},
				{2, 2, 1, 2, 3, 4, 5},
				{3, 3, 2, 1, 2, 3, 4},
				{4, 4, 3, 2, 1, 2, 3},
				{5, 5, 4, 3, 2, 2, 3},
				{6, 6, 5, 4, 3, 3, 2},
				{7, 7, 6, 5, 4, 4, 3},
			},
			expectedOps: []*EditOp{
				{Type: EditOpTypeReplace, KeyPart: "s", ReplaceWith: "k"},
				{Type: EditOpTypeNoEdit, KeyPart: "i"},
				{Type: EditOpTypeNoEdit, KeyPart: "t"},
				{Type: EditOpTypeNoEdit, KeyPart: "t"},
				{Type: EditOpTypeReplace, KeyPart: "i", ReplaceWith: "e"},
				{Type: EditOpTypeNoEdit, KeyPart: "n"},
				{Type: EditOpTypeDelete, KeyPart: "g"},
			},
		},
		{
			fromKeyColumn: strings.Split("Sunday", ""),
			toKey:         strings.Split("Saturday", ""),
			rows: [][]int{
				{0, 1, 2, 3, 4, 5, 6, 7, 8},
				{1, 0, 1, 2, 3, 4, 5, 6, 7},
				{2, 1, 1, 2, 2, 3, 4, 5, 6},
				{3, 2, 2, 2, 3, 3, 4, 5, 6},
				{4, 3, 3, 3, 3, 4, 3, 4, 5},
				{5, 4, 3, 4, 4, 4, 4, 3, 4},
				{6, 5, 4, 4, 5, 5, 5, 4, 3},
			},
			expectedOps: []*EditOp{
				{Type: EditOpTypeNoEdit, KeyPart: "S"},
				{Type: EditOpTypeInsert, KeyPart: "a"},
				{Type: EditOpTypeInsert, KeyPart: "t"},
				{Type: EditOpTypeNoEdit, KeyPart: "u"},
				{Type: EditOpTypeReplace, KeyPart: "n", ReplaceWith: "r"},
				{Type: EditOpTypeNoEdit, KeyPart: "d"},
				{Type: EditOpTypeNoEdit, KeyPart: "a"},
				{Type: EditOpTypeNoEdit, KeyPart: "y"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("from %v to %v", tc.fromKeyColumn, tc.toKey), func(t *testing.T) {
			tri := New()
			actual := tri.getEditOps(tc.rows, tc.fromKeyColumn, tc.toKey)
			assert.Equal(t, tc.expectedOps, actual)
		})
	}
}
