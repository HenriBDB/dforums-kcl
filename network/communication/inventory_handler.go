package communication

import (
	"dforum-app/security"
	"sync"
)

type InventoryHandler struct {
	sync.RWMutex
	inv map[security.HashSignature]*InventoryMutex
}

// Get the lock for the corresponding hash value
// If doesn't exist, create new lock
func (ih *InventoryHandler) GetLock(id security.HashSignature) *InventoryMutex {
	ih.RLocker().Lock()
	if v, ok := ih.inv[id]; ok {
		ih.RLock()
		return v
	}
	ih.RUnlock()
	// No lock found, create new lock
	ih.Lock()
	defer ih.Unlock()
	newLock := newInventoryMutex()
	ih.inv[id] = newLock
	return newLock
}

func (ih *InventoryHandler) DeleteLock(id security.HashSignature) {
	ih.Lock()
	defer ih.Unlock()
	delete(ih.inv, id)
}

func NewInventoryHandler() InventoryHandler {
	return InventoryHandler{
		RWMutex: sync.RWMutex{},
		inv:     make(map[security.HashSignature]*InventoryMutex),
	}
}

type InventoryMutex struct {
	sync.Mutex
	success bool
}

func newInventoryMutex() *InventoryMutex {
	return &InventoryMutex{
		Mutex:   sync.Mutex{},
		success: false,
	}
}

type inventoryAction func() bool

func (inv *InventoryMutex) CompleteActionUntilSuccessful(action inventoryAction) {
	inv.Lock()
	if !inv.success {
		inv.success = action()
	}
	defer inv.Unlock()
}
