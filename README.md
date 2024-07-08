
# DateTime Client

This package provides a simple HTTP client for interacting with a datetime server. It's designed to retrieve the current date and time from a specified server endpoint.

## Features

- Easy-to-use client for fetching current date and time
- Configurable base URL and port
- Built-in retry mechanism with backoff strategy
- Timeout handling

## Installation

To use this package in your Go project, you can install it using:

```
go get github.com/codescalersinternships/datetime-client-eyadhussein
```

## Usage

Here's a basic example of how to use the DateTime Client:

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/codescalersinternships/datetime-client-eyadhussein/pkg/datetimeclient"
)

func main() {
    // Create a new client
    client := datetimeclient.NewRealClient("http://localhost", "8080", 10*time.Second)

    // Get the current date and time
    dateTime, err := client.GetCurrentDateTime()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Current Date and Time: %q\n", string(dateTime))
}
```

If environment variables are defined, just create a new client with empty string arguments:
```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/codescalersinternships/datetime-client-eyadhussein/pkg/datetimeclient"
)

func main() {
    // Create a new client
    client := datetimeclient.NewRealClient("", "", 10*time.Second)

    // Get the current date and time
    dateTime, err := client.GetCurrentDateTime()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Current Date and Time: %q\n", string(dateTime))
}
```

In the terminal, run:
```bash
SERVER_URL=http://localhost PORT=8080 go run main.go
```

Terminal output:
```bash
2024-07-04 15:11:44
```

# How to Test

Run

```bash
go test -v ./...
```
