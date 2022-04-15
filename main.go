package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = `<!doctype html>
	<html>
		<head>
			<meta http-equiv="content-type" content="html/text; charset=utf-8">
			<title>Markdown Preview Tool</title>
		</head>
		<body>
`

	footer = `
		</body>
	</html>`
)

func main() {
	filename := flag.String("file", "", "Markdown file to preview")
	flag.Parse()

	if *filename == "" {
		flag.Usage()

		os.Exit(1)
	}

	if err := run(*filename); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

func run(filename string) error {
	// Read all the data from the provided input file and check for potential read errors.
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseContent(fileContent)

	outFName := fmt.Sprintf("%s.html", filepath.Base(filename))
	fmt.Println(outFName)

	return saveHTML(outFName, htmlData)
}

func parseContent(input []byte) []byte {
	// Parse the markdown input to generate valid and safe HTML.
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// Create a buffer of bytes to write to file.
	var buffer bytes.Buffer

	// Write HTML to this buffer including the header and footer constants.
	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	return buffer.Bytes()
}

func saveHTML(filename string, data []byte) error {
	// Write the data contents to file.
	return os.WriteFile(filename, data, 0644)
}