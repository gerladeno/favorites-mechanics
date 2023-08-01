package list

type DeLinkedList[T any] struct {
	Head *Node[T]
	Tail *Node[T]
	len  int
}

type Node[T any] struct {
	Prev  *Node[T]
	Next  *Node[T]
	Value T
}

func (l *DeLinkedList[T]) AddElement(value T, prev, next *Node[T]) {
	switch {
	case l.len == 0:
		l.init(value)
	case prev == nil && next == nil:
		l.insertLast(value)
	case prev != nil:
		l.insertAfter(value, prev)
	case next != nil:
		l.insertBefore(value, next)
	}
}

func (l *DeLinkedList[T]) insertBefore(value T, next *Node[T]) {
	node := Node[T]{
		Prev:  next.Prev,
		Next:  next,
		Value: value,
	}

	if next.Prev != nil {
		next.Prev.Next = &node
	} else {
		l.Head = &node
	}

	next.Prev = &node
	l.len++
}

func (l *DeLinkedList[T]) insertAfter(value T, prev *Node[T]) {
	node := Node[T]{
		Prev:  prev,
		Next:  prev.Next,
		Value: value,
	}

	if prev.Next != nil {
		prev.Next.Prev = &node
	} else {
		l.Tail = &node
	}

	prev.Next = &node
	l.len++
}

func (l *DeLinkedList[T]) insertLast(value T) {
	node := Node[T]{Value: value}
	l.Tail.Next = &node
	l.Tail = &node
	l.len++
}

func (l *DeLinkedList[T]) init(value T) {
	l.Head = &Node[T]{Value: value}
	l.Tail = l.Head
	l.len++
}

func (l *DeLinkedList[T]) DeleteElement(target *Node[T]) {
	if target == nil || l.len == 0 {
		return
	}

	if l.Head == target {
		l.Head = target.Next
	}

	if l.Tail == target {
		l.Tail = target.Prev
	}

	if target.Prev != nil {
		target.Prev.Next = target.Next
	}

	if target.Next != nil {
		target.Next.Prev = target.Prev
	}

	l.len--
}

func (l *DeLinkedList[T]) Len() int {
	return l.len
}

func (l *DeLinkedList[T]) List() []T {
	result := make([]T, 0, l.len)

	node := l.Head
	for node != nil {
		result = append(result, node.Value)
		node = node.Next
	}

	return result
}
