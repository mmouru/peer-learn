package main

import (
	"compress/gzip"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

const helloProtocol = "/hello/1.0.0"

func readHelloProtocol(s network.Stream) error {
	log.Printf("New incoming connection from peer %s\n", s.Conn().RemotePeer())
	gzr, err := gzip.NewReader(s)
	if err != nil {
		fmt.Println("Error creating gzip reader:", err)
		return err
	}
	defer gzr.Close()

	decompressedData, err := io.ReadAll(gzr)
	if err != nil {
		fmt.Println("Error reading from gzip reader:", err)
		return err
	}

	fmt.Println(string(decompressedData))

	err = os.WriteFile("moro.txt", decompressedData, 0666)
	if err != nil {
		fmt.Println("Error reading from gzip reader:", err)
		return err
	}

	//connection := s.Conn()

	//log.Printf("Message from %s: %s", connection.RemotePeer().String(), message)
	return nil
}

func main() {
	sp := flag.String("sp", "3001", "Source port for local host")
	peerAddr := flag.String("peer-address", "", "peer address")
	flag.Parse()

	h, err := startPeer(*sp) // start listening on port 3001
	if err != nil {
		return
	}

	h.SetStreamHandler(helloProtocol, func(s network.Stream) {
		log.Printf(("/hello/1.00 stream created"))
		err := readHelloProtocol(s)
		if err != nil {
			s.Reset()
		} else {
			s.Close()
		}
	})

	defer h.Close()

	if *peerAddr != "" {
		connectPeer(h, *peerAddr)
	}

	select {} // run indefinetly

}

func startPeer(sourcePort string) (host.Host, error) {

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
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", sourcePort), // regular tcp connections
		),
	)

	fmt.Println("ConnectAdress:", fmt.Sprintf("%s/p2p/%s", h.Addrs()[0], h.ID()))

	return h, nil

}

// setup connection to peer
func connectPeer(h host.Host, peerAddr string) {
	peerMA, err := multiaddr.NewMultiaddr(peerAddr)

	if err != nil {
		panic(err)
	}
	peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
	if err != nil {
		panic(err)
	}

	fmt.Println(peerAddrInfo)

	// Connect to the node at the given address.
	if err := h.Connect(context.Background(), *peerAddrInfo); err != nil {
		panic(err)
	}
	fmt.Println("Connected to", peerAddrInfo.String())
	s, err := h.NewStream(context.Background(), peerAddrInfo.ID, "/hello/1.0.0")
	if err != nil {
		panic(err)
	}

	f, _ := os.Open("./test.txt")
	content, _ := io.ReadAll(f)
	gzw := gzip.NewWriter(s)

	_, err = gzw.Write(content)
	if err != nil {
		fmt.Println("Error writing to gzip writer:", err)
		return
	}
	gzw.Close()
	s.Close()
	select {}
}
