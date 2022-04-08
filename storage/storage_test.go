package storage

import (
	"crypto/sha256"
	"fmt"
	"log"
	"testing"
)

func TestNodeStorage(*testing.T) {
	sut := NewStorageModule("../test/")
	defer sut.TearDown()

	node := NewNode("", "", 5, sha256.Sum224([]byte("Hello World!")))
	fmt.Println("Node size:", len(node.GetBytes()))
	sut.StoreNode(node)

	node2 := sut.GetNode(node.GetFingerprint(), true)
	log.Println(node2.String())
}

func TestParentRelation(*testing.T) {
	// Setup
	sut := NewStorageModule("../test/")
	defer sut.TearDown()

	root := NewNode("Root", "about root", 5, [28]byte{})
	child1 := NewNode("Child 1", "This is the first child node", 4, root.GetFingerprint())
	child2 := NewNode("Child 2", "This is the second child node", 3, root.GetFingerprint())
	sut.StoreNode(root)
	sut.StoreNode(child1)
	sut.StoreNode(child2)

	// Retrieve
	tree := sut.GetChildrenNodes(root.GetFingerprint(), true, 3)
	for _, v := range tree {
		log.Println(v)
	}
}
