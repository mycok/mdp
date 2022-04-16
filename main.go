package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"
	"html/template"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

type content struct {
	Title string
	Body template.HTML
}

func main() {
	filename := flag.String("file", "", "Markdown file to preview")
	templateFName := flag.String("t", "", "Name of the template file provided by the user")
	skipPreview := flag.Bool("s", false, "Skip auto file preview with the default browser")
	flag.Parse()

	if *filename == "" {
		flag.Usage()

		os.Exit(1)
	}

	if err := run(*filename,  *templateFName, *skipPreview, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

func run(filename, templateFName string, skipPreview bool, w io.Writer) error {
	// Read all the data from the provided input file and check for potential read errors.
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData, err := parseContent(fileContent, templateFName)
	if err != nil {
		return err
	}

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

	defer os.Remove(outFName)

	return preview(outFName)
}

func parseContent(input []byte, templateFName string) ([]byte, error) {
	// Parse the markdown input to generate valid and safe HTML.
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	tmpl, err := template.ParseFiles("default.html.tmpl")
	if err != nil {
		return nil, err
	}

	// If a user provides a custom template, we create a new template
	// with parsed content from the user provided template
	if templateFName != "" {
		tmpl, err = template.ParseFiles(templateFName)
		if err != nil {
			return nil, err
		}
	}

	c := content{
		Title: "Markdown Preview Tool",
		Body: template.HTML(body),
	}


	// Create a buffer of bytes to write to file.
	var buffer bytes.Buffer

	if err = tmpl.Execute(&buffer, c); err != nil {
		return nil, err
	}


	return buffer.Bytes(), nil
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
	err = exec.Command(commandPath, commandParams...).Run()

	// Add a delay to give the browser time to open the file before returning
	// from the function. Once the function returns, the calling function calls all
	// pending defer statements which in this case deletes the mentioned preview file.
	// TODO: replace the sleep delay functionality with signal or some better way.
	time.Sleep(5 * time.Second)

	return err
}
