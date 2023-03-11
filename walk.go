package trie

type WalkFunc func(key []string, node *Node) error

// Walk traverses the Trie and calls walker function. If walker function returns an error, Walk early-returns with that error.
// Traversal follows insertion order.
func (t *Trie) Walk(key []string, walker WalkFunc) error {
	node := t.root
	for _, keyPart := range key {
		child, ok := node.children[keyPart]
		if !ok {
			return nil
		}
		node = child
	}
	return t.walk(node, &key, walker)
}

func (t *Trie) walk(node *Node, prefixKey *[]string, walker WalkFunc) error {
	if node.isTerminal {
		key := make([]string, len(*prefixKey))
		copy(key, *prefixKey)
		if err := walker(key, node); err != nil {
			return err
		}
	}

	for dllNode := node.childrenDLL.head; dllNode != nil; dllNode = dllNode.next {
		child := dllNode.trieNode
		*prefixKey = append(*prefixKey, child.keyPart)
		err := t.walk(child, prefixKey, walker)
		*prefixKey = (*prefixKey)[:len(*prefixKey)-1]
		if err != nil {
			return err
		}
	}
	return nil
}
