package trie_test

import (
	"testing"

	"github.com/shivamMg/ppds/tree"
	"github.com/shivamMg/trie"
	"github.com/stretchr/testify/assert"
)

func TestNode_SetValue(t *testing.T) {
	tri := trie.New()
	tri.Put([]string{"a", "b"}, 1)
	tri.Put([]string{"a", "b", "c"}, 2)

	node := tri.Root()
	node.SetValue(10)
	assert.Equal(t, nil, node.Value())

	node = tri.Root().ChildNodes()[0]
	node.SetValue(10)
	assert.Equal(t, nil, node.Value())

	node = tri.Root().ChildNodes()[0].ChildNodes()[0]
	node.SetValue(10)
	assert.Equal(t, 10, node.Value())
}

func TestTrie_Put(t *testing.T) {
	tri := trie.New()
	existed := tri.Put([]string{"an", "umbrella"}, 2)
	assert.False(t, existed)
	existed = tri.Put([]string{"the"}, 1)
	assert.False(t, existed)
	existed = tri.Put([]string{"the", "swimmer"}, 4)
	assert.False(t, existed)
	existed = tri.Put([]string{"the", "tree"}, 3)
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
	tri.Put([]string{"an", "apple", "tree"}, 5)
	tri.Put([]string{"an", "umbrella"}, 6)
	tri.Put([]string{"the"}, 1)
	tri.Put([]string{"the", "green", "tree"}, 4)
	tri.Put([]string{"the", "quick", "brown", "fox"}, 2)
	tri.Put([]string{"the", "quick", "swimmer"}, 3)

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
