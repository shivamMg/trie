package trie_test

import (
	"fmt"

	"github.com/shivamMg/trie"
)

func printEditOps(ops []*trie.EditOp) {
	for _, op := range ops {
		switch op.Type {
		case trie.EditOpTypeNoEdit:
			fmt.Printf("- don't edit %q\n", op.KeyPart)
		case trie.EditOpTypeInsert:
			fmt.Printf("- insert %q\n", op.KeyPart)
		case trie.EditOpTypeDelete:
			fmt.Printf("- delete %q\n", op.KeyPart)
		case trie.EditOpTypeReplace:
			fmt.Printf("- replace %q with %q\n", op.KeyPart, op.ReplaceWith)
		}
	}
}

func Example() {
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
}
