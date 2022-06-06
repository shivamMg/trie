package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strings"

	"github.com/shivamMg/trie"
)

//go:generate go-bindata -nometadata -nocompress words.txt

var data []byte

func init() {
	data = MustAsset("words.txt")
}

func main() {
	tri := trie.New()
	r := bufio.NewReader(bytes.NewReader(data))
	for {
		word, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		word = strings.TrimRight(word, "\n")
		key := strings.Split(word, "")
		tri.Put(key, word)
	}

	log.Println("begin")
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		text := s.Text()
		if text == "" {
			break
		}
		key := strings.Split(text, "")
		results := tri.Search(key, trie.WithMaxEditDistance(4), trie.WithTopKLeastEdited(10))
		for _, res := range results.Results {
			log.Println(res.Key, res.EditCount)
		}
	}
	if s.Err() != nil {
		log.Fatal(s.Err())
	}
}
