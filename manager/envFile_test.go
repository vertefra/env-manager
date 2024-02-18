package manager

import (
	"fmt"
	"os"
	"testing"
)

func createTestFolder() string {
	var FOLDER_NAME = ".env-manager-test"
	CreateInitFolderIfNotExist(&FOLDER_NAME)
	return FOLDER_NAME
}

func destroyTestFolder() {
	var FOLDER_NAME = ".env-manager-test"
	os.RemoveAll(FOLDER_NAME)
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
		destroyTestFolder()
		deleteEnvFile(ENV_FILE_PATH)
	}()

	content := getEnvFileContent(ENV_FILE_IDENTIFIER, ENV_FILE_CONTENT)
	createEnvFile(ENV_FILE_PATH, content)

	// folderPath := createTestFolder()

	e := ReadEnvFile(ENV_FILE_PATH)

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

func TestSaveEnvFile(t *testing.T) {
	const ENV_FILE_PATH = ".env-test"
	const ENV_FILE_CONTENT = "HELLO=WORLD\n"
	const ENV_FILE_IDENTIFIER = "test"
	const ENCRYPT_SECRET = "asvcdswfjmrewlfgjrflvqmcewcme"

	defer func() {
		destroyTestFolder()
		deleteEnvFile(ENV_FILE_PATH)
	}()

	content := getEnvFileContent(ENV_FILE_IDENTIFIER, ENV_FILE_CONTENT)
	createEnvFile(ENV_FILE_PATH, content)

	e := ReadEnvFile(ENV_FILE_PATH)

	folderPath := createTestFolder()

	SaveEnvFile(e, ENCRYPT_SECRET)

	// Rea env file in the folder
	e = ReadEnvFile(fmt.Sprintf("%s/%s%s", folderPath, SAVED_PREFIX, ENV_FILE_IDENTIFIER))

	wantIdentifier := ENV_FILE_IDENTIFIER == e.identifier

	if !wantIdentifier {
		t.Errorf("SaveEnvFile() = %v, want %v", e.identifier, ENV_FILE_IDENTIFIER)
	}
}
