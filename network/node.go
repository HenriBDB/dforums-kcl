package network

import (
	"context"
	"dforum-app/configuration"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	mplex "github.com/libp2p/go-libp2p-mplex"
	libp2ptls "github.com/libp2p/go-libp2p-tls"
	"github.com/libp2p/go-tcp-transport"
)

/*
	Node Factories
*/

func CreateDefaultNode(port int, priv crypto.PrivKey) (host.Host, context.Context, *dht.IpfsDHT) {
	// Inspired by https://github.com/libp2p/go-libp2p/blob/master/examples/libp2p-host/host.go
	// and https://github.com/libp2p/go-libp2p/blob/master/examples/ipfs-camp-2019
	ctx := context.Background()

	// Create Node
	identity := libp2p.Identity(priv)

	listeningAddresses := libp2p.ListenAddrStrings(
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port), // regular tcp connections
		fmt.Sprintf("/ip6/::/tcp/%d", port),      // include IPv6 support
	)
	transport := libp2p.Transport(tcp.NewTCPTransport)

	muxers := libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport)
	security := libp2p.Security(libp2ptls.ID, libp2ptls.New)

	low, high := configuration.GetConnectionLimits()
	cm, _ := connmgr.NewConnManager(
		low,                                  // Lowwater
		high,                                 // HighWater
		connmgr.WithGracePeriod(time.Minute), // GracePeriod
	)
	connectionManager := libp2p.ConnectionManager(cm)

	// Let this host use the DHT to find other hosts
	var idht *dht.IpfsDHT
	dhtRouting := func(h host.Host) (routing.PeerRouting, error) {
		var err error
		idht, err = dht.New(ctx, h)
		return idht, err
	}
	routing := libp2p.Routing(dhtRouting)

	h, err := libp2p.New(
		identity,
		listeningAddresses,
		transport,
		muxers,
		security,
		connectionManager,
		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),
		routing,
	)
	if err != nil {
		panic(err)
	}

	peerInfo := peer.AddrInfo{
		ID:    h.ID(),
		Addrs: h.Addrs(),
	}
	addrs, _ := peer.AddrInfoToP2pAddrs(&peerInfo)
	configuration.Logger.Info("Listening in on:", addrs[0].String())

	return h, ctx, idht
}
