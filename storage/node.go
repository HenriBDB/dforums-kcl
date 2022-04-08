package storage

import (
	"dforum-app/configuration"
	"dforum-app/security"
	"encoding/json"
	"time"
)

// Public struct DataObject to be used to represent a node in the network.
type DataObject struct {
	Parent    [28]byte
	Timestamp int64
	Topic     string
	Indicator int8
	Content   string
}

func (do DataObject) GetBytes() []byte {
	res, err := json.Marshal(do)
	if err != nil {
		configuration.Logger.Error("could not convert node to bytes")
		return make([]byte, 0)
	}
	return res
}

type Node struct {
	SecObj security.SecurityObject
	DatObj DataObject
}

func NewNode(topic string, detail string, indicator int8, parentHash [28]byte) *Node {
	do := DataObject{
		Parent:    parentHash,
		Timestamp: time.Now().Unix(),
		Topic:     topic,
		Indicator: indicator,
		Content:   detail,
	}
	security, err := security.GenSecurityObject(do.GetBytes())
	if err != nil {
		configuration.Logger.Errorf("failed to create node with title: %s - %ss", topic, err.Error())
		return nil
	}
	configuration.Logger.Info("successfully created node with title: ", topic)
	return &Node{SecObj: security, DatObj: do}
}

func (n *Node) Verify() bool {
	return n.SecObj.Verify(n.DatObj.GetBytes())
}

func (n Node) GetBytes() []byte {
	res, err := json.Marshal(n)
	if err != nil {
		configuration.Logger.Error("could not convert node to bytes")
		return make([]byte, 0)
	}
	return res
}

func ParseNode(bytes []byte) *Node {
	var node Node
	err := json.Unmarshal(bytes, &node)
	if err != nil {
		configuration.Logger.Error("could not parse node from bytes")
		return nil
	}
	return &node
}

func (n *Node) GetFingerprint() [28]byte {
	return n.SecObj.Fingerprint
}

func (n *Node) GetTimestamp() int64 {
	return n.DatObj.Timestamp
}

func (n Node) String() string {
	res, _ := json.Marshal(n)
	return string(res)
}
