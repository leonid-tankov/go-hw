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
	firstItem *ListItem
	lastItem  *ListItem
	length    int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.firstItem
}

func (l *list) Back() *ListItem {
	return l.lastItem
}

func (l *list) PushFront(value interface{}) *ListItem {
	if l.firstItem == nil {
		l.firstItem = &ListItem{
			Value: value,
			Next:  nil,
			Prev:  nil,
		}
		l.lastItem = l.firstItem
	} else {
		newItem := &ListItem{
			Value: value,
			Next:  l.firstItem,
			Prev:  nil,
		}
		l.firstItem = newItem
		newItem.Next.Prev = newItem
	}
	l.length++
	return l.firstItem
}

func (l *list) PushBack(value interface{}) *ListItem {
	if l.lastItem == nil {
		l.lastItem = &ListItem{
			Value: value,
			Next:  nil,
			Prev:  nil,
		}
		l.firstItem = l.lastItem
	} else {
		newItem := &ListItem{
			Value: value,
			Next:  nil,
			Prev:  l.lastItem,
		}
		l.lastItem = newItem
		newItem.Prev.Next = newItem
	}
	l.length++
	return l.firstItem
}

func (l *list) Remove(item *ListItem) {
	if l.length == 0 {
		return
	}
	if item.Prev == nil {
		l.firstItem = item.Next
	} else {
		item.Prev.Next = item.Next
	}
	if item.Next == nil {
		l.lastItem = item.Prev
	} else {
		item.Next.Prev = item.Prev
	}
	l.length--
}

func (l *list) MoveToFront(item *ListItem) {
	if item.Prev == nil {
		return
	}
	l.firstItem.Prev = item
	item.Prev.Next = item.Next
	if item.Next == nil {
		l.lastItem = item.Prev
		item.Next = l.firstItem
		item.Prev = nil
		l.firstItem = item
		return
	}
	item.Next.Prev = item.Prev
	item.Next = l.firstItem
	item.Prev = nil
	l.firstItem = item
}
