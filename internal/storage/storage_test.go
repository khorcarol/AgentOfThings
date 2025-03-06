package storage

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/khorcarol/AgentOfThings/internal/api"
)

// mockDirProvider can be used in tests
type mockDirProvider struct {
	homeDir   string
	configDir string
	err       error
}

func (p mockDirProvider) GetHomeDir() (string, error) {
	return p.homeDir, p.err
}

func (p mockDirProvider) GetConfigDir() (string, error) {
	return p.configDir, p.err
}

// TestGetStorageDir tests the GetStorageDir function
func TestGetStorageDir(t *testing.T) {
	// Save the original provider to restore later
	originalProvider := provider
	defer func() { provider = originalProvider }()

	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Set up the mock provider
	provider = mockDirProvider{
		configDir: tempDir,
		err:       nil,
	}

	// Test successful directory creation
	dir, err := GetStorageDir()
	if err != nil {
		t.Fatalf("GetStorageDir() error = %v", err)
	}

	expectedDir := filepath.Join(tempDir, appDirName)
	if dir != expectedDir {
		t.Errorf("GetStorageDir() = %v, want %v", dir, expectedDir)
	}

	// Verify directory was created
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("Directory was not created: %v", expectedDir)
	}

	// Test error handling
	provider = mockDirProvider{
		configDir: "",
		err:       os.ErrPermission,
	}

	_, err = GetStorageDir()
	if err != os.ErrPermission {
		t.Errorf("Expected permission error, got %v", err)
	}
}

func TestSaveLoadFriends(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()
	testStorageDir := filepath.Join(tempDir, appDirName)
	if err := os.MkdirAll(testStorageDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Save original provider to restore later
	originalProvider := provider
	defer func() { provider = originalProvider }()

	// Set up the mock provider
	provider = mockDirProvider{
		configDir: tempDir,
		err:       nil,
	}

	// Sample data
	friends := map[api.ID]api.User{
		api.ID{Address: "user1"}: {UserID: api.ID{Address: "user1"}, CommonInterests: []api.Interest{}, Seen: false},
		api.ID{Address: "user2"}: {UserID: api.ID{Address: "user2"}, CommonInterests: []api.Interest{}, Seen: false},
	}

	// Test saving friends
	if err := SaveFriends(friends); err != nil {
		t.Fatalf("SaveFriends() error = %v", err)
	}

	// Verify file was created
	filePath := filepath.Join(testStorageDir, friendsFileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("File was not created: %v", filePath)
	}

	// Test loading friends
	loadedFriends, err := LoadFriends()
	if err != nil {
		t.Fatalf("LoadFriends() error = %v", err)
	}

	if !reflect.DeepEqual(friends, loadedFriends) {
		t.Errorf("LoadFriends() = %v, want %v", loadedFriends, friends)
	}

	// Test loading empty friends file
	os.Remove(filePath)
	emptyFriends, err := LoadFriends()
	if err != nil {
		t.Fatalf("LoadFriends() error = %v", err)
	}
	if len(emptyFriends) != 0 {
		t.Errorf("Expected empty map, got %v", emptyFriends)
	}

	// Test error cases
	// Invalid JSON
	if err := os.WriteFile(filePath, []byte("{invalid json}"), 0644); err != nil {
		t.Fatalf("Failed to write invalid JSON: %v", err)
	}
	_, err = LoadFriends()
	if err == nil {
		t.Error("LoadFriends() should have failed with invalid JSON")
	}
}

func TestSaveLoadFriends_EdgeCases(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()
	testStorageDir := filepath.Join(tempDir, appDirName)
	if err := os.MkdirAll(testStorageDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Save original provider to restore later
	originalProvider := provider
	defer func() { provider = originalProvider }()

	// Set up the mock provider
	provider = mockDirProvider{
		configDir: tempDir,
		err:       nil,
	}

	// Test with a directory that already exists
	if err := os.MkdirAll(testStorageDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Test with a file that has incorrect permissions
	filePath := filepath.Join(testStorageDir, friendsFileName)
	if err := os.WriteFile(filePath, []byte("{}"), 0000); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	_, err := LoadFriends()
	if err == nil {
		t.Error("LoadFriends() should have failed with incorrect file permissions")
	}

	// Test with a corrupted JSON file
	if err := os.WriteFile(filePath, []byte("{invalid json}"), 0644); err != nil {
		t.Fatalf("Failed to write invalid JSON: %v", err)
	}

	_, err = LoadFriends()
	if err == nil {
		t.Error("LoadFriends() should have failed with invalid JSON")
	}
}
