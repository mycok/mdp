package main

import (
	"bytes"
	"os"
	"testing"
)

const (
	inputFilePath = "./testdata/test1.md"
	goldenFilePath = "./testdata/test1.md.html"
	resultFile = "test1.md.html"
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
	if err := run(inputFilePath); err != nil {
		t.Fatal(err)
	}


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