package manager

import (
	"fmt"
	"os"
)

var FOLDER_NAME = ".env-manager"

func CreateInitFolderIfNotExist(folderName *string) {
	if folderName == nil {
		folderName = &FOLDER_NAME
	}
	fmt.Println("Creating init folder... ")
	// Check if folder exists
	// If not, create it
	if _, err := os.Stat(*folderName); os.IsNotExist(err) {
		fmt.Printf("\nFolder does not exist: %s\n", *folderName)
		fmt.Println("\nCreating folder... ")
		err := os.Mkdir(*folderName, 0755)
		if err != nil {
			fmt.Println("Error creating folder")
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Folder created: %s\n", *folderName)
	} else {
		fmt.Printf("Folder exists: %s\n", *folderName)
	}
}
