package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	byItems  map[*ListItem]Key
	mu       sync.Mutex
}

// Clear - Очистить кэш.
func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
	l.byItems = make(map[*ListItem]Key, l.capacity)
}

// Set - Добавить значение в кэш по ключу.
func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	// если элемент присутствует в словаре, то обновить его значение и переместить элемент в начало очереди;
	if el, ok := l.items[key]; ok {
		delete(l.byItems, el)
		el.Value = value
		l.queue.MoveToFront(el)
		l.byItems[el] = key
		return true
	}

	// если элемента нет в словаре, то добавить в словарь и в начало очереди (при этом, если размер очереди больше
	// ёмкости кэша, то необходимо удалить последний элемент из очереди и его значение из словаря);
	l.queue.PushFront(value)
	l.items[key] = l.queue.Front()
	l.byItems[l.queue.Front()] = key

	if l.queue.Len() > l.capacity {
		k := l.byItems[l.queue.Back()]
		delete(l.items, k)
		delete(l.byItems, l.queue.Back())
		l.queue.Remove(l.queue.Back())
	}

	return false
}

// Get - Получить значение из кэша по ключу.
func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// если элемент присутствует в словаре, то переместить элемент в начало очереди и вернуть его значение и true;
	if el, ok := l.items[key]; ok {
		l.queue.MoveToFront(el)
		return el.Value, true
	}

	return nil, false
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		byItems:  make(map[*ListItem]Key, capacity),
	}
}
