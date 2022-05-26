package trie_test

import (
	"testing"

	"github.com/shivamMg/ppds/tree"
	"github.com/shivamMg/trie"
	"github.com/stretchr/testify/assert"
)

func TestTrie_Search_Put(t *testing.T) {
	tri := trie.New()
	tri.Put([]string{"the"}, 0)
	tri.Put([]string{"the", "quick", "brown", "fox"}, 1)
	tri.Put([]string{"the", "quick", "swimmer"}, 2)
	tri.Put([]string{"the", "green", "tree"}, 3)
	tri.Put([]string{"an", "apple", "on", "the", "tree"}, 4)
	tri.Put([]string{"an", "umbrella"}, 5)

	tree.PrintHrn(tri.Root())
	expected := `*
├─ an
│  ├─ apple
│  │  └─ on
│  │     └─ the
│  │        └─ tree
│  └─ umbrella
└─ the
   ├─ green
   │  └─ tree
   └─ quick
      ├─ brown
      │  └─ fox
      └─ swimmer
`

	actual := tree.SprintHrn(tri.Root())

	assert.Equal(t, expected, actual)
}
