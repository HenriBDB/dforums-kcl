package storage

import (
	"dforum-app/configuration"
	"dforum-app/security"
	"encoding/binary"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LevelDbImpl struct {
	nodeDB *leveldb.DB
	// This database stores parent - child relationships by storing key value pairs where the key = parent hash + child hash and the value is empty.
	// This allows to iterate through indexes in a lexicographic order and get all children of a node.
	edgeDB *leveldb.DB
	// This database indexes node hash values by time stamp.
	timestampDB *leveldb.DB
}

func NewLevelDbImpl() *LevelDbImpl {
	return &LevelDbImpl{}
}

func (db *LevelDbImpl) HasNode(id security.HashSignature) bool {
	ok, _ := db.nodeDB.Has(id[:], nil)
	return ok
}

func (db *LevelDbImpl) GetNode(id security.HashSignature) (*Node, bool) {
	nodeBytes, err := db.nodeDB.Get(id[:], nil)
	if err == leveldb.ErrNotFound {
		return &Node{}, false
	}
	//TODO check for parsing errors
	node := ParseNode(nodeBytes)
	return node, true
}

func (db *LevelDbImpl) GetChildren(id security.HashSignature) []security.HashSignature {
	children := []security.HashSignature{}

	iter := db.edgeDB.NewIterator(util.BytesPrefix(id[:]), nil)
	for iter.Next() {
		children = append(children, *(*[28]byte)(iter.Key()[28:56]))
	}
	iter.Release()

	return children
}

func (db *LevelDbImpl) GetAllNodesSince(t time.Time) []security.HashSignature {
	nodes := []security.HashSignature{}

	timeStart := make([]byte, 8)
	binary.BigEndian.PutUint64(timeStart, uint64(t.Unix()))
	timeEnd := make([]byte, 8)
	binary.BigEndian.PutUint64(timeEnd, uint64(time.Now().Unix()))

	iter := db.timestampDB.NewIterator(&util.Range{Start: timeStart, Limit: timeEnd}, nil)

	for iter.Next() {
		nodes = append(nodes, *(*[28]byte)(iter.Key()[8:]))
	}
	iter.Release()

	return nodes
}

func (db *LevelDbImpl) StoreNode(n *Node) bool {
	nodeId := n.GetFingerprint()
	// Add node's timestamp index to the time indexed table
	time := make([]byte, 8)
	binary.BigEndian.PutUint64(time, uint64(n.GetTimestamp()))
	if err := db.timestampDB.Put(append(time, nodeId[:]...), nil, nil); err != nil {
		configuration.Logger.Errorf("could not add the timestamp of node %s to the database: %s", nodeId[0:4], err.Error())
		return false
	}
	// Add the node's parent relationship to the database
	if err := db.edgeDB.Put(append(n.DatObj.Parent[:], nodeId[:]...), nil, nil); err != nil {
		configuration.Logger.Errorf("could not add the edge of node %s to the database: %s", nodeId[0:4], err.Error())
		return false
	}
	if err := db.nodeDB.Put(nodeId[:], n.GetBytes(), nil); err != nil {
		configuration.Logger.Errorf("could not add the node %s to the database: %s", nodeId[0:4], err.Error())
		return false
	}
	return true
}

func (db *LevelDbImpl) TimeOfMostRecentNode() time.Time {
	iter := db.timestampDB.NewIterator(nil, nil)
	iter.Last()
	k := iter.Key()
	if k == nil {
		// No key found
		return time.Now().AddDate(0, 0, -14)
	}
	epochTime := int64(binary.BigEndian.Uint32(k[0:8]))
	return time.Unix(epochTime, 0)
}

func (db *LevelDbImpl) InitDatabase(pathToFiles string) error {
	// Open DB Files
	nodes, err := leveldb.OpenFile(pathToFiles+"appdata.db", nil)
	if err != nil {
		return err
	}
	edges, err := leveldb.OpenFile(pathToFiles+"datarelation.db", nil)
	if err != nil {
		return err
	}
	timestamps, err := leveldb.OpenFile(pathToFiles+"datatimestamps.db", nil)
	if err != nil {
		return err
	}
	// Set the databases
	db.nodeDB = nodes
	db.edgeDB = edges
	db.timestampDB = timestamps

	return nil
}

func (db *LevelDbImpl) Close() {
	db.nodeDB.Close()
	db.edgeDB.Close()
	db.timestampDB.Close()
}
