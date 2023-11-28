package files

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFiles_SafeWriteObjectToJSONFile_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "file.json")

	data := map[string]interface{}{
		"foo": "bar",
	}

	err := SafeWriteObjectToJSONFile(path, data)
	if err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var read map[string]interface{}
	err = json.Unmarshal(raw, &read)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, data, read)
}

func TestFiles_SafeWriteObjectToJSONFile_InvalidMarshall(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "file.json")
	err := SafeWriteObjectToJSONFile(path, make(chan int))

	assert.Error(t, err)
	assert.ErrorContains(t, err, "json: unsupported type: chan int")
}

func TestFiles_SafeWriteObjectToJSONFile_InvalidPath(t *testing.T) {
	dir := t.TempDir()
	// create a file in the temp dir
	file, err := os.CreateTemp(dir, "test")
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	path := filepath.Join(file.Name(), "file.json")
	err = SafeWriteObjectToJSONFile(path, nil)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "not a directory")
}

func TestFiles_SafeWriteObjectToJsonFile_InvalidPath2(t *testing.T) {
	dir := t.TempDir()
	err := SafeWriteObjectToJSONFile(dir, nil)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "is a directory")
}
