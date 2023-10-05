package stack

//last in, first out stack

import (

)

type Stack[V any] struct {
	stack []V
	size int
}

func New[V any](sz int) *Stack[V] {
	return &Stack[V] {
		//stack: []V{},
		stack: make([]V, 0, sz),
	}
}

//Pushed value onto the top of the stack and returns index of recently added element
func (s *Stack[V]) Push(val V) (index int) {
	s.stack = append(s.stack, val)
	s.size++
	index = s.size-1
	return
}

//Returns the index and value of the last element on the stack.
func (s *Stack[V]) Back() (int, V) {
	return s.size-1, s.stack[s.size-1]
}

//Returns the index and value of the last element on the stack.
func (s *Stack[V]) BackIndex() (int) {
	return s.size-1
}

//Returns the index of the first element on the stack.
func (s *Stack[V]) Front() (int, V) {
	return 0, s.stack[0]
}

//Returns the index of the first element on the stack.
func (s *Stack[V]) FrontIndex() (int) {
	return 0
}

//Remove from the bottom of the stack. Returns index of element that was removed.
func (s *Stack[V]) Pop() (index int, val V) {
	if s.size == 0 {
		index = 0
		return
	}

	index, val = s.size-1, s.stack[s.size-1]
	s.stack = s.stack[:s.size-1]
	s.size--
	return
}

//Remove a element by the index on the stack.
func (s *Stack[V]) Remove(index int) {
	s.stack = append(s.stack[:index], s.stack[index+1:]...)
	s.size--
}

//Moves a element to the back of the stack.
func (s *Stack[V]) MoveToBack(index int) (int) {
	if s.size <= 1 {
		return 0
	}

	val := s.stack[index]
	s.Remove(index)
	return s.Push(val)
}

//Clears the entire stack but keep allocated memory
func (s *Stack[V]) Clear() {
	s.stack = s.stack[:0]
}

func (s *Stack[V]) Stack() []V {
	return s.stack
}