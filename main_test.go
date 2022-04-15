package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

const (
	inputFilePath  = "./testdata/test1.md"
	goldenFilePath = "./testdata/test1.md.html"
)

func TestParseContent(t *testing.T) {
	input, err := os.ReadFile(inputFilePath)
	if err != nil {
		t.Fatal(err)
	}

	result := parseContent(input)

	expected, err := os.ReadFile(goldenFilePath)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(result, expected) {
		t.Logf("golden:/n%s/n", expected)
		t.Logf("result:/n%s/n", result)

		t.Errorf("Result content does not match golden file")
	}
}

func TestRun(t *testing.T) {
	var mockStdOut bytes.Buffer

	// Call the run method and write the file path on the provided buffer.
	if err := run(inputFilePath, &mockStdOut, true); err != nil {
		t.Fatal(err)
	}

	// Convert the bytes buffer into a space trimmed file path string.
	resultFile := strings.TrimSpace(mockStdOut.String())

	result, err := os.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile(goldenFilePath)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.EqualFold(result, expected) {
		t.Logf("golden:/n%s/n", expected)
		t.Logf("result:/n%s/n", result)

		t.Errorf("Result content does not match golden file")
	}

	os.Remove(resultFile)
}
