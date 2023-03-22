package trie_test

import (
	"errors"
	"testing"

	"github.com/shivamMg/trie"
	"github.com/stretchr/testify/assert"
)

func TestTrie_WalkErr(t *testing.T) {
	tri := trie.New()
	tri.Put([]string{"d", "a", "l", "i"}, 1)
	tri.Put([]string{"d", "a", "l", "i", "b"}, 2)
	tri.Put([]string{"d", "a", "l", "i", "b", "e"}, 3)
	tri.Put([]string{"d", "a", "l", "i", "b", "e", "r", "t"}, 4)

	var selected []string
	walker := func(key []string, node *trie.Node) error {
		what := node.Value().(int)
		if what == 3 {
			selected = key
			return errors.New("found")
		}
		return nil
	}

	err := tri.Walk(nil, walker)
	assert.EqualError(t, err, "found")
	assert.EqualValues(t, []string{"d", "a", "l", "i", "b", "e"}, selected)
}

func TestTrie_Walk(t *testing.T) {
	tri := trie.New()
	tri.Put([]string{"d", "a", "l", "i"}, []int{0, 1, 2, 4, 5})
	tri.Put([]string{"d", "a", "l", "i", "b"}, []int{1, 2, 4, 5})
	tri.Put([]string{"d", "a", "l", "i", "b", "e"}, []int{1, 0, 2, 4, 5, 0})
	tri.Put([]string{"d", "a", "l", "i", "b", "e", "r", "t"}, []int{1, 2, 4, 5})
	type KVPair struct {
		key   []string
		value []int
	}
	var selected []KVPair
	walker := func(key []string, node *trie.Node) error {
		what := node.Value().([]int)
		for _, i := range what {
			if i == 0 {
				selected = append(selected, KVPair{key, what})
				break
			}
		}
		return nil
	}

	err := tri.Walk(nil, walker)
	assert.NoError(t, err)
	expected := []KVPair{
		{[]string{"d", "a", "l", "i"}, []int{0, 1, 2, 4, 5}},
		{[]string{"d", "a", "l", "i", "b", "e"}, []int{1, 0, 2, 4, 5, 0}},
	}
	assert.EqualValues(t, expected, selected)

	selected = nil
	err = tri.Walk([]string{"d", "a", "l", "i", "b"}, walker)
	assert.NoError(t, err)
	expected = []KVPair{
		{[]string{"d", "a", "l", "i", "b", "e"}, []int{1, 0, 2, 4, 5, 0}},
	}
	assert.EqualValues(t, expected, selected)

	selected = nil
	err = tri.Walk([]string{"a", "b", "c"}, walker)
	assert.NoError(t, err)
	assert.Nil(t, selected)
}
