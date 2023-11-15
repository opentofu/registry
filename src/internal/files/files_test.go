package files

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
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

	if !reflect.DeepEqual(data, read) {
		t.Errorf("Expected %#v, got %#v", data, read)
	}
}

func TestFiles_WriteToJsonFile_InvalidMarshall(t *testing.T) {
	dir := t.TempDir()

	path := filepath.Join(dir, "subdir", "file.json")

	err := WriteToJsonFile(path, make(chan int))

	if err == nil {
		t.Error("Expected marshal error, got <nil>")
	} else if !strings.Contains(err.Error(), "json: unsupported type: chan int") {
		t.Errorf("Expected marshal error, got %s", err.Error())
	}
}

func TestFiles_WriteToJsonFile_InvalidPath(t *testing.T) {
	// TODO this might not be valid for non-posix systems
	path := "/dev/null/foo"

	err := WriteToJsonFile(path, nil)

	if err == nil {
		t.Error("Expected directory error, got <nil>")
	} else if err.Error() != "failed to create directory for /dev/null/foo: mkdir /dev/null: not a directory" {
		t.Errorf("Expected directory error, got %s", err.Error())
	}
}

func TestFiles_WriteToJsonFile_InvalidPath2(t *testing.T) {
	dir := t.TempDir()

	err := WriteToJsonFile(dir, nil)

	if err == nil {
		t.Error("Expected file error, got <nil>")
	} else if err.Error() != fmt.Sprintf("failed to write to file: open %s: is a directory", dir) {
		t.Errorf("Expected file error, got %s", err.Error())
	}
}
