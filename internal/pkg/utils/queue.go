package utils

import (
	"container/list"
	"fmt"
)

// Queue implemented using container/list:
type Queue struct {
	data *list.List
}

// NewQueue creates and returns a new Queue
func NewQueue() *Queue {
	return &Queue{data: list.New()}
}

// Enqueue:
func (q *Queue) Enqueue(value int) {
	q.data.PushBack(value)
}

// Dequeue:
func (q *Queue) Dequeue() (int, error) {
	if q.IsEmpty() {
		return 0, fmt.Errorf("queue is empty")
	}
	front := q.data.Front()
	q.data.Remove(front)
	return front.Value.(int), nil
}

// Front :
func (q *Queue) Front() (int, error) {
	if q.IsEmpty() {
		return 0, fmt.Errorf("queue is empty")
	}
	return q.data.Front().Value.(int), nil
}

// IsEmpty :
func (q *Queue) IsEmpty() bool {
	return q.data.Len() == 0
}

// Size :
func (q *Queue) Size() int {
	return q.data.Len()
}
