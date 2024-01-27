package cache

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

const cacheDir = ".gogetty"
const ClientList = "clients.csv"
const moduleDir = "modules"

// Returns the cache directory's absolute path.
func CacheDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory.", err)
	}

	return filepath.Join(homeDir, cacheDir)
}

// Returns the module directory's absolute path
func ModuleDir() string {
	return filepath.Join(CacheDir(), moduleDir)
}

// Creates the directory structure, and writes the cache map file.
func Init() error {
	// Create the modules directory
	if err := os.MkdirAll(ModuleDir(), 0755); err != nil {
		return err
	}

	// Initialize the client list
	if err := writeClients([]string{}); err != nil {
		return err
	}

	return nil
}

func GetClients() ([]string, error) {
	return readClients()
}

func AddClient(clientPath string) error {

	clients, err := readClients()
	if err != nil {
		fmt.Println("Error reading client list.")
		return err
	}

	// Check if the client already exists in the list
	for _, existingClient := range clients {
		if existingClient == clientPath {
			return fmt.Errorf("client already exists: %s", clientPath)
		}
	}

	clients = append(clients, clientPath)
	return writeClients(clients)
}

func Remove(modulePath string) error {
	err := os.RemoveAll(modulePath)
	if err != nil {
		// Handle the error if needed
		return err
	}
	return nil // Return nil to indicate success
}

func readClients() ([]string, error) {
	cacheFilePath := filepath.Join(CacheDir(), ClientList)

	file, err := os.Open(cacheFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var clients []string
	for _, record := range records {
		// Assuming each record contains a single client name
		if len(record) > 0 {
			clients = append(clients, record[0])
		}
	}

	return clients, nil
}

func writeClients(clients []string) error {
	cacheFilePath := filepath.Join(CacheDir(), ClientList)

	file, err := os.Create(cacheFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, client := range clients {
		if err := writer.Write([]string{client}); err != nil {
			return err
		}
	}

	return nil
}
