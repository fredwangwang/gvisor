package stack

// ElementMapper provides an identity mapping by default.
//
// This can be replaced to provide a struct that maps elements to linker
// objects, if they are not the same. An ElementMapper is not typically
// required if: Linker is left as is, Element is left as is, or Linker and
// Element are the same type.
type groPacketElementMapper struct{}

// linkerFor maps an Element to a Linker.
//
// This default implementation should be inlined.
//
//go:nosplit
func (groPacketElementMapper) linkerFor(elem *groPacket) *groPacket { return elem }

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
type groPacketList struct {
	head *groPacket
	tail *groPacket
}

// Reset resets list l to the empty state.
func (l *groPacketList) Reset() {
	l.head = nil
	l.tail = nil
}

// Empty returns true iff the list is empty.
//
//go:nosplit
func (l *groPacketList) Empty() bool {
	return l.head == nil
}

// Front returns the first element of list l or nil.
//
//go:nosplit
func (l *groPacketList) Front() *groPacket {
	return l.head
}

// Back returns the last element of list l or nil.
//
//go:nosplit
func (l *groPacketList) Back() *groPacket {
	return l.tail
}

// Len returns the number of elements in the list.
//
// NOTE: This is an O(n) operation.
//
//go:nosplit
func (l *groPacketList) Len() (count int) {
	for e := l.Front(); e != nil; e = (groPacketElementMapper{}.linkerFor(e)).Next() {
		count++
	}
	return count
}

// PushFront inserts the element e at the front of list l.
//
//go:nosplit
func (l *groPacketList) PushFront(e *groPacket) {
	linker := groPacketElementMapper{}.linkerFor(e)
	linker.SetNext(l.head)
	linker.SetPrev(nil)
	if l.head != nil {
		groPacketElementMapper{}.linkerFor(l.head).SetPrev(e)
	} else {
		l.tail = e
	}

	l.head = e
}

// PushFrontList inserts list m at the start of list l, emptying m.
//
//go:nosplit
func (l *groPacketList) PushFrontList(m *groPacketList) {
	if l.head == nil {
		l.head = m.head
		l.tail = m.tail
	} else if m.head != nil {
		groPacketElementMapper{}.linkerFor(l.head).SetPrev(m.tail)
		groPacketElementMapper{}.linkerFor(m.tail).SetNext(l.head)

		l.head = m.head
	}
	m.head = nil
	m.tail = nil
}

// PushBack inserts the element e at the back of list l.
//
//go:nosplit
func (l *groPacketList) PushBack(e *groPacket) {
	linker := groPacketElementMapper{}.linkerFor(e)
	linker.SetNext(nil)
	linker.SetPrev(l.tail)
	if l.tail != nil {
		groPacketElementMapper{}.linkerFor(l.tail).SetNext(e)
	} else {
		l.head = e
	}

	l.tail = e
}

// PushBackList inserts list m at the end of list l, emptying m.
//
//go:nosplit
func (l *groPacketList) PushBackList(m *groPacketList) {
	if l.head == nil {
		l.head = m.head
		l.tail = m.tail
	} else if m.head != nil {
		groPacketElementMapper{}.linkerFor(l.tail).SetNext(m.head)
		groPacketElementMapper{}.linkerFor(m.head).SetPrev(l.tail)

		l.tail = m.tail
	}
	m.head = nil
	m.tail = nil
}

// InsertAfter inserts e after b.
//
//go:nosplit
func (l *groPacketList) InsertAfter(b, e *groPacket) {
	bLinker := groPacketElementMapper{}.linkerFor(b)
	eLinker := groPacketElementMapper{}.linkerFor(e)

	a := bLinker.Next()

	eLinker.SetNext(a)
	eLinker.SetPrev(b)
	bLinker.SetNext(e)

	if a != nil {
		groPacketElementMapper{}.linkerFor(a).SetPrev(e)
	} else {
		l.tail = e
	}
}

// InsertBefore inserts e before a.
//
//go:nosplit
func (l *groPacketList) InsertBefore(a, e *groPacket) {
	aLinker := groPacketElementMapper{}.linkerFor(a)
	eLinker := groPacketElementMapper{}.linkerFor(e)

	b := aLinker.Prev()
	eLinker.SetNext(a)
	eLinker.SetPrev(b)
	aLinker.SetPrev(e)

	if b != nil {
		groPacketElementMapper{}.linkerFor(b).SetNext(e)
	} else {
		l.head = e
	}
}

// Remove removes e from l.
//
//go:nosplit
func (l *groPacketList) Remove(e *groPacket) {
	linker := groPacketElementMapper{}.linkerFor(e)
	prev := linker.Prev()
	next := linker.Next()

	if prev != nil {
		groPacketElementMapper{}.linkerFor(prev).SetNext(next)
	} else if l.head == e {
		l.head = next
	}

	if next != nil {
		groPacketElementMapper{}.linkerFor(next).SetPrev(prev)
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
type groPacketEntry struct {
	next *groPacket
	prev *groPacket
}

// Next returns the entry that follows e in the list.
//
//go:nosplit
func (e *groPacketEntry) Next() *groPacket {
	return e.next
}

// Prev returns the entry that precedes e in the list.
//
//go:nosplit
func (e *groPacketEntry) Prev() *groPacket {
	return e.prev
}

// SetNext assigns 'entry' as the entry that follows e in the list.
//
//go:nosplit
func (e *groPacketEntry) SetNext(elem *groPacket) {
	e.next = elem
}

// SetPrev assigns 'entry' as the entry that precedes e in the list.
//
//go:nosplit
func (e *groPacketEntry) SetPrev(elem *groPacket) {
	e.prev = elem
}

// RingInit instantiates an Element to be an item in a ring (circularly-linked
// list).
//
//go:nosplit
func groPacketRingInit(e *groPacket) {
	linker := groPacketElementMapper{}.linkerFor(e)
	linker.SetNext(e)
	linker.SetPrev(e)
}

// RingAdd adds new to old's ring.
//
//go:nosplit
func groPacketRingAdd(old *groPacket, new *groPacket) {
	oldLinker := groPacketElementMapper{}.linkerFor(old)
	newLinker := groPacketElementMapper{}.linkerFor(new)
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
func groPacketRingRemove(e *groPacket) {
	eLinker := groPacketElementMapper{}.linkerFor(e)
	next := eLinker.Next()
	prev := eLinker.Prev()
	next.SetPrev(prev)
	prev.SetNext(next)
	groPacketRingInit(e)
}

// RingEmpty returns true if there are no other elements in the list.
//
//go:nosplit
func groPacketRingEmpty(e *groPacket) bool {
	linker := groPacketElementMapper{}.linkerFor(e)
	return linker.Next() == e
}
