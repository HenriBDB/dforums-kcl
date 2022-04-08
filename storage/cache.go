// The cache file is used to control cache mechanisms: adding to cache, cache limits and querying the cache
package storage

import "dforum-app/security"

// The purpose of this cache is to store nodes recently shared or requested
type StorageCache struct {
	nodeHashCache ConcurrentHashMap
	nodeCache     ConcurrentHashMap
}

func NewStorageCache() StorageCache {
	return StorageCache{
		nodeHashCache: NewConcurrentHashMap(),
		nodeCache:     NewConcurrentHashMap(),
	}
}

func (s *StorageCache) containsNode(id security.HashSignature) bool {
	// Check the hash cache first
	if s.nodeHashCache.Has(id) {
		return true
	}
	// Check the node cache second
	if s.nodeCache.Has(id) {
		return true
	}
	return false
}

func (s *StorageCache) getNode(id security.HashSignature) (*Node, bool) {
	return s.nodeCache.Get(id)
}

func (s *StorageCache) addNodeHash(id security.HashSignature) {
	s.nodeHashCache.Set(id, nil)
}

func (s *StorageCache) addNode(node *Node) {
	s.nodeHashCache.Set(node.GetFingerprint(), node)
}
