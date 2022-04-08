package storage

import (
	"dforum-app/security"
	"sync"
)

/*
	Thread safe and type safe implementation of a hashmap
	from node fingerprint to node pointer
*/

// Inspired by https://stackoverflow.com/questions/36167200/how-safe-are-golang-maps-for-concurrent-read-write-operations
// The above answer is a simplification of Go's sync map core library.
// A similar implementation can also be found here under MIT licence: https://github.com/orcaman/concurrent-map/blob/master/concurrent_map.go
type ConcurrentHashMap struct {
	content map[security.HashSignature]*Node
	lock    sync.RWMutex
}

func NewConcurrentHashMap() ConcurrentHashMap {
	return ConcurrentHashMap{
		content: make(map[security.HashSignature]*Node),
		lock:    sync.RWMutex{},
	}
}

/* NodeCache functionalities */

func (n *ConcurrentHashMap) Get(id security.HashSignature) (*Node, bool) {
	n.lock.RLock()
	defer n.lock.RUnlock()
	node, exists := n.content[id]
	return node, exists
}

func (n *ConcurrentHashMap) Has(id security.HashSignature) bool {
	n.lock.RLock()
	defer n.lock.RUnlock()
	_, ok := n.content[id]
	return ok
}

func (n *ConcurrentHashMap) Set(id security.HashSignature, node *Node) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.content[id] = node
}
