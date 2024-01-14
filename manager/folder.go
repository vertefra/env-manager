package manager

import (
	"fmt"
	"os"
)

const FOLDER_NAME = ".env-manager"

func CreateInitFolderIfNotExist() {
	fmt.Println("Creating init folder... ")
	// Check if folder exists
	// If not, create it
	if _, err := os.Stat(FOLDER_NAME); os.IsNotExist(err) {
		fmt.Printf("\nFolder does not exist: %s\n", FOLDER_NAME)
		fmt.Println("\nCreating folder... ")
		err := os.Mkdir(FOLDER_NAME, 0755)
		if err != nil {
			fmt.Println("Error creating folder")
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Folder exists: %s\n", FOLDER_NAME)
	}
}
