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
	"sync"

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

var wg sync.WaitGroup

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

	filepath := fmt.Sprintf("weights/model_%s", s.Conn().RemotePeer())

	err = os.WriteFile(fmt.Sprintf("%s.pth", filepath), modelpth, 0666)
	if err != nil {
		log.Printf("something went wrong saving model parameters")
		panic(err)
	}

	defer wg.Done()

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
	zipFileWriteName := "train_set.zip"

	err = os.WriteFile(zipFileWriteName, decompressedData, 0666)
	if err != nil {
		fmt.Println("Error reading from gzip reader:", err)
		return "", err
	}
	fmt.Println("ready with receiving data")

	// unzip the training set
	helper.UnzipFile(zipFileWriteName, "data")

	defer os.RemoveAll("data")
	defer os.RemoveAll(zipFileWriteName)

	fmt.Println("START MODEL TRAINING")
	// logic to run the training on current computer
	cmd := exec.Command("python3", "helper/trainer.py")

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Println(err)
		os.RemoveAll("data")
		os.RemoveAll(zipFileWriteName)
	}
	fmt.Println("Finished with training the model:", string(output))

	// if no err in training return the model parameters to the requester
	pth, err := os.Open("./data/model_state_dict.pth")

	if err != nil {
		log.Printf("JOO EI VOINU AVATA MODEL PTHs")
	}

	s2, err := h.NewStream(context.Background(), s.Conn().RemotePeer(), mdlp)

	if err != nil {
		log.Printf("Joo ei pygeny avaa uuttaa connection")
	}

	buffer := make([]byte, 1024)
	gzw := gzip.NewWriter(s2)

	for {
		n, err := pth.Read(buffer)
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
		http.HandleFunc("/filetransfer", func(w http.ResponseWriter, r *http.Request) {
			filetransferWebSocket(w, r)
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
func connectPeer(h host.Host, peerIp string, peerPort string, peerId string, data_split string) {
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

	zip, _ := os.Open(data_split)
	defer zip.Close()
	// defer os.RemoveAll(data_split)

	buffer := make([]byte, 1024)
	gzw := gzip.NewWriter(s)
	for {
		n, err := zip.Read(buffer)

		if err == io.EOF {
			break // End of file
		}
		if err != nil {
			log.Fatal(err, "kallen orokarinen")
		}

		// Write the chunk to the gzip writer
		_, err = gzw.Write(buffer[:n])
		if err != nil {
			log.Fatal(err, "malle molokarinen")
		}
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
	fmt.Println("MOIRO")
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

		// Split training data between peers.

		n_peers := len(peersToConnect) + 1 // + 1 for self learning
		helper.SplitTrainingDataAmongPeers(n_peers, "./all_training_data")

		// start own learning routine with own split
		wg.Add(1)
		go func() {
			helper.LocalLearningProcess("./splits/split_1.zip")
			defer wg.Done()
		}()

		for i, peer := range peersToConnect {
			wg.Add(1)
			train_split := fmt.Sprintf("./splits/split_%d.zip", i+2)
			go func(peer PeerInfo, trainSplit string) {
				connectPeer(h, peer.Ip, peer.Port, peer.PeerId, trainSplit)
			}(peer, train_split)
		}

		wg.Wait()
		fmt.Println("All trainers have completed.")

		conn.WriteMessage(websocket.TextMessage, []byte("Peer learning completed!"))

		ensemble := exec.Command("python3", "helper/ensemble.py")

		output, err := ensemble.CombinedOutput()

		if err != nil {
			fmt.Println(output)
			panic(err)
		}

		fmt.Println(output)

		// TARJOA ENSEMBLE MALLI UI:sta ladattavaksi
		file, err := os.ReadFile("averaged_model_state_dict.pth")
		if err != nil {
			log.Println("Error reading file:", err)
			return
		}

		// Write file data to WebSocket connection
		err = conn.WriteMessage(websocket.BinaryMessage, file)
		if err != nil {
			log.Println("Error sending file:", err)
			return
		}
	}
}

func filetransferWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Read the message which should be the file content
	_, fileData, err := conn.ReadMessage()
	if err != nil {
		log.Println("Read error:", err)
		return
	}

	// Save the received file data to a file
	datasetZipFileName := "train_set.zip"
	err = os.WriteFile(datasetZipFileName, fileData, 0644)
	if err != nil {
		log.Println("Error saving file:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("File transfer failure"))
		return
	}

	log.Println("File saved successfully")

	helper.UnzipFile(datasetZipFileName, "all_training_data")

	err = conn.WriteMessage(websocket.TextMessage, []byte("File transfer completed successfully"))
	if err != nil {
		log.Println("Error sending completion message:", err)
		return
	}
}
