package trie_test

import (
	"github.com/shivamMg/trie"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrie_Filter(t *testing.T) {
	tri := trie.New()
	tri.Put([]string{"z"}, 2)
	tri.Put([]string{"d", "a", "l", "i"}, 2)
	tri.Put([]string{"d", "a", "l", "i", "b"}, 1)
	tri.Put([]string{"d", "a", "l", "i", "b", "e"}, 2)
	tri.Put([]string{"d", "a", "l", "i", "b", "e", "r", "t"}, 1)
	tri.Put([]string{"a", "a"}, 2)

	rs := tri.Filter(func(val interface{}) bool {
		what := val.(int)
		if what == 2 {
			return true
		}
		return false
	})

	for _, res := range rs.Results {
		assert.Equal(t, 2, res.Value)
	}
	assert.Equal(t, 4, len(rs.Results))

	dOrders := []struct {
		key []string
	}{
		{
			key: []string{"z"},
		},
		{
			key: []string{"d", "a", "l", "i"},
		},
		{
			[]string{"d", "a", "l", "i", "b", "e"},
		},
		{
			[]string{"a", "a"},
		},
	}

	for i, ord := range dOrders {
		assert.Equal(t, ord.key, rs.Results[i].Key)
	}

}

func TestTrie_FilterComplex(t *testing.T) {
	tri := trie.New()
	tri.Put([]string{"d", "a", "l", "i"}, []int{0, 1, 2, 4, 5})
	tri.Put([]string{"d", "a", "l", "i", "b"}, []int{1, 2, 4, 5})
	tri.Put([]string{"d", "a", "l", "i", "b", "e"}, []int{0, 1, 2, 4, 5})
	tri.Put([]string{"d", "a", "l", "i", "b", "e", "r", "t"}, []int{1, 2, 4, 5})

	rs := tri.Filter(func(val interface{}) bool {
		what := val.([]int)
		for _, i := range what {
			if i == 0 {
				return true
			}
		}
		return false
	})

	for _, res := range rs.Results {
		what := res.Value.([]int)
		ok := false
		for _, i := range what {
			if i == 0 {
				ok = true
			}
		}
		assert.True(t, ok)
	}
	assert.Equal(t, 2, len(rs.Results))

}
