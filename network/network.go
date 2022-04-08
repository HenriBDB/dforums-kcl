package network

import (
	"context"
	"dforum-app/configuration"
	"dforum-app/network/communication"
	"dforum-app/storage"
	"math/rand"
	"sync"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	disc "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
)

type NetworkModule struct {
	communicationMgr *communication.CommunicationManager
}

func NewNetworkModule(sM *storage.StorageModule) *NetworkModule {
	return &NetworkModule{
		communicationMgr: communication.NewCommunicationManager(sM),
	}
}

func (n *NetworkModule) CreateAndStartHost() {
	//TODO Security store private keys and load them dynamically
	priv, _, err := crypto.GenerateKeyPair(
		crypto.RSA, // Select your key type. Ed25519 are nice short
		2048,       // Select key length when possible (i.e. RSA).
	)
	if err != nil {
		panic(err)
	}

	host, ctx, dht := CreateDefaultNode(configuration.GetNetworkPort(), priv)
	n.communicationMgr.SetHost(host, ctx)

	host.SetStreamHandler(n.communicationMgr.GetProtocolID(), n.communicationMgr.GetMessageHandler())

	bootstrap(host, ctx)
	setPeerRouting(host, ctx, dht, n.communicationMgr.GetProtocolID())

	// Sync with random peer, try max 10 peers
	if len(host.Network().Peers()) > 0 {
		for i := 0; i < 10; i++ {
			randIdx := rand.Intn(len(host.Network().Peers()))
			randPeer := host.Network().Peers()[randIdx]
			if _, err := host.Peerstore().SupportsProtocols(randPeer, string(n.communicationMgr.GetProtocolID())); err == nil {
				n.communicationMgr.Sync(randPeer)
				break
			}
		}
	}
}

func bootstrap(h host.Host, ctx context.Context) {
	// Boostrap onto the network
	targetPeers := append(configuration.GetNetworkSeeds(), configuration.GetNetworkPeers()...)

	wg := new(sync.WaitGroup)
	wg.Add(len(targetPeers))
	for _, address := range targetPeers {
		go func(address string) {
			defer wg.Done()
			seedAddr, err := multiaddr.NewMultiaddr(address)
			if err != nil {
				configuration.Logger.Errorf("connecting to bootstrap: %s", err)
			}
			seedInfo, err := peer.AddrInfoFromP2pAddr(seedAddr)
			if err != nil {
				configuration.Logger.Errorf("connecting to bootstrap: %s", err)
			}
			err = h.Connect(ctx, *seedInfo)
			if err != nil {
				configuration.Logger.Errorf("connecting to bootstrap: %s", err)
			} else {
				configuration.Logger.Info("connected to", seedInfo.ID.ShortString())
			}
		}(address)
	}
	go func() {
		wg.Wait()
		// All connection attemps to peers failed
		if len(h.Network().Peers()) == 0 {
			configuration.Logger.Error("could not connect to any configured peers")
		}
	}()
}

func (n *NetworkModule) GetAddress() multiaddr.Multiaddr {
	h, _ := n.communicationMgr.GetHost()
	peerInfo := peer.AddrInfo{
		ID:    h.ID(),
		Addrs: h.Addrs(),
	}
	addrs, _ := peer.AddrInfoToP2pAddrs(&peerInfo)
	return addrs[0]
}

func (n *NetworkModule) TearDown() {
	n.communicationMgr.TearDown()
}

func setPeerRouting(h host.Host, ctx context.Context, idht *dht.IpfsDHT, messageProtocol protocol.ID) {
	err := idht.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}
	// Set discovery and find peers
	routingDiscovery := disc.NewRoutingDiscovery(idht)
	disc.Advertise(ctx, routingDiscovery, string(messageProtocol))
	peers, err := disc.FindPeers(ctx, routingDiscovery, string(messageProtocol))
	if err != nil {
		panic(err)
	}
	for _, peer := range peers {
		if h.Network().Connectedness(peer.ID) != network.Connected {
			configuration.Logger.Infof("Found peer: %s", peer.ID.ShortString())
			h.Connect(ctx, peer)
		}
	}
}
