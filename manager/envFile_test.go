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
	content := fmt.Sprintf("#- identifier: %s\n#- restore-as: %s\n", identifier, DEFAULT_RESTORE_AS)

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

	e := ReadEnvFile(ENV_FILE_PATH)

	wantIdentifier := ENV_FILE_IDENTIFIER == e.Identifier()
	wantContent := content == e.fileContent
	wantEncrypted := e.encrypted == ""

	if !wantIdentifier {
		t.Errorf("ReadEnvFile() = %v, want %v", e.Identifier(), ENV_FILE_IDENTIFIER)
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
	e := ReadEnvFile(ENV_FILE_PATH)
	f, err := GetOrCreateFolder(&FOLDER_PATH)
	if err != nil {
		t.Errorf("GetOrCreateFolder() = %v, want %v", err, nil)
	}
	f.AddFileIdentifier(EnvFilePath(ENV_FILE_PATH), EnvFileIdentifier(ENV_FILE_IDENTIFIER))
	// Inject custom folder path
	SaveEnvFile(e, ENCRYPT_SECRET, &f.FolderPath)

	// Read env file in the folder
	e, err = GetEnvFile(ENV_FILE_IDENTIFIER, &f.FolderPath)
	if err != nil {
		t.Errorf("GetEnvFile() = %v, want %v", err, nil)
	}
	println(e.encrypted)
	// File is read and encrypted, the identifier is unknown
	wantIdentifier := e.Identifier() == ENV_FILE_IDENTIFIER
	wantEncrypted := e.encrypted != ""

	if !wantIdentifier {
		t.Errorf("SaveEnvFile() = %v, want %v", e.Identifier(), ENV_FILE_IDENTIFIER)
	}

	if !wantEncrypted {
		t.Errorf("SaveEnvFile() = %v, want %v", e.encrypted, "")
	}

	// Simulate a get operation
	toRestore, err := GetEnvFile(ENV_FILE_IDENTIFIER, &f.FolderPath)
	if err != nil {
		t.Errorf("GetEnvFile() = %v, want %v", err, nil)
	}

	wantRestoredIdentifier := toRestore.Identifier() == ENV_FILE_IDENTIFIER

	if !wantRestoredIdentifier {
		t.Errorf("GetEnvFile() = %v, want %v", toRestore.Identifier(), ENV_FILE_IDENTIFIER)
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

func TestInitEnvFile(t *testing.T) {
	const ENV_FILE_IDENTIFIER = "production"
	const ENV_FILE_RESTORE_AS = ".env.production"
	const ENCRYPT_SECRET = "488c447d4919b142c80c82832cef7f18"
	var FOLDER_PATH = ".env-manager-test-init"

	defer func() {
		fmt.Println("Destroying test folder")
		destroyTestFolder(&FOLDER_PATH)
		deleteEnvFile(ENV_FILE_RESTORE_AS)
	}()

	// Create env file using InitEnvFile
	e := InitEnvFile(ENV_FILE_IDENTIFIER, ENV_FILE_RESTORE_AS)

	// Verify identifier and restoreAs are set correctly
	if e.Identifier() != ENV_FILE_IDENTIFIER {
		t.Errorf("InitEnvFile() identifier = %v, want %v", e.Identifier(), ENV_FILE_IDENTIFIER)
	}

	if e.header.RestoreAs != ENV_FILE_RESTORE_AS {
		t.Errorf("InitEnvFile() restoreAs = %v, want %v", e.header.RestoreAs, ENV_FILE_RESTORE_AS)
	}

	// Set content (this will add headers automatically)
	rawContent := "DB_HOST=localhost\nDB_PORT=5432\n"
	e.SetContent(rawContent)

	// Create folder and save
	f, err := GetOrCreateFolder(&FOLDER_PATH)
	if err != nil {
		t.Errorf("GetOrCreateFolder() = %v, want %v", err, nil)
	}

	f.AddFileIdentifier(EnvFilePath("manual"), EnvFileIdentifier(ENV_FILE_IDENTIFIER))
	SaveEnvFile(e, ENCRYPT_SECRET, &f.FolderPath)

	// Read it back
	e2, err := GetEnvFile(ENV_FILE_IDENTIFIER, &f.FolderPath)
	if err != nil {
		t.Errorf("GetEnvFile() = %v, want %v", err, nil)
	}

	if e2.Identifier() != ENV_FILE_IDENTIFIER {
		t.Errorf("Retrieved identifier = %v, want %v", e2.Identifier(), ENV_FILE_IDENTIFIER)
	}

	// Restore and verify content
	RestoreEnvFile(e2, ENCRYPT_SECRET)

	restoredContent, err := os.ReadFile(ENV_FILE_RESTORE_AS)
	if err != nil {
		t.Errorf("RestoreEnvFile() = %v, want %v", err, nil)
	}

	// The restored file should contain the headers + raw content
	expectedContent := fmt.Sprintf("#- identifier: %s\n#- restore-as: %s\n%s", ENV_FILE_IDENTIFIER, ENV_FILE_RESTORE_AS, rawContent)
	if string(restoredContent) != expectedContent {
		t.Errorf("Restored content = %v, want %v", string(restoredContent), expectedContent)
	}
}
