package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

const fileTransferProtocol = "/file-transfer/1.0.0"

func main() {
	peerAddr := flag.String("peer-address", "", "peer address")
	h, err := startPeer() // start listening on port 3001
	if err != nil {
		return
	}

	defer h.Close()

	if *peerAddr != "" {
		connectPeer(h, *peerAddr)
	}

	select {} // run indefinetly

}

func startPeer() (host.Host, error) {

	// Set your own keypair
	priv, _, err := crypto.GenerateKeyPair(
		crypto.Ed25519, // Select your key type. Ed25519 are nice short
		-1,             // Select key length when possible (i.e. RSA).
	)
	if err != nil {
		panic(err)
	}

	h, err := libp2p.New(
		// Use the keypair we generated
		libp2p.Identity(priv),
		// Multiple listen addresses
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/3001", // regular tcp connections
		),
	)

	fmt.Println("ConnectAdress:", fmt.Sprintf("%s/p2p/%s", h.Addrs()[0], h.ID()))

	return h, nil

}

// setup connection to peer
func connectPeer(h host.Host, peerAddr string) error {
	peerMA, err := multiaddr.NewMultiaddr(peerAddr)

	if err != nil {
		panic(err)
	}
	peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
	if err != nil {
		panic(err)
	}

	// Connect to the node at the given address.
	if err := h.Connect(context.Background(), *peerAddrInfo); err != nil {
		panic(err)
	}
	fmt.Println("Connected to", peerAddrInfo.String())
	return nil
}
