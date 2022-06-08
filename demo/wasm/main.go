package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"syscall/js"

	"github.com/shivamMg/trie"
)

var (
	mu  sync.Mutex
	tri *trie.Trie
)

func init() {
	initTrie()
}

func initTrie() {
	mu.Lock()
	defer mu.Unlock()
	tri = trie.New()
	data := MustAsset("words.txt")
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

func searchWord(this js.Value, args []js.Value) interface{} {
	mu.Lock()
	defer mu.Unlock()
	word := args[0].String()
	key := strings.Split(word, "")
	results := tri.Search(key, trie.WithMaxEditDistance(3), trie.WithEditOps(),
		trie.WithTopKLeastEdited(10))
	n := len(results.Results)
	words := make([]interface{}, n)
	noEdits := make([]interface{}, n)
	for i, res := range results.Results {
		words[i] = strings.Join(res.Key, "")
		noEdits[i] = getNoEdits(res.Key, res.EditOps)
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
