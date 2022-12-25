package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"strings"
	"sync"
	"syscall/js"

	"github.com/shivamMg/trie"
)

//go:embed words.txt
var data []byte

var (
	mu  sync.Mutex
	tri *trie.Trie

	longestWordLen int
)

func init() {
	initTrie()
}

func initTrie() {
	mu.Lock()
	defer mu.Unlock()
	tri = trie.New()
	r := bufio.NewReader(bytes.NewReader(data))
	for {
		word, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		word = strings.TrimRight(word, "\n")
		if len(word) > longestWordLen {
			longestWordLen = len(word)
		}
		key := strings.Split(word, "")
		tri.Put(key, struct{}{})
	}
}

func getNoEdits(key []string, ops []*trie.EditOp) []interface{} {
	uneditedLetters := make([]string, 0)
	for _, op := range ops {
		if op.Type == trie.EditOpTypeNoEdit {
			uneditedLetters = append(uneditedLetters, op.KeyPart)
		}
	}
	noEdits := make([]interface{}, len(key))
	j := 0
	for i, letter := range key {
		unedited := false
		if j < len(uneditedLetters) && letter == uneditedLetters[j] {
			unedited = true
			j += 1
		}
		noEdits[i] = unedited
	}
	return noEdits
}

func getNoEditsForPrefixSearch(wordLen int, keyLen int) []interface{} {
	noEdits := make([]interface{}, keyLen)
	for i := 0; i < keyLen; i++ {
		noEdits[i] = i < wordLen
	}
	return noEdits
}

func searchWord(this js.Value, args []js.Value) interface{} {
	mu.Lock()
	defer mu.Unlock()
	word := args[0].String()
	if len(word) > longestWordLen {
		return map[string]interface{}{
			"words":   []interface{}{},
			"noEdits": []interface{}{},
		}
	}
	approximate := args[1].Bool()
	key := strings.Split(word, "")
	opts := []func(*trie.SearchOptions){trie.WithMaxResults(10)}
	if approximate {
		opts = append(opts, trie.WithMaxEditDistance(3), trie.WithEditOps(), trie.WithTopKLeastEdited())
	}
	results := tri.Search(key, opts...)
	n := len(results.Results)
	words := make([]interface{}, n)
	noEdits := make([]interface{}, n)
	for i, res := range results.Results {
		words[i] = strings.Join(res.Key, "")
		if approximate {
			noEdits[i] = getNoEdits(res.Key, res.EditOps)
		} else {
			noEdits[i] = getNoEditsForPrefixSearch(len(word), len(res.Key))
		}
	}
	return map[string]interface{}{
		"words":   words,
		"noEdits": noEdits,
	}
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("searchWord", js.FuncOf(searchWord))
	<-c
}
