package helper

import (
	"fmt"
	"os"
	"os/exec"
)

func UnzipTrainingSet(zipFileName string, nameForFolder string) {
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
