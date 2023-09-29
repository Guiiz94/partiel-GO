package scanner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"sync"
)

const (
	baseURL    = "http://10.49.122.144:"
	startPort  = 1
	endPort    = 65535
	timeout    = 1 * time.Second
	concurrent = 1000
)

type ScanConfig struct {
    StartPort   int
    EndPort     int
    Concurrency int
    Client      *http.Client
}

type PortScanner struct {
    Client      *http.Client
    StartPort   int
    EndPort     int
    Concurrency int
    wg          sync.WaitGroup
}

type HTTPClientConfig struct {
    Timeout time.Duration
    BaseURL string
}



func New(config ScanConfig) *PortScanner {
    return &PortScanner{
        Client:      config.Client,
        StartPort:   config.StartPort,
        EndPort:     config.EndPort,
        Concurrency: config.Concurrency,
    }
}


func NewPortScanner() *PortScanner {
	return &PortScanner{
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (ps *PortScanner) handlePort(port int) {
    resp, err := ps.Client.Get(fmt.Sprintf("%s%d/ping", baseURL, port))
    if err != nil {
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        fmt.Printf("Réponse réussie du port %d (statut %d): %s\n", port, resp.StatusCode, string(body))

        jsonData := `{
            "User": "Matteo"
        }`

        endpoints := []string{"/signup", "/check", "/getUserSecret", "/getUserLevel", "/getUserPoints"}
        for _, endpoint := range endpoints {
            req, err := http.NewRequest("POST", fmt.Sprintf("%s%d%s", baseURL, port, endpoint), strings.NewReader(jsonData))
            if err != nil {
                fmt.Printf("Erreur lors de la création de la requête POST (%s) pour le port %d: %s\n", endpoint, port, err)
                continue
            }

            req.Header.Set("Content-Type", "application/json")
            postResp, err := ps.Client.Do(req)
            if err != nil {
                fmt.Printf("Erreur lors de l'envoi de la requête POST (%s) sur le port %d: %s\n", endpoint, port, err)
                continue
            }
            defer postResp.Body.Close()

            responseBody, _ := ioutil.ReadAll(postResp.Body)
            fmt.Printf("Réponse de la requête POST (%s) depuis le port %d: %s\n", endpoint, port, string(responseBody))

            if endpoint == "/getUserSecret" {
                const prefix = "User secret: "
                if strings.HasPrefix(string(responseBody), prefix) {
                    secret := strings.TrimPrefix(string(responseBody), prefix)
                    secret = strings.TrimSpace(secret)
                    fmt.Println("Extracted secret:", secret)
                    
                    data := map[string]string{
                        "User":   "Matteo",
                        "Secret": secret,
                    }
                    jsonDataBytes, _ := json.Marshal(data)
                    jsonData = string(jsonDataBytes)
                } else {
                    fmt.Println("Unexpected format for getUserSecretBody")
                }
            }
        }
    }
}


func (ps *PortScanner) ScanPorts() {
	var wg sync.WaitGroup
	ports := make(chan int, concurrent)

	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range ports {
				resp, err := ps.Client.Get(fmt.Sprintf("%s%d/ping", baseURL, port))
				if err == nil && resp.StatusCode == http.StatusOK {
					ps.handlePort(port)
					return 
				}
			}
		}()
	}

	for port := startPort; port <= endPort; port++ {
		ports <- port
	}
	close(ports)

	wg.Wait()
}

func NewHTTPClient(config HTTPClientConfig) *http.Client {
    return &http.Client{
        Timeout: config.Timeout,
    }
}

func (ps *PortScanner) Start() {
    var wg sync.WaitGroup
    ports := make(chan int, ps.Concurrency)

    for i := 0; i < ps.Concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for port := range ports {
                ps.handlePort(port) 
            }
        }()
    }

    for port := ps.StartPort; port <= ps.EndPort; port++ {
        ports <- port
    }
    close(ports)

    wg.Wait()
}
