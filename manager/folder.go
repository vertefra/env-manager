package manager

import (
	"encoding/json"
	"fmt"
	"os"
)

type EnvFilePath string
type EnvFileIdentifier string

type Manifest struct {
	Identifiers map[EnvFilePath]EnvFileIdentifier `json:"identifiers"`
}

func (m *Manifest) Write(folderPath string) error {
	manifestPath := fmt.Sprintf("%s/%s", folderPath, "manifest.json")
	f, err := os.Create(manifestPath)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(m)
}

func (m *Manifest) Load(folderPath string) error {
	manifestPath := fmt.Sprintf("%s/%s", folderPath, "manifest.json")
	f, err := os.Open(manifestPath)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manifest) addFileIdentifier(filePath EnvFilePath, identifier EnvFileIdentifier, folderPath string) error {
	if m.Identifiers == nil {
		m.Identifiers = make(map[EnvFilePath]EnvFileIdentifier)
	}
	m.Identifiers[filePath] = identifier
	return m.Write(folderPath)
}

func (m *Manifest) EvictFileIdentifier(filePath EnvFilePath, folderPath string) error {
	if m.Identifiers != nil {
		delete(m.Identifiers, filePath)
		return m.Write(folderPath)
	}
	return nil
}

type Folder struct {
	manifest   *Manifest
	FolderPath string
}

func (f *Folder) AddFileIdentifier(filePath EnvFilePath, identifier EnvFileIdentifier) error {
	return f.manifest.addFileIdentifier(filePath, identifier, f.FolderPath)
}

func (f *Folder) EvictFileIdentifier(filePath EnvFilePath) error {
	return f.manifest.EvictFileIdentifier(filePath, f.FolderPath)
}

func (f *Folder) GetIdentifiers() map[EnvFilePath]EnvFileIdentifier {
	return f.manifest.Identifiers
}

func GetOrCreateFolder(folderName *string) (*Folder, error) {
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
			return nil, err
		}
	} else {
		fmt.Printf("Folder exists: %s\n", *folderName)
	}

	folder := &Folder{
		manifest: &Manifest{
			Identifiers: make(map[EnvFilePath]EnvFileIdentifier),
		},
		FolderPath: *folderName,
	}

	// Try to load existing manifest
	folder.manifest.Load(folder.FolderPath)

	return folder, nil
}
