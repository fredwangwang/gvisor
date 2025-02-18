package gofer

// ElementMapper provides an identity mapping by default.
//
// This can be replaced to provide a struct that maps elements to linker
// objects, if they are not the same. An ElementMapper is not typically
// required if: Linker is left as is, Element is left as is, or Linker and
// Element are the same type.
type stringElementMapper struct{}

// linkerFor maps an Element to a Linker.
//
// This default implementation should be inlined.
//
//go:nosplit
func (stringElementMapper) linkerFor(elem *stringListElem) *stringListElem { return elem }

// List is an intrusive list. Entries can be added to or removed from the list
// in O(1) time and with no additional memory allocations.
//
// The zero value for List is an empty list ready to use.
//
// To iterate over a list (where l is a List):
//
//	for e := l.Front(); e != nil; e = e.Next() {
//		// do something with e.
//	}
//
// +stateify savable
type stringList struct {
	head *stringListElem
	tail *stringListElem
}

// Reset resets list l to the empty state.
func (l *stringList) Reset() {
	l.head = nil
	l.tail = nil
}

// Empty returns true iff the list is empty.
//
//go:nosplit
func (l *stringList) Empty() bool {
	return l.head == nil
}

// Front returns the first element of list l or nil.
//
//go:nosplit
func (l *stringList) Front() *stringListElem {
	return l.head
}

// Back returns the last element of list l or nil.
//
//go:nosplit
func (l *stringList) Back() *stringListElem {
	return l.tail
}

// Len returns the number of elements in the list.
//
// NOTE: This is an O(n) operation.
//
//go:nosplit
func (l *stringList) Len() (count int) {
	for e := l.Front(); e != nil; e = (stringElementMapper{}.linkerFor(e)).Next() {
		count++
	}
	return count
}

// PushFront inserts the element e at the front of list l.
//
//go:nosplit
func (l *stringList) PushFront(e *stringListElem) {
	linker := stringElementMapper{}.linkerFor(e)
	linker.SetNext(l.head)
	linker.SetPrev(nil)
	if l.head != nil {
		stringElementMapper{}.linkerFor(l.head).SetPrev(e)
	} else {
		l.tail = e
	}

	l.head = e
}

// PushFrontList inserts list m at the start of list l, emptying m.
//
//go:nosplit
func (l *stringList) PushFrontList(m *stringList) {
	if l.head == nil {
		l.head = m.head
		l.tail = m.tail
	} else if m.head != nil {
		stringElementMapper{}.linkerFor(l.head).SetPrev(m.tail)
		stringElementMapper{}.linkerFor(m.tail).SetNext(l.head)

		l.head = m.head
	}
	m.head = nil
	m.tail = nil
}

// PushBack inserts the element e at the back of list l.
//
//go:nosplit
func (l *stringList) PushBack(e *stringListElem) {
	linker := stringElementMapper{}.linkerFor(e)
	linker.SetNext(nil)
	linker.SetPrev(l.tail)
	if l.tail != nil {
		stringElementMapper{}.linkerFor(l.tail).SetNext(e)
	} else {
		l.head = e
	}

	l.tail = e
}

// PushBackList inserts list m at the end of list l, emptying m.
//
//go:nosplit
func (l *stringList) PushBackList(m *stringList) {
	if l.head == nil {
		l.head = m.head
		l.tail = m.tail
	} else if m.head != nil {
		stringElementMapper{}.linkerFor(l.tail).SetNext(m.head)
		stringElementMapper{}.linkerFor(m.head).SetPrev(l.tail)

		l.tail = m.tail
	}
	m.head = nil
	m.tail = nil
}

// InsertAfter inserts e after b.
//
//go:nosplit
func (l *stringList) InsertAfter(b, e *stringListElem) {
	bLinker := stringElementMapper{}.linkerFor(b)
	eLinker := stringElementMapper{}.linkerFor(e)

	a := bLinker.Next()

	eLinker.SetNext(a)
	eLinker.SetPrev(b)
	bLinker.SetNext(e)

	if a != nil {
		stringElementMapper{}.linkerFor(a).SetPrev(e)
	} else {
		l.tail = e
	}
}

// InsertBefore inserts e before a.
//
//go:nosplit
func (l *stringList) InsertBefore(a, e *stringListElem) {
	aLinker := stringElementMapper{}.linkerFor(a)
	eLinker := stringElementMapper{}.linkerFor(e)

	b := aLinker.Prev()
	eLinker.SetNext(a)
	eLinker.SetPrev(b)
	aLinker.SetPrev(e)

	if b != nil {
		stringElementMapper{}.linkerFor(b).SetNext(e)
	} else {
		l.head = e
	}
}

// Remove removes e from l.
//
//go:nosplit
func (l *stringList) Remove(e *stringListElem) {
	linker := stringElementMapper{}.linkerFor(e)
	prev := linker.Prev()
	next := linker.Next()

	if prev != nil {
		stringElementMapper{}.linkerFor(prev).SetNext(next)
	} else if l.head == e {
		l.head = next
	}

	if next != nil {
		stringElementMapper{}.linkerFor(next).SetPrev(prev)
	} else if l.tail == e {
		l.tail = prev
	}

	linker.SetNext(nil)
	linker.SetPrev(nil)
}

// Entry is a default implementation of Linker. Users can add anonymous fields
// of this type to their structs to make them automatically implement the
// methods needed by List.
//
// +stateify savable
type stringEntry struct {
	next *stringListElem
	prev *stringListElem
}

// Next returns the entry that follows e in the list.
//
//go:nosplit
func (e *stringEntry) Next() *stringListElem {
	return e.next
}

// Prev returns the entry that precedes e in the list.
//
//go:nosplit
func (e *stringEntry) Prev() *stringListElem {
	return e.prev
}

// SetNext assigns 'entry' as the entry that follows e in the list.
//
//go:nosplit
func (e *stringEntry) SetNext(elem *stringListElem) {
	e.next = elem
}

// SetPrev assigns 'entry' as the entry that precedes e in the list.
//
//go:nosplit
func (e *stringEntry) SetPrev(elem *stringListElem) {
	e.prev = elem
}

// RingInit instantiates an Element to be an item in a ring (circularly-linked
// list).
//
//go:nosplit
func stringRingInit(e *stringListElem) {
	linker := stringElementMapper{}.linkerFor(e)
	linker.SetNext(e)
	linker.SetPrev(e)
}

// RingAdd adds new to old's ring.
//
//go:nosplit
func stringRingAdd(old *stringListElem, new *stringListElem) {
	oldLinker := stringElementMapper{}.linkerFor(old)
	newLinker := stringElementMapper{}.linkerFor(new)
	next := oldLinker.Next()
	prev := old

	next.SetPrev(new)
	newLinker.SetNext(next)
	newLinker.SetPrev(prev)
	oldLinker.SetNext(new)
}

// RingRemove removes e from its ring.
//
//go:nosplit
func stringRingRemove(e *stringListElem) {
	eLinker := stringElementMapper{}.linkerFor(e)
	next := eLinker.Next()
	prev := eLinker.Prev()
	next.SetPrev(prev)
	prev.SetNext(next)
	stringRingInit(e)
}

// RingEmpty returns true if there are no other elements in the list.
//
//go:nosplit
func stringRingEmpty(e *stringListElem) bool {
	linker := stringElementMapper{}.linkerFor(e)
	return linker.Next() == e
}
