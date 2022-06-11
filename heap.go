package trie

type searchResultMaxHeap []*SearchResult

func (s searchResultMaxHeap) Len() int {
	return len(s)
}

func (s searchResultMaxHeap) Less(i, j int) bool {
	if s[i].EditDistance == s[j].EditDistance {
		return s[i].tiebreaker > s[j].tiebreaker
	}
	return s[i].EditDistance > s[j].EditDistance
}

func (s searchResultMaxHeap) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *searchResultMaxHeap) Push(x interface{}) {
	*s = append(*s, x.(*SearchResult))
}

func (s *searchResultMaxHeap) Pop() interface{} {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}
