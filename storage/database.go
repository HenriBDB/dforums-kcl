package storage

import (
	"dforum-app/security"
	"time"
)

type Database interface {
	HasNode(security.HashSignature) bool
	GetNode(security.HashSignature) (*Node, bool)
	GetChildren(security.HashSignature) []security.HashSignature
	GetAllNodesSince(time.Time) []security.HashSignature
	StoreNode(*Node) bool
	TimeOfMostRecentNode() time.Time
	InitDatabase(pathToFiles string) error
	Close()
}
