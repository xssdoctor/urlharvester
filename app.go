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
	proxyURL := flag.String("proxy", "", "Proxy server to use (e.g., 127.0.0.1:8080)")
	flag.Parse()

	// Check if arguments were provided as positional arguments instead of flags
	if *inputFile == "" && *outputDir == "" && len(flag.Args()) == 2 {
		*inputFile = flag.Args()[0]
		*outputDir = flag.Args()[1]
	}

	// Ensure both input and output are provided.
	if *inputFile == "" || *outputDir == "" {
		fmt.Println("Usage: downloadFilesFromList -input=<inputfile> -output=<outputfolder> [-proxy=<proxyaddress>]")
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
		err = downloadFile(line, *outputDir, *proxyURL)
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
func downloadFile(fileURL, outputDir, proxyURL string) error {
	// Create a new request
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the User-Agent header
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36")

	// Create HTTP client with optional proxy
	client := &http.Client{}
	
	// Configure proxy if specified
	if proxyURL != "" {
		proxyURLParsed, err := url.Parse("http://" + proxyURL)
		if err != nil {
			return fmt.Errorf("error parsing proxy URL: %v", err)
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURLParsed),
		}
		client.Transport = transport
	}
	
	// Send the request
	resp, err := client.Do(req)
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
