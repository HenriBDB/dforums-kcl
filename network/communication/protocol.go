package communication

import (
	"bufio"
	"context"
	"dforum-app/configuration"
	"dforum-app/security"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

const MessageProtocol = protocol.ID("/libp2p/DDF/0.0.1")

// Inspired by: https://github.com/aethereans/aether-app
type ProtocolAction uint

func (a ProtocolAction) String() string {
	actionCode := [...]string{
		"InvalidMessage",
		"SyncRequest",
		"InventoryMessage",
		"DataRequest",
		// This set has to match the set in const() and its order.
	}
	if !a.isValid() {
		return "Invalid Network Action."
	}
	return actionCode[a]
}

func (a ProtocolAction) isValid() bool {
	return InvalidMessage <= a && a <= DataRequest
}

// Available actions matching action codes above
const (
	InvalidMessage ProtocolAction = iota
	SyncRequest
	InventoryMessage
	DataRequest
)

func parseActionByte(actionCode byte) ProtocolAction {
	action := ProtocolAction(uint8(actionCode))
	if !action.isValid() {
		return InvalidMessage
	}
	return action
}

/*
	HELPER METHODS
*/

func messageHandler(cm *CommunicationManager, s network.Stream) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	defer s.Close()
	configuration.Logger.Info(s.ID(), "- new stream from:", s.Conn().RemotePeer().ShortString())
	metaSize := 3
	meta := make([]byte, metaSize)
	err := readFromStreamWithContext(ctx, meta, s, metaSize)
	if err != nil {
		configuration.Logger.Error(s.ID(), "failed to handle message header:", err.Error())
		return
	}
	action := parseActionByte(meta[0])
	if action == InvalidMessage {
		configuration.Logger.Info(s.ID(), "peer sent an InvalidMessage message")
		return
	}
	contentSize := binary.BigEndian.Uint16(meta[1:])
	// The following actions require data
	if contentSize <= 0 {
		sendInvalidMessage(s)
		return
	}
	content := make([]byte, contentSize)
	readFromStreamWithContext(ctx, content, s, int(contentSize))
	if err != nil {
		configuration.Logger.Error(s.ID(), "failed to handle message body:", err.Error())
		return
	}
	switch action {
	case SyncRequest:
		cm.handleSyncRequest(content, s)
	case InventoryMessage:
		cm.handleInventoryMessage(content, s)
	case DataRequest:
		cm.handleDataRequest(content, s)
	default:
		configuration.Logger.Error(s.ID(), "message received did not conform to the communication protocol")
		return
	}
}

func readFromStreamWithContext(ctx context.Context, buffer []byte, s network.Stream, size int) error {
	readDone := make(chan error, 1)
	go func() {
		n, err := io.ReadFull(s, buffer)
		if n != size {
			readDone <- errors.New("invalid message length received")
			return
		}
		readDone <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-readDone:
		return err
	}
}

func readAllFromStreamWithContext(ctx context.Context, s network.Stream) ([]byte, error) {
	readDone := make(chan error, 1)
	var data []byte
	go func() {
		var err error
		data, err = io.ReadAll(s)
		readDone <- err
	}()

	select {
	case <-ctx.Done():
		return data, ctx.Err()
	case err := <-readDone:
		return data, err
	}
}

func simpleSend(msg []byte, s network.Stream) error {
	w := bufio.NewWriter(s)
	n, err := w.Write(msg)
	if n != len(msg) {
		return fmt.Errorf("expected to write %d bytes, wrote %d", len(msg), n)
	}
	if err != nil {
		return err
	}
	if err = w.Flush(); err != nil {
		return err
	}
	return nil
}

func sendInvalidMessage(s network.Stream) error {
	msg := []byte{byte(InvalidMessage)}
	return simpleSend(msg, s)
}

// https://github.com/libp2p/go-libp2p/blob/master/examples/chat-with-mdns/main.go
func sendRequestWithResponse(msg []byte, s network.Stream) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	defer s.Close()

	err := simpleSend(msg, s)
	if err != nil {
		return nil, err
	}
	data, err := readAllFromStreamWithContext(ctx, s)
	if err != nil {
		return nil, err
	}
	if len(data) > 0 {
		return data, nil
	}
	return nil, nil
}

func broadcastToAllPeers(msg []byte, host host.Host, ctx context.Context) {
	for _, peer := range host.Network().Peers() {
		if _, err := host.Peerstore().SupportsProtocols(peer, string(MessageProtocol)); err == nil {
			s, err := getPeerStream(peer, host, ctx)
			if err != nil {
				continue
			}
			err = simpleSend(msg, s)
			if err != nil {
				configuration.Logger.Error(s.ID(), "failed to broadcast message to peer:", err.Error())
			}
			s.Close()
		}
	}
}

func getPeerStream(peer peer.ID, host host.Host, ctx context.Context) (network.Stream, error) {
	_, err := host.Peerstore().SupportsProtocols(peer, string(MessageProtocol))
	if err == nil {
		s, err := host.NewStream(ctx, peer, MessageProtocol)
		if err == nil {
			configuration.Logger.Info(s.ID(), "connect to peer:", peer.ShortString())
		}
		return s, err
	}
	return nil, err
}

func getHashSignatureFromMessage(msg []byte) (security.HashSignature, error) {
	if len(msg) != 28 { // Signatures are 28 bytes in length, include integration test for this
		return security.HashSignature{}, errors.New("invalid msg length provided")
	}
	dataID := [28]byte{}
	copy(dataID[:], msg[0:28])
	return dataID, nil
}

/*
	MESSAGE BUILDERS
*/

func buildSimpleActionMessage(action ProtocolAction, size uint16) []byte {
	byteSize := make([]byte, 2)
	binary.BigEndian.PutUint16(byteSize, size)
	return []byte{byte(action), byteSize[0], byteSize[1]}
}

func BuildSyncRequest(t time.Time) []byte {
	epochTime, _ := json.Marshal(t.Unix())
	return append(buildSimpleActionMessage(SyncRequest, uint16(len(epochTime))), epochTime...)
}

func BuildInventoryMessage(id security.HashSignature) []byte {
	return append(buildSimpleActionMessage(InventoryMessage, uint16(len(id))), (id[:])...)
}

func BuildDataRequestMsg(id security.HashSignature) []byte {
	return append(buildSimpleActionMessage(DataRequest, uint16(len(id))), (id[:])...)
}
