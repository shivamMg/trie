## trie [![godoc](https://godoc.org/github.com/shivammg/trie?status.svg)](https://godoc.org/github.com/shivamMg/trie) ![Build](https://github.com/shivamMg/trie/actions/workflows/ci.yml/badge.svg?branch=master) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

An implementation of the [Trie](https://en.wikipedia.org/wiki/Trie) data structure in Go. It provides more features than the usual Trie prefix-search, and is meant to be used for auto-completion.

### Demo

Auto-completion demo can be tried at [shivamMg.github.io/trie](https://shivammg.github.io/trie/).

### Features

- Keys are `[]string` instead of `string`, thereby supporting more use cases - e.g. []string{the quick brown fox} can be a key where each word will be a node in the Trie
- Support for Put key and Delete key
- Support for Prefix search - e.g. searching for _nation_ might return _nation_, _national_, _nationalism_, _nationalist_, etc.
- Support for Edit distance search (aka Levenshtein distance) - e.g. searching for _wheat_ might return similar looking words like _wheat_, _cheat_, _heat_, _what_, etc.
- Order of search results is deterministic. It follows insertion (Put()) order.

### Examples

```go
tri := trie.New()
// Put keys ([]string) and values (any)
tri.Put([]string{"the"}, 1)
tri.Put([]string{"the", "quick", "brown", "fox"}, 2)
tri.Put([]string{"the", "quick", "sports", "car"}, 3)
tri.Put([]string{"the", "green", "tree"}, 4)
tri.Put([]string{"an", "apple", "tree"}, 5)
tri.Put([]string{"an", "umbrella"}, 6)

tri.Root().Print()
// Output (full trie with terminals ending with ($)):
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
// Output (prefix-based search):
// [the quick brown fox] 2
// [the quick sports car] 3

key := []string{"the", "tree"}
results = tri.Search(key, trie.WithMaxEditDistance(2), // An edit can be insert, delete, replace
	trie.WithEditOps())
for _, res := range results.Results {
	fmt.Println(res.Key, res.EditDistance) // EditDistance is number of edits needed to convert to [the tree]
}
// Output (results not more than 2 edits away from [the tree]):
// [the] 1
// [the green tree] 1
// [an apple tree] 2
// [an umbrella] 2

result := results.Results[2]
fmt.Printf("To convert %v to %v:\n", result.Key, key)
printEditOps(result.EditOps)
// Output (edit operations needed to covert a result to [the tree]):
// To convert [an apple tree] to [the tree]:
// - delete "an"
// - replace "apple" with "the"
// - don't edit "tree"

results = tri.Search(key, trie.WithMaxEditDistance(2), trie.WithTopKLeastEdited(), trie.WithMaxResults(2))
for _, res := range results.Results {
	fmt.Println(res.Key, res.Value, res.EditDistance)
}
// Output (top 2 least edited results):
// [the] 1 1
// [the green tree] 4 1
```

### References

* https://en.wikipedia.org/wiki/Levenshtein_distance#Iterative_with_full_matrix
* http://stevehanov.ca/blog/?id=114
* https://gist.github.com/jlherren/d97839b1276b9bd7faa930f74711a4b6