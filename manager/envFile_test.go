package manager

import (
	"fmt"
	"os"
	"testing"
)

func destroyTestFolder(folderPath *string) {
	os.RemoveAll(*folderPath)
}

func createEnvFile(path string, content string) {
	f, err := os.Create(path)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err = f.WriteString(content)

	if err != nil {
		panic(err)
	}

}

func deleteEnvFile(path string) {
	os.Remove(path)
}

func getEnvFileContent(identifier string, keyValuePairs ...string) string {
	content := fmt.Sprintf("#- identifier: %s\n", identifier)

	for _, pair := range keyValuePairs {
		content += fmt.Sprintf("%s\n", pair)
	}

	return content
}

func TestReadEnvFile(t *testing.T) {
	const ENV_FILE_PATH = ".env-test"
	const ENV_FILE_CONTENT = "HELLO=WORLD\n"
	const ENV_FILE_IDENTIFIER = "test"

	defer func() {
		deleteEnvFile(ENV_FILE_PATH)
	}()

	content := getEnvFileContent(ENV_FILE_IDENTIFIER, ENV_FILE_CONTENT)
	createEnvFile(ENV_FILE_PATH, content)

	// folderPath := createTestFolder()

	e := ReadEnvFile(ENV_FILE_PATH, false)

	wantIdentifier := ENV_FILE_IDENTIFIER == e.identifier
	wantContent := content == e.fileContent
	wantEncrypted := e.encrypted == ""

	if !wantIdentifier {
		t.Errorf("ReadEnvFile() = %v, want %v", e.identifier, ENV_FILE_IDENTIFIER)
	}

	if !wantContent {
		t.Errorf("ReadEnvFile() = %v, want %v", e.fileContent, ENV_FILE_CONTENT)
	}

	if !wantEncrypted {
		t.Errorf("ReadEnvFile() = %v, want %v", e.encrypted, "")
	}
}

func TestSaveAndReadEnvFile(t *testing.T) {
	const ENV_FILE_PATH = ".env-test"
	const ENV_FILE_CONTENT = "HELLO=WORLD\n"
	const ENV_FILE_IDENTIFIER = "test"
	const ENCRYPT_SECRET = "488c447d4919b142c80c82832cef7f18"
	var FOLDER_PATH = ".env-manager-test"

	defer func() {
		fmt.Println("Destroying test folder")
		destroyTestFolder(&FOLDER_PATH)
		deleteEnvFile(ENV_FILE_PATH)
		deleteEnvFile(".env")
	}()

	content := getEnvFileContent(ENV_FILE_IDENTIFIER, ENV_FILE_CONTENT)
	createEnvFile(ENV_FILE_PATH, content)

	// Simulate an init operation
	// Read the file given by the user (created as a fixture)
	e := ReadEnvFile(ENV_FILE_PATH, false)
	CreateInitFolderIfNotExist(&FOLDER_PATH)
	// Inject custom folder path
	SaveEnvFile(e, ENCRYPT_SECRET, &FOLDER_PATH)

	// Test List Env Files
	envFiles := GetEnvFiles(&FOLDER_PATH)

	wantEnvFiles := len(envFiles) == 1

	if !wantEnvFiles {
		t.Errorf("GetEnvFiles() = %v, want %v", len(envFiles), 1)
	}

	// Read env file in the folder
	e = ReadEnvFile(fmt.Sprintf("%s/%s%s", FOLDER_PATH, SAVED_PREFIX, ENV_FILE_IDENTIFIER), true)
	println(e.encrypted)
	// File is read and encrypted, the identifier is unknown
	wantIdentifier := e.identifier == ENV_FILE_IDENTIFIER
	wantEncrypted := e.encrypted != ""

	if !wantIdentifier {
		t.Errorf("SaveEnvFile() = %v, want %v", e.identifier, ENV_FILE_IDENTIFIER)
	}

	if !wantEncrypted {
		t.Errorf("SaveEnvFile() = %v, want %v", e.encrypted, "")
	}

	// Simulate a get operation
	toRestore := GetEnvFile(ENV_FILE_IDENTIFIER, &FOLDER_PATH)

	wantRestoredIdentifier := toRestore.identifier == ENV_FILE_IDENTIFIER

	if !wantRestoredIdentifier {
		t.Errorf("GetEnvFile() = %v, want %v", toRestore.identifier, ENV_FILE_IDENTIFIER)
	}

	wantRestoredEncrypted := toRestore.encrypted != ""

	if !wantRestoredEncrypted {
		t.Errorf("GetEnvFile() = %v, want %v", toRestore.encrypted, "")
	}

	wantDecryptedContent := toRestore.fileContent == ""

	if !wantDecryptedContent {
		t.Errorf("GetEnvFile() = %v, want %v", toRestore.fileContent, "")
	}

	// Restore env file
	RestoreEnvFile(toRestore, ENCRYPT_SECRET)

	RESTORED := ".env"

	// Check if the file was restored
	restoredContent, err := os.ReadFile(RESTORED)

	if err != nil {
		t.Errorf("RestoreEnvFile() = %v, want %v", err, nil)
	}

	// Make sure the content is the same
	wantRestoredContent := string(restoredContent) == content

	if !wantRestoredContent {
		t.Errorf("RestoreEnvFile() = %v, want %v", string(restoredContent), content)
	}
}
