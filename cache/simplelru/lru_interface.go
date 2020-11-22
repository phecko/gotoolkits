package simplelru

type LRUCache interface {

	Add(key, value interface{}) bool

	Get(key interface{}) (value interface{},ok bool)

	Contains(key interface{}) (ok bool)

	Peek(key interface{}) (value interface{}, ok bool)

	Remove(key interface{}) bool

	GetOldest() (interface{}, interface{}, bool)

	Keys() []interface{}

	Len() int

	Purge()

	Resize(int) int
}
