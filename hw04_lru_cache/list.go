package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

// Len - длина списка.
func (l *list) Len() int {
	return l.len
}

// Front - первый элемент списка.
func (l *list) Front() *ListItem {
	return l.front
}

// Back -последний элемент списка.
func (l *list) Back() *ListItem {
	return l.back
}

// PushFront - добавить значение в начало.
func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.front == nil {
		l.front = item
		l.back = item
	} else {
		item.Next, l.front.Prev = l.front, item
		l.front = item
	}
	l.len++
	return l.front
}

// PushBack - добавить значение в конец.
func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.back == nil {
		l.back = item
		l.front = item
	} else {
		item.Prev, l.back.Next = l.back, item
		l.back = item
	}
	l.len++
	return l.back
}

// Remove - удалить элемент.
func (l *list) Remove(i *ListItem) {
	if i == l.front {
		l.front, l.front.Next.Prev = l.front.Next, nil
	} else if i == l.back {
		l.back, l.back.Prev.Next = l.back.Prev, nil
	} else {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
		i.Next = nil
		i.Prev = nil
	}
	l.len--
}

// MoveToFront - переместить элемент в начало
func (l *list) MoveToFront(i *ListItem) {
	if i == l.front {
		return
	} else if i == l.back {
		i.Prev.Next, l.back = nil, i.Prev
	} else {
		i.Prev.Next, i.Next.Prev = i.Next, i.Prev
	}

	i.Next, i.Prev, l.front.Prev = l.front, nil, i
	l.front = i
}

func NewList() List {
	return new(list)
}
