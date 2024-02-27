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

	// Command to run the unzip command
	cmd := exec.Command("unzip", zipFileName, "-d", folderName)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing unzip command:", err)
		return
	}

	// Print the output of the command
	fmt.Println(string(output))

	fmt.Println("Extraction completed successfully.")
}

func SplitTrainingDataAmongPeers(n_peers int, trainingDataFolder string) {

	defer os.RemoveAll(trainingDataFolder)
	fmt.Println("NO MOROOOOO")
	cmd := exec.Command("./helper/splits.py", fmt.Sprint(n_peers), trainingDataFolder)

	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
	}
}
