package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	// Define command-line flags for input file and output folder.
	inputFile := flag.String("input", "", "Input file containing list of URLs")
	outputDir := flag.String("output", "", "Output folder to save downloaded files")
	flag.Parse()

	// Check if arguments were provided as positional arguments instead of flags
	if *inputFile == "" && *outputDir == "" && len(flag.Args()) == 2 {
		*inputFile = flag.Args()[0]
		*outputDir = flag.Args()[1]
	}

	// Ensure both input and output are provided.
	if *inputFile == "" || *outputDir == "" {
		fmt.Println("Usage: downloadFilesFromList -input=<inputfile> -output=<outputfolder>")
		fmt.Println("   or: downloadFilesFromList <inputfile> <outputfolder>")
		os.Exit(1)
	}

	// Ensure output directory exists (create it if it doesn't).
	err := os.MkdirAll(*outputDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Open the input file.
	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Read the file line by line.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// Download each URL.
		err = downloadFile(line, *outputDir)
		if err != nil {
			fmt.Printf("Failed to download %s: %v\n", line, err)
		} else {
			fmt.Printf("Downloaded: %s\n", line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}
}

// downloadFile fetches the content from fileURL and saves it into outputDir.
func downloadFile(fileURL, outputDir string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("http get error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Extract the file name from the URL.
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return fmt.Errorf("error parsing url: %v", err)
	}
	filename := path.Base(parsedURL.Path)
	if filename == "/" || filename == "." || filename == "" {
		// Default to "index.html" if the URL doesn't specify a file.
		filename = "index.html"
	}

	// Build the full file path and create the file.
	filePath := filepath.Join(outputDir, filename)
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer outFile.Close()

	// Write the content to the file.
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}
