// Package storage provides the basic API for storing data persistently or caching it in memory.
package storage

import (
	"dforum-app/security"
	"math/rand"
	"time"
)

type NewNodeListener interface {
	RegisterNewNode(*Node)
}

type StorageModule struct {
	cache     StorageCache
	db        Database
	listeners []NewNodeListener
}

func NewStorageModule(pathToDb string) *StorageModule {
	db := NewLevelDbImpl()
	db.InitDatabase(pathToDb)

	return &StorageModule{
		cache: NewStorageCache(),
		db:    db,
	}
}

/*
PUBLIC API for the storage package
*/

// Check whether a given node is already stored locally.
// Hashes searched are added to teh cache.
func (s *StorageModule) NodeExists(nodeHash security.HashSignature) bool {
	if s.cache.containsNode(nodeHash) {
		return true
	}
	// Check the database last
	exists := s.db.HasNode(nodeHash)
	if exists {
		//Add node hash to cache
		s.cache.addNodeHash(nodeHash)
		return true
	}
	return false
}

func (s *StorageModule) StoreAndRegisterNewNode(n *Node) {
	s.StoreNode(n)
	s.PublishNode(n)
}

// Store a given node in the database
func (s *StorageModule) StoreNode(n *Node) {
	// Add the node to the database
	s.db.StoreNode(n)
	s.cache.addNode(n)
}

func (s *StorageModule) PublishNode(n *Node) {
	for _, v := range s.listeners {
		v.RegisterNewNode(n)
	}
}

func (s *StorageModule) Subscribe(n NewNodeListener) {
	s.listeners = append(s.listeners, n)
}

// Retrieve a node from the local storage to be shared with peers.
// First checks the cache and then the database in case of a miss on the cache.
// All nodes retrieved for sharing are cached in case another peer asks for them.
// Returns an empty node if not found.
func (s *StorageModule) GetNode(id security.HashSignature, shouldCache bool) *Node {
	// Return node from cache if found
	if node, exists := s.cache.getNode(id); exists {
		return node
	}
	// Else look for node in database
	if node, exists := s.db.GetNode(id); exists {
		// Add node to cache before sharing
		if shouldCache {
			s.cache.addNode(node)
		}
		return node
	}
	return nil
}

func (s *StorageModule) GetTopLevelNodes() []*Node {
	// Children of a 0 hash are top level nodes
	nodeSlice := []*Node{}
	topLevelNodes := s.db.GetChildren(security.HashSignature{})
	for _, n := range topLevelNodes {
		nodeSlice = append(nodeSlice, s.GetNode(n, false))
	}
	return nodeSlice
}

// Retrieve nodes from the database to be displayed on the GUI.
// Paginates by breadth first.
func (s *StorageModule) GetChildrenNodes(parent security.HashSignature, includeParent bool, max int8) []*Node {
	// Default to 50 child nodes per query
	if max <= 0 {
		max = 50
	}

	fetchedNodes := []*Node{}
	childrenHashes := s.db.GetChildren(parent)

	// Shuffle children
	for i := range fetchedNodes {
		j := rand.Intn(i + 1)
		fetchedNodes[i], fetchedNodes[j] = fetchedNodes[j], fetchedNodes[i]
	}

	if includeParent {
		fetchedNodes = append(fetchedNodes, s.GetNode(parent, false))
	}

	for i, v := range childrenHashes {
		if i >= int(max) {
			break
		}
		fetchedNodes = append(fetchedNodes, s.GetNode(v, false))
	}

	return fetchedNodes
}

func (s *StorageModule) GetNodesSince(t time.Time) []security.HashSignature {
	return s.db.GetAllNodesSince(t)
}

func (s *StorageModule) TimeOfMostRecentNode() time.Time {
	return s.db.TimeOfMostRecentNode()
}

func (s *StorageModule) TearDown() {
	// Close Database
	s.db.Close()
}
