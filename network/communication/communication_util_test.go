package communication_test

import (
	"context"
	"dforum-app/network/communication"
	"dforum-app/storage"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	mplex "github.com/libp2p/go-libp2p-mplex"
	libp2ptls "github.com/libp2p/go-libp2p-tls"
	"github.com/libp2p/go-tcp-transport"
	"github.com/multiformats/go-multiaddr"
)

func createAndInitCommMgr(port int, pathToDbs string) (*communication.CommunicationManager, string, *storage.StorageModule) {
	ctx := context.Background()
	priv, _, _ := crypto.GenerateKeyPair(
		crypto.Ed25519, // Select your key type. Ed25519 are nice short
		-1,             // Select key length when possible (i.e. RSA).
	)
	identity := libp2p.Identity(priv)
	listeningAddresses := libp2p.ListenAddrStrings(
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port), // regular tcp connections
		fmt.Sprintf("/ip6/::/tcp/%d", port),      // include IPv6 support
	)
	transport := libp2p.Transport(tcp.NewTCPTransport)
	muxers := libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport)
	security := libp2p.Security(libp2ptls.ID, libp2ptls.New)

	h, err := libp2p.New(
		identity,
		listeningAddresses,
		transport,
		muxers,
		security,
	)
	if err != nil {
		panic(err)
	}
	sM := storage.NewStorageModule(pathToDbs)
	cM := communication.NewCommunicationManager(sM)
	h.SetStreamHandler(cM.GetProtocolID(), cM.GetMessageHandler())
	cM.SetHost(h, ctx)
	peerInfo := peer.AddrInfo{
		ID:    h.ID(),
		Addrs: h.Addrs(),
	}
	addrs, _ := peer.AddrInfoToP2pAddrs(&peerInfo)
	return cM, addrs[0].String(), sM
}

func connectNodes(cM *communication.CommunicationManager, address string) {
	host, ctx := cM.GetHost()
	seedAddr, err := multiaddr.NewMultiaddr(address)
	if err != nil {
		panic(err)
	}
	seedInfo, err := peer.AddrInfoFromP2pAddr(seedAddr)
	if err != nil {
		panic(err)
	}
	err = host.Connect(ctx, *seedInfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connecting to bootstrap: %s", err)
	} else {
		fmt.Println("Connected to", seedInfo.ID)
	}
}
