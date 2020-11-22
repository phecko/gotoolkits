package container_ext

import (
	"container/list"
	"sync"
)

type SafeList struct {
	sync.RWMutex
	L *list.List
}

func NewSafeList() *SafeList {
	return &SafeList{L: list.New()}
}

func (l *SafeList) Len() int{
	l.Lock()
	defer l.Unlock()
	return l.L.Len()
}

func (l *SafeList) Front() *list.Element {
	l.Lock()
	defer l.Unlock()
	return l.L.Front()
}

func (l *SafeList) Back() *list.Element {
	l.Lock()
	defer l.Unlock()
	return l.L.Back()
}

func (l *SafeList) Remove(e *list.Element) interface{} {
	l.Lock()
	defer l.Unlock()
	return l.L.Remove(e)
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *SafeList) PushFront(v interface{}) interface{} {
	l.Lock()
	defer l.Unlock()
	return l.L.PushFront(v)
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *SafeList) PushBack(v interface{}) *list.Element {
	l.Lock()
	defer l.Unlock()
	return l.L.PushBack(v)
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *SafeList) InsertBefore(v interface{}, mark *list.Element) *list.Element {
	l.Lock()
	defer l.Unlock()
	return l.L.InsertBefore(v, mark)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *SafeList) InsertAfter(v interface{}, mark *list.Element) *list.Element {
	l.Lock()
	defer l.Unlock()
	return l.L.InsertAfter(v, mark)
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *SafeList) MoveToFront(e *list.Element) {
	l.Lock()
	defer l.Unlock()
	l.L.MoveToFront(e)
}

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *SafeList) MoveToBack(e *list.Element) {
	l.Lock()
	defer l.Unlock()
	l.L.MoveToBack(e)
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *SafeList) MoveBefore(e, mark *list.Element) {
	l.Lock()
	defer l.Unlock()
	l.L.MoveBefore(e, mark)
}

// MoveAfter moves element e to its new position after mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *SafeList) MoveAfter(e, mark *list.Element) {
	l.Lock()
	defer l.Unlock()
	l.L.MoveAfter(e, mark)
}

// PushBackBatch
func (l *SafeList) PushBackBatch(others []interface{}) {
	l.Lock()
	defer l.Unlock()
	for _, other := range others{
		l.L.PushBack(other)
	}
}

// PushFrontBatch inserts a copy of an other list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (l *SafeList) PushFrontBatch(others []interface{}) {
	l.Lock()
	defer l.Unlock()
	for i := l.Len()-1; i>=0; i-- {
		l.L.PushFront(others[i])
	}
}
