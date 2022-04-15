package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

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
	skipPreview := flag.Bool("s", false, "Skip auto file preview with the default browser")
	flag.Parse()

	if *filename == "" {
		flag.Usage()

		os.Exit(1)
	}

	if err := run(*filename, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

func run(filename string, w io.Writer, skipPreview bool) error {
	// Read all the data from the provided input file and check for potential read errors.
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseContent(fileContent)

	// Create a permanent file in the current root dir with the generated path.
	// outFName :=  fmt.Sprintf("%s.html", filepath.Base(filename))

	// Create a temp file in a temp dir based on the current os.
	tempF, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}

	if err := tempF.Close(); err != nil {
		return err
	}

	outFName := tempF.Name()
	fmt.Fprintln(w, outFName)

	if err := saveHTML(outFName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	return preview(outFName)
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

func preview(fileName string) error {
	commandName := ""
	commandParams := []string{}

	// Define executable based on os.
	switch runtime.GOOS {
	case "linux":
		commandName = "xdg-open"
	case "windows":
		commandName = "cmd.exe"
		commandParams = []string{"/c", "start"}
	case "darwin":
		commandName = "open"
	default:
		return fmt.Errorf("Os not supported")
	}

	// Append filename to params slice.
	commandParams = append(commandParams, fileName)

	// Locate executable based on path.
	commandPath, err := exec.LookPath(commandName)
	if err != nil {
		return err
	}

	// Open the file using the default program.
	return exec.Command(commandPath, commandParams...).Run()

}
