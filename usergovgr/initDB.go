package usergovgr

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/mutecomm/go-sqlcipher"
	"golang.org/x/crypto/pbkdf2"
)

func deriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New)
}

func getDatabasePath() (string, error) {
	var baseDir string

	switch runtime.GOOS {
	case "windows":
		baseDir = filepath.Join(os.Getenv("APPDATA"), "YourAppName")
	case "darwin":
		baseDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "YourAppName")
	case "linux":
		baseDir = filepath.Join(os.Getenv("HOME"), ".local", "share", "YourAppName")
	default:
		return "", fmt.Errorf("unsupported platform")
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return "", err
	}

	// Return the full path to the SQLite database file
	return filepath.Join(baseDir, "encrypted.db"), nil
}

func openEncryptedDB(dbPath, password string) (*sql.DB, error) {
	// Open the SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Set the encryption key
	salt := []byte("randomSaltValue") // Should be securely generated and stored
	key := deriveKey(password, salt)
	keyHex := fmt.Sprintf("%x", key)

	_, err = db.Exec(fmt.Sprintf("PRAGMA key = '%s';", keyHex))
	if err != nil {
		return nil, err
	}

	// Verify the key
	_, err = db.Exec("PRAGMA cipher_memory_security = ON;")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	dbPath, err := getDatabasePath()
	if err != nil {
		log.Fatal("Failed to determine database path:", err)
	}

	password := "your-strong-password"
	db, err := openEncryptedDB(dbPath, password)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Perform your database operations here
	// ...
}
