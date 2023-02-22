package structure

import "sync"

// ConcurrentMapMap 并发安全的map[interface{}]map[interface{}]interface{}
type ConcurrentMapMap struct {
	storage sync.Map
}

// Push push val
func (c *ConcurrentMapMap) Push(key1, key2, val interface{}) {
	v, _ := c.storage.LoadOrStore(key1, &sync.Map{})
	m := v.(*sync.Map)
	m.Store(key2, val)
}

// Delete delete the key
func (c *ConcurrentMapMap) Delete(key interface{}) {
	c.storage.Delete(key)
}

// Pop pop val
func (c *ConcurrentMapMap) Pop(key1, key2 interface{}) {
	if val, ok := c.storage.Load(key1); ok {
		val.(*sync.Map).Delete(key2)
	}
}

// GetValue get val
func (c *ConcurrentMapMap) GetValue(key1, key2 interface{}) (interface{}, bool) {
	if val, ok := c.storage.Load(key1); ok {
		return val.(*sync.Map).Load(key2)
	}
	return nil, false
}

// Get get map
func (c *ConcurrentMapMap) Get(key1 interface{}) (*sync.Map, bool) {
	if val, ok := c.storage.Load(key1); ok {
		return val.(*sync.Map), true
	}
	return nil, false
}
