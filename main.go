package main

import (
	"time"
	"github.com/Guiiz94/partiel-GO/scanner"
)

const (
    baseURL    = "http://10.49.122.144:"
    startPort  = 1
    endPort    = 65535
    timeout    = 1 * time.Second
    concurrent = 1000
)

func main() {
    httpClientConfig := scanner.HTTPClientConfig{
        Timeout: timeout,
        BaseURL: baseURL,
    }
    client := scanner.NewHTTPClient(httpClientConfig)

    scanConfig := scanner.ScanConfig{
        StartPort:   startPort,
        EndPort:     endPort,
        Concurrency: concurrent,
        Client:      client,
    }
    scannerInstance := scanner.New(scanConfig)
    scannerInstance.Start()
}