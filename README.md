# UrlHarvester

A lightweight Go tool for downloading files from a list of URLs.

## Features

- Download multiple files from a list of URLs in a text file
- Custom user-agent to mimic browser requests
- Optional proxy support
- Simple command-line interface with both flag and positional argument support

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/UrlHarvester.git
cd UrlHarvester

# Build the binary
go build -o UrlHarvester app.go
```

Alternatively, download the pre-built binary from the releases page.

## Usage

UrlHarvester supports two syntax formats:

```bash
# Using flags
./UrlHarvester -input=urls.txt -output=downloads -proxy=127.0.0.1:8080

# Using positional arguments
./UrlHarvester urls.txt downloads
```

### Parameters

- `input`: Path to a text file containing one URL per line
- `output`: Directory where downloaded files will be saved
- `proxy` (optional): Proxy server address in the format host:port

## Example

Create a file named `urls.txt` with content:

```
https://example.com/file1.pdf
https://example.com/images/logo.png
https://another-site.com/document.docx
```

Run the tool:

```bash
./UrlHarvester -input=urls.txt -output=downloads
```

Files will be saved to the `downloads` directory, preserving their original filenames.

## License

MIT License
