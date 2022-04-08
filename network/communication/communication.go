package communication

import (
	"context"
	"dforum-app/configuration"
	"dforum-app/security"
	"dforum-app/storage"
	"encoding/json"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

type CommunicationManager struct {
	ctx              context.Context
	host             host.Host
	localStorage     *storage.StorageModule
	inventoryHandler InventoryHandler
}

func NewCommunicationManager(sm *storage.StorageModule) *CommunicationManager {
	cM := &CommunicationManager{
		localStorage:     sm,
		inventoryHandler: NewInventoryHandler(),
	}
	sm.Subscribe(cM)
	return cM
}

func (cm *CommunicationManager) TearDown() {
	// Store addresses of current peers before shutting down
	currentPeerAddresses := []string{}
	for _, peerID := range cm.host.Network().Peers() {
		peerInfo := peer.AddrInfo{
			ID:    peerID,
			Addrs: cm.host.Network().Peerstore().Addrs(peerID),
		}
		addrs, _ := peer.AddrInfoToP2pAddrs(&peerInfo)
		currentPeerAddresses = append(currentPeerAddresses,
			addrs[0].String())
	}
	configuration.SetNetworkPeers(currentPeerAddresses)
	cm.host.Close()
}

// Return a message handler specific for a communication manager
func (cm *CommunicationManager) GetMessageHandler() func(network.Stream) {
	return func(s network.Stream) {
		messageHandler(cm, s)
	}
}

func (cm *CommunicationManager) GetProtocolID() protocol.ID {
	return MessageProtocol
}

func (cm *CommunicationManager) SetHost(h host.Host, ctx context.Context) {
	cm.host = h
	cm.ctx = ctx
}

func (cm *CommunicationManager) GetHost() (h host.Host, ctx context.Context) {
	return cm.host, cm.ctx
}

func (cm *CommunicationManager) Sync(p peer.ID) {
	cm.SendSyncRequest(p)
}

func (cm *CommunicationManager) RegisterNewNode(n *storage.Node) {
	cm.SendInventoryMessage(n.GetFingerprint())
}

/*
	Communication Actions
*/

func (cm *CommunicationManager) handleInventoryMessage(msg []byte, s network.Stream) {
	// Parse
	id, err := getHashSignatureFromMessage(msg)
	if err != nil {
		configuration.Logger.Error(s.ID(), "received invalid inventory message")
		return
	}
	configuration.Logger.Info(s.ID(), "receieved inventory message:", id[0:4])
	cm.registerNodeInv(id, s.Conn().RemotePeer())
}

// Shares locally created posts with other users
// Broadcasts the post's identity to all peers
// Peers can then request the full post in a second request
func (cm *CommunicationManager) SendInventoryMessage(id security.HashSignature) {
	configuration.Logger.Info("sending inventory message:", id[0:4])
	msg := BuildInventoryMessage(id)
	broadcastToAllPeers(msg, cm.host, cm.ctx)
}

func (cm *CommunicationManager) handleDataRequest(msg []byte, s network.Stream) {
	// Parse
	id, err := getHashSignatureFromMessage(msg)
	if err != nil {
		configuration.Logger.Error(s.ID(), "received invalid data request")
		sendInvalidMessage(s)
		return
	}
	configuration.Logger.Info(s.ID(), "received data request:", id[0:4])
	// Respond
	node := cm.localStorage.GetNode(id, true)
	if node == nil {
		sendInvalidMessage(s)
		return
	}
	err = simpleSend(node.GetBytes(), s)
	if err != nil {
		configuration.Logger.Error(s.ID(), "failed to respond to data request:", err.Error())
	}
}

func (cm *CommunicationManager) SendDataRequest(id security.HashSignature, peer peer.ID) bool {
	for i := 0; i < 5; i++ { // try 5 times
		s, err := getPeerStream(peer, cm.host, cm.ctx)
		if err != nil {
			configuration.Logger.Error("failed to get stream for data request from peer:", peer.ShortString(), err.Error())
			continue
		}
		configuration.Logger.Info(s.ID(), "sending data request:", id[0:4])
		msg := BuildDataRequestMsg(id)
		response, err := sendRequestWithResponse(msg, s)
		if err != nil {
			configuration.Logger.Errorf("%s failed to complete data request %d/5: %s", s.ID(), i, err.Error())
			continue
		}
		node := storage.ParseNode(response)
		if node == nil { // Failed to receive a valid node
			configuration.Logger.Error(s.ID(), "invalid node received")
			continue
		}
		if ok := node.Verify(); ok {
			cm.localStorage.StoreNode(node)
			cm.localStorage.PublishNode(node)
			return true
		} else {
			configuration.Logger.Error(s.ID(), "node received did not meet security verifications")
		}
		return false
	}
	return false
}

func (cm *CommunicationManager) handleSyncRequest(msg []byte, s network.Stream) {
	// Decipher request
	var unixTime int64
	json.Unmarshal(msg, &unixTime)
	rTime := time.Unix(unixTime, 0)
	// Time provided has to be in the past
	if time.Now().Before(rTime) {
		configuration.Logger.Error(s.ID(), "received invalid sync request")
		sendInvalidMessage(s)
		return
	}
	//Limit catch-up to two weeks
	twoWeeksAgo := time.Now().AddDate(0, 0, -14)
	if rTime.Before(twoWeeksAgo) {
		rTime = twoWeeksAgo
	}
	configuration.Logger.Info(s.ID(), "received sync request for date:", rTime.Format(time.RFC822Z))
	// Send Inv Messages
	dataItems := cm.localStorage.GetNodesSince(rTime)
	jsonItems, _ := json.Marshal(dataItems)
	err := simpleSend(jsonItems, s)
	if err != nil {
		configuration.Logger.Error(s.ID(), "failed to respond to sync request:", err.Error())
	}
}

func (cm *CommunicationManager) SendSyncRequest(peer peer.ID) {
	t := cm.localStorage.TimeOfMostRecentNode()
	msg := BuildSyncRequest(t)
	s, err := getPeerStream(peer, cm.host, cm.ctx)
	if err != nil {
		configuration.Logger.Error("failed to get stream for sync request from peer:", peer.ShortString(), err.Error())
		return
	}
	configuration.Logger.Info(s.ID(), "sending sync request for date:", t.Format(time.RFC822Z))
	data, err := sendRequestWithResponse(msg, s)
	if err != nil {
		configuration.Logger.Error(s.ID(), "failed to complete sync request:", err.Error())
		return
	}
	var dataItems []security.HashSignature
	json.Unmarshal(data, &dataItems)
	for _, v := range dataItems {
		cm.registerNodeInv(v, peer)
	}
}

/*
	Helper Methods
*/

func (cm *CommunicationManager) registerNodeInv(id security.HashSignature, peer peer.ID) {
	if cm.localStorage.NodeExists(id) {
		return // Ignore inv if node already exists
	}
	lock := cm.inventoryHandler.GetLock(id)
	lock.CompleteActionUntilSuccessful(func() bool {
		if cm.SendDataRequest(id, peer) {
			// Node stored successfully
			cm.inventoryHandler.DeleteLock(id)
			return true
		}
		return false
	})
}
