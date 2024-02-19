package manager

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

type EnvFile struct {
	header      []string
	fileContent string
	encrypted   string
	identifier  string
	folderPath  string // Where the encrypted file is saved
	restoreAs   string // Name of the decrypted file
}

func (e *EnvFile) readRestoreAs() {
	// search for #- restore-as: <filename>
	// in the header
	// If not found, use the default name
	if len(e.header) == 0 {
		e.readHeader()
	}
	for _, line := range e.header {
		if strings.HasPrefix(line, RESTORE_AS_HEADER) {
			r := strings.TrimPrefix(line, RESTORE_AS_HEADER)
			e.restoreAs = strings.Trim(r, " ")
			return
		}
	}

	if e.restoreAs == "" {
		e.restoreAs = DEFAULT_RESTORE_AS
	}
}

func (e *EnvFile) readIdentifier() {
	// Search for #- identifier: <identifier>
	// in the header
	// If not found, exit with error

	for _, line := range e.header {
		if strings.HasPrefix(line, IDENTIFIER_HEADER) {
			i := strings.TrimPrefix(line, IDENTIFIER_HEADER)
			e.identifier = strings.Trim(i, " ")
			return
		}
	}

	if e.identifier == "" {
		fmt.Println("Identifier not found")
		os.Exit(1)
	}
}

func (e *EnvFile) readHeader() {
	if e.fileContent == "" {
		fmt.Println("File content is empty")
		os.Exit(1)
	}

	// Read first line
	// If it starts with #-, it's a header

	lines := strings.Split(e.fileContent, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "#-") {
			e.header = append(e.header, line)
		}
	}
}

func (e *EnvFile) encrypt(key string) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	plaintext := []byte(e.fileContent)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err = rand.Read(iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	e.encrypted = hex.EncodeToString(ciphertext)
}

func (e *EnvFile) decrypt(key string) {
	ciphertext, err := hex.DecodeString(e.encrypted)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	e.fileContent = string(ciphertext)
}

func getEnvFileList(folderPath *string) []string {
	files, err := os.ReadDir(*folderPath)
	if err != nil {
		panic(err)
	}

	var fileNames []string

	for _, file := range files {
		var name = file.Name()
		// File gets saved as .env.<identifier>
		// split the name and keep the identifier
		name = strings.TrimPrefix(name, SAVED_PREFIX)
		fileNames = append(fileNames, name)
	}

	return fileNames
}

func GetEnvFile(identifier string, folder *string) *EnvFile {

	allIdentifiers := getEnvFileList(folder)

	validIdentifier := false
	for _, id := range allIdentifiers {
		if id == identifier {
			validIdentifier = true
			break
		}
	}

	if !validIdentifier {
		panic("Invalid identifier - " + identifier)
	}

	filePath := fmt.Sprintf("%s/%s%s", *folder, SAVED_PREFIX, identifier)
	e := ReadEnvFile(filePath, true)
	return e
}

///
/// Public functions
///

func RestoreEnvFile(e *EnvFile, decryptSecret string) {
	e.decrypt(decryptSecret)

	e.readRestoreAs()

	fmt.Printf("Restoring file %s as %s\n", e.folderPath, e.restoreAs)

	f, err := os.Create(e.restoreAs)

	if err != nil {
		fmt.Println("Error creating file")
		fmt.Println(err)
		os.Exit(1)
	}

	defer f.Close()

	// Write file

	_, err = f.WriteString(e.fileContent)

	if err != nil {
		fmt.Println("Error writing file")
		fmt.Println(err)
		os.Exit(1)
	}
}

// SaveEnvFile saves the environment file to the env-manager folder
// in the encrypted format
func SaveEnvFile(e *EnvFile, encryptSecret string, folderPath *string) {
	e.encrypt(encryptSecret)

	if e.folderPath == "" {
		e.folderPath = *folderPath
	}

	filePath := fmt.Sprintf("%s/%s%s", e.folderPath, SAVED_PREFIX, e.identifier)
	fmt.Printf("Saving file: %s\n", filePath)

	f, err := os.Create(filePath)

	if err != nil {
		fmt.Println("Error creating file")
		fmt.Println(err)
		os.Exit(1)
	}

	defer f.Close()

	// Write file

	_, err = f.WriteString(e.encrypted)

	if err != nil {
		fmt.Println("Error writing file")
		fmt.Println(err)
		os.Exit(1)
	}
}

func ReadEnvFile(filePath string, fromEncrypted bool) *EnvFile {
	fmt.Printf("Reading file: %s\n", filePath)

	f, err := os.Open(filePath)

	if err != nil {
		fmt.Println("Error opening file")
		fmt.Println(err)
		os.Exit(1)
	}

	defer f.Close()

	// Read file

	fileBytes, err := io.ReadAll(f)

	if err != nil {
		fmt.Println("Error reading file")
		fmt.Println(err)
		os.Exit(1)
	}

	c := string(fileBytes)

	e := EnvFile{}
	// If filepath starts with the .env-manager folder
	// then the content is encrypted
	if fromEncrypted {
		// Identifier is <filename>.<identifier>
		sgmts := strings.Split(filePath, ".")
		identifier := sgmts[len(sgmts)-1]
		e.encrypted = c
		e.identifier = identifier
	} else {
		e.fileContent = c
		e.readHeader()
		e.readIdentifier()
	}
	return &e
}
