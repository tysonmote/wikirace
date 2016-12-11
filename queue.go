package main

type StringQueue struct {
	strings []string
}

func NewStringQueue() StringQueue {
	return StringQueue{[]string{}}
}

func (q *StringQueue) Enqueue(s string) {
	q.strings = append(q.strings, s)
}

func (q *StringQueue) Dequeue() (s string) {
	s = q.strings[0]
	q.strings = q.strings[1:]
	return s
}

func (q *StringQueue) IsEmpty() bool {
	return len(q.strings) == 0
}
