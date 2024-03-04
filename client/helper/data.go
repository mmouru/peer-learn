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
	fmt.Println("zipfileName:", zipFileName, "folder:", folderName)
	ls := exec.Command("ls")
	outputs, err := ls.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the output
	fmt.Println(string(outputs))

	cmd := exec.Command("unzip", "-o", zipFileName, "-d", nameForFolder)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing unzip command:", err)
		fmt.Println("Command output:", string(output))
		return
	}

	// Print the output of the command
	fmt.Println(string(output))

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
	fmt.Println("Start dataset splitting", fmt.Sprintf("%s %s", fmt.Sprint(n_peers), trainingDataFolder))
	cmd := exec.Command("./helper/splits.py", fmt.Sprint(n_peers), trainingDataFolder)

	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
	}
}
