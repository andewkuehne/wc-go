package main

import (
	"bytes"
	"os"
	"testing"
)

func TestCountBytes(t *testing.T) {
	input := bytes.NewBufferString("This is a test.")
	inCounts := countInput(input.Bytes(), true, false, false)
	expected := counts{byteCount: 15}
	if inCounts != expected {
		t.Errorf("Expected %v, but got %v", expected, inCounts)
	}
}

func TestCountWords(t *testing.T) {
	input := bytes.NewBufferString("This is a test.")
	inCounts := countInput(input.Bytes(), false, true, false)
	expected := counts{wordCount: 4}
	if inCounts != expected {
		t.Errorf("Expected %v, but got %v", expected, inCounts)
	}
}

func TestCountLines(t *testing.T) {
	input := bytes.NewBufferString("This\nis\na\ntest.")
	inCounts := countInput(input.Bytes(), false, false, true)
	expected := counts{lineCount: 4}
	if inCounts != expected {
		t.Errorf("Expected %v, but got %v", expected, inCounts)
	}
}

func TestCountAll(t *testing.T) {
	input := bytes.NewBufferString("This is a test.\nThis is only a test.")
	inCounts := countInput(input.Bytes(), true, true, true)
	expected := counts{byteCount: 34, wordCount: 9, lineCount: 2}
	if inCounts != expected {
		t.Errorf("Expected %v, but got %v", expected, inCounts)
	}
}

func TestCountFile(t *testing.T) {
	// Create temporary file with test data
	data := []byte("This is a test.\nThis is only a test.")
	tmpfile, err := os.CreateTemp("", "wc-go-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Count using command-line arguments
	args := []string{"-lwc", tmpfile.Name()}
	inCounts, err := countArgs(args)
	if err != nil {
		t.Fatal(err)
	}

	// Check counts
	expected := counts{byteCount: 34, wordCount: 9, lineCount: 2}
	if inCounts != expected {
		t.Errorf("Expected %v, but got %v", expected, inCounts)
	}
}

func TestCountNoArgs(t *testing.T) {
	input := bytes.NewBufferString("This is a test.\nThis is only a test.")
	inCounts := countInput(input.Bytes(), false, false, false)
	expected := counts{byteCount: 34, wordCount: 9, lineCount: 2}
	if inCounts != expected {
		t.Errorf("Expected %v, but got %v", expected, inCounts)
	}
}

func TestCountInvalidFile(t *testing.T) {
	args := []string{"-lwc", "invalid-file.txt"}
	_, err := countArgs(args)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
}
