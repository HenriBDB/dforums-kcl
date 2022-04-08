package communication_test

import (
	"bufio"
	"context"
	"crypto/rand"
	"dforum-app/network/communication"
	"io/ioutil"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/network"
)

func TestHangingStreamReceiver(t *testing.T) {
	// Attempt to create a hanging stream connection that could lead to a DoS attack
	if testing.Short() {
		t.Skip()
	}
	cM1, _, _ := createAndInitCommMgr(7068, "../../test/test1/")
	cM2, addr2, _ := createAndInitCommMgr(8069, "../../test/test2/")
	connectNodes(cM1, addr2)
	defer cM1.TearDown()
	defer cM2.TearDown()

	h1, ctx := cM1.GetHost()
	s, err := h1.NewStream(ctx, h1.Network().Peers()[0], communication.MessageProtocol)
	if err != nil {
		t.Fatal(err.Error())
	}
	w := bufio.NewWriter(s)
	msg := []byte("of")
	w.Write(msg)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	testSuccessful := make(chan bool, 1)
	go func() {
		// Wait for stream to close
		ioutil.ReadAll(s)
		testSuccessful <- true
	}()

	select {
	case <-ctx.Done():
		t.Fatal("stream is hanging for over 5 seconds, Denial of Service attacks are possible")
	case <-testSuccessful:
	}
}

func TestHangingStreamSender(t *testing.T) {
	// Attempt to create a hanging stream connection that could lead to a DoS attack
	if testing.Short() {
		t.Skip()
	}
	cM1, _, _ := createAndInitCommMgr(7068, "../../test/test1/")
	cM2, addr2, _ := createAndInitCommMgr(8069, "../../test/test2/")
	connectNodes(cM1, addr2)
	defer cM1.TearDown()
	defer cM2.TearDown()

	h1, _ := cM1.GetHost() // Set handler that will ignore requests and let streams hang
	h1.SetStreamHandler(communication.MessageProtocol, func(s network.Stream) { return })
	h2, _ := cM2.GetHost()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	testSuccessful := make(chan bool, 1)
	go func() {
		cM2.SendSyncRequest(h2.Network().Peers()[0])
		testSuccessful <- true
	}()

	select {
	case <-ctx.Done():
		t.Fatal("stream is hanging for over 5 seconds, Denial of Service attacks are possible")
	case <-testSuccessful:
	}
}

func TestStreamBufferOverload(t *testing.T) {
	// Attempt to create a hanging stream connection that could lead to a DoS attack
	if testing.Short() {
		t.Skip()
	}
	cM1, _, _ := createAndInitCommMgr(7068, "../../test/test1/")
	cM2, addr2, _ := createAndInitCommMgr(8069, "../../test/test2/")
	connectNodes(cM1, addr2)
	defer cM1.TearDown()
	defer cM2.TearDown()

	h1, _ := cM1.GetHost() // Set handler that will ignore requests and let streams hang
	h1.SetStreamHandler(communication.MessageProtocol, func(s network.Stream) {
		// Spam peer with random data
		w := bufio.NewWriter(s)
		for i := 0; i < 5000000; i++ {
			b := make([]byte, 1000)
			rand.Read(b)
			w.Write(b)
			w.Flush()
		}
	})
	h2, _ := cM2.GetHost()

	cM2.SendSyncRequest(h2.Network().Peers()[0])
}
