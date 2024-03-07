package helper

import (
	"fmt"
	"os"
	"os/exec"
)

func UnzipFile(zipFileName string, nameForFolder string) {
	folderName := nameForFolder
	if folderName == "" {
		folderName = "data"
	}

	err := os.MkdirAll(folderName, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	cmd := exec.Command("unzip", "-o", zipFileName, "-d", nameForFolder)

	done := make(chan bool)
	go Spinner("Extracting training data", done)

	// Run the command
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error executing unzip command:", err)
		return
	}

	done <- true

	fmt.Println("Extraction completed successfully.")

}

func ZipFile(filepath string, zipfileName string) {

	cmd := exec.Command(fmt.Sprintf("zip %s %s", zipfileName, filepath))

	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
	}
}

func SplitTrainingDataAmongPeers(n_peers int, trainingDataFolder string) {

	defer os.RemoveAll(trainingDataFolder)

	done := make(chan bool)
	go Spinner("Splitting the training set", done)

	cmd := exec.Command("./helper/splits.py", fmt.Sprint(n_peers), trainingDataFolder)

	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
	}

	done <- true

	fmt.Println(fmt.Printf("Training set split into %d sets.\n", n_peers))
}
