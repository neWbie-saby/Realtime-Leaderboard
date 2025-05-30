package ds

func NewHashMap[K comparable, V any]() *HashMap[K, V] {
	return &HashMap[K, V]{
		data: make(map[K]V),
	}
}

func (hm *HashMap[K, V]) Set(key K, value V) {
	hm.data[key] = value
}

func (hm *HashMap[K, V]) Get(key K) (V, bool) {
	val, ok := hm.data[key]
	return val, ok
}

func (hm *HashMap[K, V]) Delete(key K) {
	delete(hm.data, key)
}

func (hm *HashMap[K, V]) Items() map[K]V {
	return hm.data
}

func (hm *HashMap[K, V]) Len() int {
	return len(hm.data)
}
