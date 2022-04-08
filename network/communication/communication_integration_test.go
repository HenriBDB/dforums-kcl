package communication_test

import (
	"dforum-app/storage"
	"fmt"
	"testing"
	"time"
)

func TestInventoryMessage(*testing.T) {
	cM1, _, sM1 := createAndInitCommMgr(7068, "../../test/test1/")
	cM2, addr2, sM2 := createAndInitCommMgr(7069, "../../test/test2/")
	connectNodes(cM1, addr2)
	defer cM1.TearDown()
	defer cM2.TearDown()

	root := storage.NewNode("Topic", "detail", 5, [28]byte{})
	sM1.StoreNode(root)

	fmt.Println("Expected:", root)
	cM1.SendInventoryMessage(root.GetFingerprint())

	time.Sleep(1 * time.Second)
	hashcode := root.GetFingerprint()
	actualNode := sM2.GetNode(hashcode, true)
	fmt.Println("Actual:", actualNode)
}

func TestSendSyncRequest(*testing.T) {
	cM1, _, _ := createAndInitCommMgr(7068, "../../test/test1/")
	cM2, addr2, _ := createAndInitCommMgr(8069, "../../test/test2/")
	connectNodes(cM1, addr2)
	defer cM1.TearDown()
	defer cM2.TearDown()

	h1, _ := cM1.GetHost()
	cM1.SendSyncRequest(h1.Network().Peers()[0])
	time.Sleep(time.Second)
}
