package trie

// TODO: tests
type doublyLinkedList struct {
	head, tail *dllNode
}

type dllNode struct {
	trieNode   *Node
	next, prev *dllNode
}

func newDLLNode(trieNode *Node) *dllNode {
	return &dllNode{trieNode: trieNode}
}

func (dll *doublyLinkedList) append(node *dllNode) {
	if dll.head == nil {
		dll.head = node
		dll.tail = node
		return
	}
	dll.tail.next = node
	node.prev = dll.tail
	dll.tail = node
}

func (dll *doublyLinkedList) pop(node *dllNode) {
	if node == dll.head && node == dll.tail {
		dll.head = nil
		dll.tail = nil
		return
	}
	if node == dll.head {
		dll.head = node.next
		dll.head.prev = nil
		node.next = nil
		return
	}
	if node == dll.tail {
		dll.tail = node.prev
		dll.tail.next = nil
		node.prev = nil
		return
	}
	prev := node.prev
	next := node.next
	prev.next = next
	next.prev = prev
	node.next = nil
	node.prev = nil
}
