package files

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFiles_WriteToJsonFile_Success(t *testing.T) {
	dir := t.TempDir()

	data := map[string]interface{}{
		"foo": "bar",
	}

	path := filepath.Join(dir, "subdir", "file.json")

	err := WriteToJsonFile(path, data)

	if err != nil {
		t.Error(err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Error(err)
	}

	var read map[string]interface{}
	err = json.Unmarshal(raw, &read)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, data, read)
}

func TestFiles_WriteToJsonFile_InvalidMarshall(t *testing.T) {
	dir := t.TempDir()

	path := filepath.Join(dir, "subdir", "file.json")

	err := WriteToJsonFile(path, make(chan int))

	if err == nil {
		t.Fatal("Expected marshal error, got <nil>")
	}
	assert.Contains(t, err.Error(), "json: unsupported type: chan int")
}

func TestFiles_WriteToJsonFile_InvalidPath(t *testing.T) {
	// TODO this might not be valid for non-posix systems
	path := "/dev/null/foo"

	err := WriteToJsonFile(path, nil)

	if err == nil {
		t.Fatal("Expected directory error, got <nil>")
	}
	assert.Equal(t, "failed to create directory for /dev/null/foo: mkdir /dev/null: not a directory", err.Error())
}

func TestFiles_WriteToJsonFile_InvalidPath2(t *testing.T) {
	dir := t.TempDir()

	err := WriteToJsonFile(dir, nil)

	if err == nil {
		t.Fatal("Expected file error, got <nil>")
	}
	assert.Equal(t, fmt.Sprintf("failed to write to file: open %s: is a directory", dir), err.Error())
}
