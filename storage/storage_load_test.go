package storage_test

import (
	"crypto/sha256"
	"dforum-app/security"
	"dforum-app/storage"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

func BenchmarkStoringNodes(b *testing.B) {
	os.RemoveAll("../test")
	sut := storage.NewStorageModule("../test/")
	defer sut.TearDown()

	for i := 0; i < 5; i++ {
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sut.StoreNode(newTestNodeForBenchmark(int64(i)))
			}
		})
	}
}

// The purpose of this test is to analyse how much space storing
// x amount of nodes would take on a user's system.
func TestStorageSize(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping testing database size due to short test...")
	}
	// Setup
	os.RemoveAll("../test")
	sut := storage.NewStorageModule("../test/")
	defer sut.TearDown()
	// Run 1 hundred thousand times
	for i := 0; i < 100000; i++ {
		sut.StoreNode(newRealisticTestNode())
	}
	// File size then needs to be checked manually
	// 5,449,872 bytes for 100,000 dummy nodes = 55 bytes / node
	// 57,389,056 bytes for 100,000 realistic nodes = 574 bytes / node
}

func newRealisticTestNode() *storage.Node {
	topic, _ := randomString(60)
	content, _ := randomString(150)
	dataObj := storage.DataObject{
		Topic:     topic,
		Content:   content,
		Parent:    sha256.Sum224([]byte("Hello World!")),
		Timestamp: time.Now().Unix(),
		Indicator: 5,
	}
	return &storage.Node{
		DatObj: dataObj,
		SecObj: security.SecurityObject{
			Fingerprint: sha256.Sum224(dataObj.GetBytes()),
			ProofOfWork: "DF1:16:3Nrb2j5bHRk=:MTM3Nw==",
		},
	}
}

func newTestNodeForBenchmark(id int64) *storage.Node {
	fingerprint := make([]byte, 8)
	binary.BigEndian.PutUint64(fingerprint, uint64(id))

	return &storage.Node{
		DatObj: storage.DataObject{
			Topic:     "A new topic!",
			Content:   "This is a very generic content of short length",
			Parent:    [28]byte{},
			Timestamp: 1649000510,
			Indicator: 5,
		},
		SecObj: security.SecurityObject{
			Fingerprint: [28]byte{fingerprint[0], fingerprint[1], fingerprint[2], fingerprint[3]},
			ProofOfWork: "DF1:16:random:1500",
		},
	}
}

func randomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
