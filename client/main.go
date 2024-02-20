package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/websocket"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	helper "github.com/mmouru/p2p-ml/helper"
)

const ftp = "/file-transfer/1.0.0"

const mdlp = "/model-transfer/1.0.0"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:3000"
	},
}

//var h host.Host

func sendModelProtocol(s network.Stream) error {
	log.Printf("Start receiving model parameters from %s\n", s.Conn().RemotePeer())
	gzr, err := gzip.NewReader(s)
	if err != nil {
		fmt.Println("Error creating gzip reader:", err)
		return err
	}
	defer gzr.Close()

	os.MkdirAll("weights", 0755)

	modelpth, _ := io.ReadAll(gzr) // virheen käsittely olisi hyväksi

	err = os.WriteFile(fmt.Sprintf("weights/model_%s.pth", s.Conn().RemotePeer()), modelpth, 0666)
	if err != nil {
		log.Printf("something went wrong saving model parameters")
		panic(err)
	}

	return nil
}

func fileTransferProtocol(s network.Stream, h host.Host) (string, error) {
	log.Printf("New incoming connection from peer %s\n", s.Conn().RemotePeer())
	gzr, err := gzip.NewReader(s)
	if err != nil {
		fmt.Println("Error creating gzip reader:", err)
		return "", err
	}
	defer gzr.Close()

	decompressedData, err := io.ReadAll(gzr)
	if err != nil {
		fmt.Println("Error reading from gzip reader:", err)
		return "", err
	}

	//fmt.Println(string(decompressedData))
	zipFileWriteName := "train_set2.zip"

	err = os.WriteFile(zipFileWriteName, decompressedData, 0666)
	if err != nil {
		fmt.Println("Error reading from gzip reader:", err)
		return "", err
	}
	fmt.Println("ready with receiving data")

	// unzip the training set
	helper.UnzipTrainingSet(zipFileWriteName, "data")

	defer os.RemoveAll("data")
	defer os.RemoveAll(zipFileWriteName)

	// logic to run the training on current computer
	cmd := exec.Command("python3", "trainer.py")

	err = cmd.Run()

	if err != nil {
		log.Println(err)
		os.RemoveAll("data")
		os.RemoveAll(zipFileWriteName)
	}
	fmt.Println("Finished with training the model")

	// if no err in training return the model parameters to the requester
	model, err := os.Open("./data/model_state_dict.pth")

	if err != nil {
		log.Printf("JOO EI VOINU AVATA MODEL PTHs")
	}

	s2, err := h.NewStream(context.Background(), s.Conn().RemotePeer(), mdlp)

	if err != nil {
		log.Printf("Joo ei pygeny avaa uuttaa conenction")
	}

	buffer := make([]byte, 1024)
	gzw := gzip.NewWriter(s2)

	for {
		n, err := model.Read(buffer)
		if err == io.EOF {
			break // End of file
		}
		if err != nil {
			log.Fatal(err)
		}

		// Write the chunk to the gzip writer
		_, err = gzw.Write(buffer[:n])
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("ready with sending data")
	s2.Close()
	// then return the saved weights or something???

	//log.Printf("Message from %s: %s", connection.RemotePeer().String(), message)
	return "ready", nil
}

func main() {
	sp := flag.String("sp", "3001", "Source port for local host")
	//peerAddr := flag.String("peer-address", "", "peer address")
	flag.Parse()

	h, err := startPeer(*sp) // start listening on port 3001 or specified
	if err != nil {
		return
	}

	go func() {
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			handleWebSocket(w, r, h) // Pass the host.Host variable 'h' here
		})
		if err := http.ListenAndServe(":7070", nil); err != nil {
			log.Fatal("WebSocket server error:", err)
		}
	}()

	h.SetStreamHandler(ftp, func(s network.Stream) {
		log.Printf(("/file-transfer/1.00 stream created"))
		_, err := fileTransferProtocol(s, h)
		if err != nil {
			s.Reset()
		} else {
			s.Close()
		}
	})

	h.SetStreamHandler(mdlp, func(s network.Stream) {
		log.Printf(("/model-transfer/1.00 stream created"))
		err := sendModelProtocol(s)
		if err != nil {
			s.Reset()
		} else {
			s.Close()
		}
	})

	defer h.Close()
	//defer joq.jotain() set Peer Disconnecting

	/*if *peerAddr != "" {
		connectPeer(h, *peerAddr)
	}*/

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

	respCode := helper.RegisterPeerToCentralList(sourcePort, h.ID().String())

	if respCode != 200 {
		log.Fatalf("something went wrong with resp code: %d", respCode)
	}

	fmt.Println("ConnectAdress:", fmt.Sprintf("%s/p2p/%s", h.Addrs()[0], h.ID()))
	fmt.Println("Public Ip Connection:", fmt.Sprintf("/ip4/%s/tcp/%s/p2p/%s", helper.GetPublicIp(), sourcePort, h.ID()))

	return h, nil

}

// setup connection to peer
func connectPeer(h host.Host, peerIp string, peerPort string, peerId string) {
	// build the addr string from pieces
	peerAddr := fmt.Sprintf("/ip4/%s/tcp/%s/p2p/%s", peerIp, peerPort, peerId)
	// create the peer connection
	peerMA, err := multiaddr.NewMultiaddr(peerAddr)

	if err != nil {
		panic(err)
	}
	peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
	if err != nil {
		panic(err)
	}

	fmt.Println(peerAddrInfo, "addr info")

	// Connect to the node at the given address.
	if err := h.Connect(context.Background(), *peerAddrInfo); err != nil {
		panic(err)
	}
	fmt.Println("Connected to", peerAddrInfo.String())
	s, err := h.NewStream(context.Background(), peerAddrInfo.ID, ftp)
	if err != nil {
		panic(err)
	}

	//f, _ := os.Open("./1.png")
	zip, _ := os.Open("./train_set.zip")
	defer zip.Close()

	buffer := make([]byte, 1024)
	gzw := gzip.NewWriter(s)
	for {
		n, err := zip.Read(buffer)
		if err == io.EOF {
			break // End of file
		}
		if err != nil {
			log.Fatal(err)
		}

		// Write the chunk to the gzip writer
		_, err = gzw.Write(buffer[:n])
		if err != nil {
			log.Fatal(err)
		}
	}

	//content, _ := io.ReadAll(f)
	//

	//_, err = gzw.Write(content)
	if err != nil {
		fmt.Println("Error writing to gzip writer:", err)
		return
	}
	gzw.Close()
	s.Close()
	select {}
}

type PeerInfo struct {
	Ip     string `json:"ip"`
	Port   string `json:"port"`
	PeerId string `json:"peer_id"`
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, h host.Host) {
	// Upgrade HTTP connection to WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Error upgrading to WebSocket", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Handle WebSocket messages
	for {
		_, message, err := conn.ReadMessage() // messageType ekaan jos haluaa palauttaa viestin
		if err != nil {
			// Handle error or connection closure
			return
		}
		fmt.Println(string(message))
		// Process received message (e.g., echo back)
		var peersToConnect []PeerInfo
		// unmarshal the stringified JSON
		if err := json.Unmarshal([]byte(string(message)), &peersToConnect); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		for _, peer := range peersToConnect {
			fmt.Println(peer.Port)
			go connectPeer(h, "87.100.213.208", "3002", "12D3KooWNBgdPvEAm7VrmLn9ciCpVT4NAQvYsc5PaNvrz2Ax5Mjb")
		}

		//connectPeer()
		if err != nil {
			// Handle error
			return
		}
	}
}
