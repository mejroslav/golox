package utils

type Stack struct {
	elements []any
}

func NewStack() *Stack {
	return &Stack{
		elements: make([]any, 0, 8), // Pre-allocate capacity
	}
}

func (s *Stack) Get(index int) (any, bool) {
	if index < 0 || index >= len(s.elements) {
		return nil, false
	}
	return s.elements[index], true
}

func (s *Stack) Push(element any) {
	s.elements = append(s.elements, element)
}

func (s *Stack) Pop() any {
	if len(s.elements) == 0 {
		return nil
	}
	last := s.elements[len(s.elements)-1]
	s.elements[len(s.elements)-1] = nil // Clear reference for GC
	s.elements = s.elements[:len(s.elements)-1]
	return last
}

func (s *Stack) Peek() any {
	if len(s.elements) == 0 {
		return nil
	}
	return s.elements[len(s.elements)-1]
}

func (s *Stack) IsEmpty() bool {
	return len(s.elements) == 0
}

func (s *Stack) Size() int {
	return len(s.elements)
}
