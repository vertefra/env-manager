package manager

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type EnvFile struct {
	header      *Header
	fileContent string
	encrypted   string
	folderPath  string // Where the encrypted file is saved
}

func (e *EnvFile) RestoreAs() string {
	return e.fileContent
}

func (e *EnvFile) Identifier() string {
	return e.header.Identifier
}

func (e *EnvFile) IsEncrypted() bool {
	return e.encrypted != ""
}

func (e *EnvFile) Headers() []string {
	return e.header.String()
}

func (e *EnvFile) readRestoreAs() {
	// search for #- restore-as: <filename>
	// in the header
	// If not found, use the default name
	if e.header.RestoreAs == "" {
		e.header.RestoreAs = DEFAULT_RESTORE_AS
	}
}

func (e *EnvFile) readHeader() error {
	h, err := InitHeader(e.fileContent)
	if err != nil {
		return err
	}
	e.header = h
	return nil
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

/// Functions

func GetEnvFile(identifier string, folder *string) (*EnvFile, error) {
	f, err := GetOrCreateFolder(folder)

	if err != nil {
		return nil, err
	}

	allIdentifiers := f.GetIdentifiers()

	for _, id := range allIdentifiers {
		if string(id) == identifier {
			return ReadEnvFile(fmt.Sprintf("%s/%s%s", f.FolderPath, SAVED_PREFIX, identifier)), nil
		}
	}

	return nil, errors.New("invalid identifier - " + identifier)
}

func GetEnvFiles(folder *string) ([]*EnvFile, error) {
	f, err := GetOrCreateFolder(folder)
	if err != nil {
		return nil, err
	}

	allIdentifiers := f.GetIdentifiers()

	var envFiles []*EnvFile

	for _, id := range allIdentifiers {
		filePath := fmt.Sprintf("%s/%s%s", *folder, SAVED_PREFIX, id)
		e := ReadEnvFile(filePath)
		envFiles = append(envFiles, e)
	}

	return envFiles, nil
}

func RestoreEnvFile(e *EnvFile, decryptSecret string) {
	e.decrypt(decryptSecret)

	// Re-parse header from decrypted content to get correct restoreAs
	h, err := InitHeader(e.fileContent)
	if err == nil {
		e.header = h
	}

	e.readRestoreAs()

	fmt.Printf("Restoring file %s as %s\n", e.folderPath, e.header.RestoreAs)

	f, err := os.Create(e.header.RestoreAs)

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

	filePath := fmt.Sprintf("%s/%s%s", e.folderPath, SAVED_PREFIX, e.header.Identifier)
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

func InitEnvFile(identifier string, restoreAs string) *EnvFile {
	h := &Header{
		Identifier: identifier,
		RestoreAs:  restoreAs,
	}
	return &EnvFile{
		header: h,
	}
}

func (e *EnvFile) SetContent(content string) {
	// Add headers to the content so they're preserved when encrypting
	headerContent := fmt.Sprintf("%s%s\n%s%s\n", IDENTIFIER_HEADER, e.header.Identifier, RESTORE_AS_HEADER, e.header.RestoreAs)
	e.fileContent = headerContent + content
}

func ReadEnvFile(filePath string) *EnvFile {
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
	if strings.Contains(filePath, DEFAULT_ENV_FOLDER) || strings.Contains(filePath, SAVED_PREFIX) {
		e.encrypted = c
		e.folderPath = filePath
		// For encrypted files, extract identifier from filename
		parts := strings.Split(filePath, "/")
		filename := parts[len(parts)-1]
		identifier := strings.TrimPrefix(filename, SAVED_PREFIX)
		e.header = &Header{
			Identifier: identifier,
			RestoreAs:  DEFAULT_RESTORE_AS,
		}
	} else {
		e.fileContent = c
		h, err := InitHeader(c)
		if err != nil {
			fmt.Println("Error initializing header")
			fmt.Println(err)
			os.Exit(1)
		}
		e.header = h
	}
	return &e
}
