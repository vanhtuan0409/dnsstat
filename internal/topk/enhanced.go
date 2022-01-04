package topk

func (s *Stream) Reset() {
	s.k = keys{m: make(map[string]int), elts: make([]Element, 0, s.n)}
	s.alphas = make([]int, s.n*6)
}
