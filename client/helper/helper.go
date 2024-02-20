package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

type IP struct {
	Query string
}

type PeerPayload struct {
	IpAddr string `json:"ip"`
	Port   string `json:"port"`
	P2PId  string `json:"peerId"`
	Cuda   bool   `json:"hasCuda"`
}

func GetPublicIp() string {
	res, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var ip IP
	json.Unmarshal(body, &ip)

	return ip.Query
}

func RegisterPeerToCentralList(port string, p2pId string) int {
	pubIpAddr := GetPublicIp()
	hasCUDA, err := HasCUDAGPU()

	if err != nil {
		fmt.Println("error in finding cuda")
		hasCUDA = false
	}

	payload := PeerPayload{
		IpAddr: pubIpAddr,
		Port:   port,
		P2PId:  p2pId,
		Cuda:   hasCUDA,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}

	apiURL := "https://peer-service-qfobv32vvq-lz.a.run.app/api/v1/register"

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	defer resp.Body.Close()

	// Check the response status
	return resp.StatusCode
}

func AppearOfflineOnList() {

}

func HasCUDAGPU() (bool, error) {
	cmd := exec.Command("nvidia-smi")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}

	// Check if the output contains information about NVIDIA GPUs
	return strings.Contains(string(output), "CUDA"), nil
}
