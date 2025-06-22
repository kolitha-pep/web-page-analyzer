# Web Analyzer

**Web Analyzer** is a fast, concurrent web application built in Go using the Gin framework. It analyzes a given web page and returns useful metadata such as HTML version, page title, heading structure, link details, and login form detection.

## Features
- Determine HTML version
- Extract page title
- Analyze heading structure (H1, H2, H3)
- Count links and categorize them (internal, external)
- Determine inaccessible links
- Detect login forms
- Concurrent link checking
- Unit tests coverage over 70% (mocked HTTP requests)

## Technologies Used
- Go (Golang)
- Gin framework
- Goquery for HTML parsing
- Logrus for logging
- Docker for containerization

## Running the application locally

### Prerequisites

- Go 1.23+
- (Optional) Docker

### Build & Run (Go)

```bash
go mod tidy
go run cmd/server/main.go
```

### Build & Run (Makefile)
```bash
make run
```

### Build & Run (Makefile + Docker)

```bash
make docker-build
make docker-run
```

### Accessing the Application
Open your web browser and navigate to `http://localhost:8080/web`.

## Running Tests

```bash
go test ./... -cover
```

## Sample Usage

**Web Analyzer** provides a simple web interface to analyze any URL. You can enter a URL in the input field and click "Analyze" to get the metadata.

Example for heavy web page:
http://localhost:8080/web/?url=https%3A%2F%2Fgo.dev%2Fref%2Fspec

Example for light web page:
http://localhost:8080/web/?url=https%3A%2F%2Fgoogle.com

or, for a fresh start: http://localhost:8080/web