package simplelru

import (
	"container/list"
	"errors"
	"sync"
)

type LRU struct {
	sync.RWMutex
	size int
	evictList *list.List
	items map[interface{}]*list.Element
	onEvict func(key interface{}, value interface{})
}

func NewLRU(size int, onEvict func(key interface{}, value interface{})) (*LRU, error) {
	if size <= 0{
		return nil, errors.New("must provide a positive size")
	}
	return &LRU{
		size: size,
		evictList: list.New(),
		items: make(map[interface{}]*list.Element),
		onEvict: onEvict,
	}, nil
}

func (L *LRU) Add(key, value interface{}) bool {
	L.Lock()
	defer L.Unlock()
	if ent, ok := L.items[key];ok{
		L.evictList.MoveToFront(ent)
		ent.Value.(*entry).value=value
		return false
	}
	ent := &entry{key, value}
	entry := L.evictList.PushFront(ent)
	L.items[key] = entry

	evict := L.evictList.Len() > L.size
	if evict {
		L.removeOldest()
	}
	return evict
}

func (L *LRU) Get(key interface{}) (value interface{}, ok bool) {
	L.RLock()
	defer L.RUnlock()
	if ent, ok := L.items[key]; ok {
		L.evictList.MoveToFront(ent)
		if ent.Value.(*entry)==nil{
			return nil, false
		}
		return ent.Value.(*entry).value, true
	}
	return
}

func (L LRU) Contains(key interface{}) (ok bool) {
	L.RLock()
	defer L.RUnlock()
	_, ok = L.items[key]
	return ok
}

func (L LRU) Peek(key interface{}) (value interface{}, ok bool) {
	L.RLock()
	defer L.RUnlock()
	ent, ok := L.items[key]
	if ok {
		return ent.Value.(*entry).value, ok
	}
	return nil, false
}

func (L *LRU) Remove(key interface{}) bool {
	L.Lock()
	defer L.Unlock()
	ent, ok := L.items[key]
	if ok{
		return L.removeElement(ent)
	}
	return false
}


func (L LRU) GetOldest() (interface{}, interface{}, bool) {
	L.RLock()
	defer L.RUnlock()
	ent := L.evictList.Back()
	if ent != nil{
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}

func (L LRU) Keys() []interface{} {
	L.RLock()
	defer L.RUnlock()
	keys := make([]interface{}, len(L.items))
	i := 0
	for ent:= L.evictList.Back(); ent != nil; ent=ent.Prev() {
		keys[i] = ent.Value.(*entry).key
		i++
	}
	return keys
}

func (L LRU) Len() int {
	L.RLock()
	defer L.RUnlock()
	return L.evictList.Len()
}

func (L LRU) Purge() {
	L.Lock()
	defer L.Unlock()
	for k, v := range L.items{
		if L.onEvict != nil {
			L.onEvict(k, v.Value.(*entry).value)
		}
		delete(L.items, k)
	}
	L.evictList.Init()
}

func (L LRU) Resize(size int) int {
	L.Lock()
	defer L.Unlock()
	diff := L.evictList.Len() - size
	for i:= 0; i < diff; i++ {
		L.removeOldest()
	}
	L.size = size
	return diff
}

func (L *LRU) removeOldest() {
	ent := L.evictList.Back()
	if ent != nil {
		L.removeElement(ent)
	}
}

func (L *LRU) removeElement(e *list.Element) bool {
	ent := e.Value.(*entry)
	L.evictList.Remove(e)
	delete(L.items, ent.key)
	if L.onEvict!=nil{
		L.onEvict(ent.key, ent.value)
	}
	return true
}


type entry struct {
	key interface{}
	value interface{}
}


