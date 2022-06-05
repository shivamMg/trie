## trie

This library provides an implementation of the [Trie](https://en.wikipedia.org/wiki/Trie) data structure. It provides more features than the usual Trie prefix-search, and is meant to be used for use cases that require auto-completion.

### Example

```go
func Example() {
	tri := trie.New()
	tri.Put([]string{"the"}, 1)
	tri.Put([]string{"the", "quick", "brown", "fox"}, 2)
	tri.Put([]string{"the", "quick", "sports", "car"}, 3)
	tri.Put([]string{"the", "green", "tree"}, 4)
	tri.Put([]string{"an", "apple", "tree"}, 5)
	tri.Put([]string{"an", "umbrella"}, 6)

	tri.Root().Print()
	// Output (full trie with terminals ($)):
	// ^
	// ├─ the ($)
	// │  ├─ quick
	// │  │  ├─ brown
	// │  │  │  └─ fox ($)
	// │  │  └─ sports
	// │  │     └─ car ($)
	// │  └─ green
	// │     └─ tree ($)
	// └─ an
	//    ├─ apple
	//    │  └─ tree ($)
	//    └─ umbrella ($)

	results := tri.Search([]string{"the", "quick"})
	for _, res := range results.Results {
		fmt.Println(res.Key, res.Value)
	}
	// Output (prefix search):
	// [the quick brown fox] 2
	// [the quick sports car] 3

	key := []string{"the", "tree"}
	results = tri.Search(key, trie.WithMaxEditDistance(2), trie.WithEditOps()) // Edit can be insert, delete, replace
	for _, res := range results.Results {
		fmt.Println(res.Key, res.Value, res.EditCount) // EditCount is number of edits needed to convert to [the tree]
	}
	// Output (results not more than 2 edits away from [the tree]):
	// [the] 1 1
	// [the green tree] 4 1
	// [an apple tree] 5 2
	// [an umbrella] 6 2

    result := results.Results[2]
	fmt.Printf("To convert %v to %v:\n", result.Key, key)
	printEditOps(result.EditOps)
	// Output (edit operations needed to covert a result to [the tree]):
	// To convert [an apple tree] to [the tree]:
	// - delete "an"
	// - replace "apple" with "the"
	// - don't edit "tree"

	results = tri.Search(key, trie.WithMaxEditDistance(2), trie.WithTopKLeastEdited(2))
	for _, res := range results.Results {
		fmt.Println(res.Key, res.Value, res.EditCount)
	}
	// Output (top 2 least edited results):
	// [the] 1 1
	// [the green tree] 4 1
}
```

### References

* https://en.wikipedia.org/wiki/Levenshtein_distance#Iterative_with_full_matrix
* http://stevehanov.ca/blog/?id=114
* https://gist.github.com/jlherren/d97839b1276b9bd7faa930f74711a4b6