package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
	"encoding/json"
	"bytes"
)

const (
	baseURL    = "http://10.49.122.144:"
	startPort  = 1
	endPort    = 65535
	timeout    = 1 * time.Second
	concurrent = 1000 
	maxAttempts = 50
)

func main() {
	client := &http.Client{
		Timeout: timeout,
	}

	var wg sync.WaitGroup
	ports := make(chan int, concurrent)

	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range ports {
				resp, err := client.Get(fmt.Sprintf("%s%d/ping", baseURL, port))
				if err == nil && resp.StatusCode == http.StatusOK {
					defer resp.Body.Close()
					body, _ := ioutil.ReadAll(resp.Body)
					fmt.Printf("Réponse réussie du port %d (statut %d): %s\n", port, resp.StatusCode, string(body))

					// Envoi de la requête POST après avoir identifié le bon port
					jsonData := `{
						"User": "Matteo"
					}`
					// Envoi de la requête POST vers /signup
					req, err := http.NewRequest("POST", fmt.Sprintf("%s%d/signup", baseURL, port), strings.NewReader(jsonData))
					if err != nil {
						fmt.Printf("Erreur lors de la création de la requête POST (signup) pour le port %d: %s\n", port, err)
						return
					}
					req.Header.Set("Content-Type", "application/json")
					signupResp, err := client.Do(req)
					if err != nil {
						fmt.Printf("Erreur lors de l'envoi de la requête POST (signup) sur le port %d: %s\n", port, err)
						return
					}
					defer signupResp.Body.Close()
					signupBody, _ := ioutil.ReadAll(signupResp.Body)
					fmt.Printf("Réponse de la requête POST (signup) depuis le port %d: %s\n", port, string(signupBody))

					// Envoi de la requête POST vers /check
					reqCheck, err := http.NewRequest("POST", fmt.Sprintf("%s%d/check", baseURL, port), strings.NewReader(jsonData))
					if err != nil {
						fmt.Printf("Erreur lors de la création de la requête POST (check) pour le port %d: %s\n", port, err)
						return
					}
					reqCheck.Header.Set("Content-Type", "application/json")
					check, err := client.Do(reqCheck)
					if err != nil {
						fmt.Printf("Erreur lors de l'envoi de la requête POST (check) sur le port %d: %s\n", port, err)
						return
					}
					defer check.Body.Close()
					checkBody, _ := ioutil.ReadAll(check.Body)
					fmt.Printf("Réponse de la requête POST (check) depuis le port %d: %s\n", port, string(checkBody))

					// Envoi de la requête POST vers /getUserSecret
					reqGetUserSecret, err := http.NewRequest("POST", fmt.Sprintf("%s%d/getUserSecret", baseURL, port), strings.NewReader(jsonData))
					if err != nil {
						fmt.Printf("Erreur lors de la création de la requête POST (getUserSecret) pour le port %d: %s\n", port, err)
						return
					}
					reqGetUserSecret.Header.Set("Content-Type", "application/json")
					getUserSecret, err := client.Do(reqGetUserSecret)
					if err != nil {
						fmt.Printf("Erreur lors de l'envoi de la requête POST (getUserSecret) sur le port %d: %s\n", port, err)
						return
					}
					defer getUserSecret.Body.Close()
					getUserSecretBody, _ := ioutil.ReadAll(getUserSecret.Body)
					fmt.Printf("Raw response from getUserSecret: %s\n", string(getUserSecretBody))

					const prefix = "User secret: "
					var secret string
					if strings.HasPrefix(string(getUserSecretBody), prefix) {
						secret = strings.TrimPrefix(string(getUserSecretBody), prefix)
						secret = strings.TrimSpace(secret)
						fmt.Println("Extracted secret:", secret)
					} else {
						fmt.Println("Unexpected format for getUserSecretBody")
						return // Si le format est inattendu, vous voudrez peut-être sortir de la boucle.
					}

					fmt.Printf("Secret pour le port %d: %s\n", port, secret)

					// Construction du JSON pour l'appel /getUserLevel
					data := map[string]string{
						"User":   "Matteo",
						"Secret": secret,
					}
					jsonLevelData, err := json.Marshal(data)
					if err != nil {
						fmt.Printf("Erreur lors de la conversion en JSON pour le port %d: %s\n", port, err)
						return
					}

					// Envoi de la requête POST vers /getUserLevel
					reqUserLevel, err := http.NewRequest("POST", fmt.Sprintf("%s%d/getUserLevel", baseURL, port), bytes.NewBuffer(jsonLevelData))
					if err != nil {
						fmt.Printf("Erreur lors de la création de la requête POST (getUserLevel) pour le port %d: %s\n", port, err)
						return
					}
					reqUserLevel.Header.Set("Content-Type", "application/json")
					respUserLevel, err := client.Do(reqUserLevel)
					if err != nil {
						fmt.Printf("Erreur lors de l'envoi de la requête POST (getUserLevel) sur le port %d: %s\n", port, err)
						return
					}
					defer respUserLevel.Body.Close()
					bodyUserLevel, _ := ioutil.ReadAll(respUserLevel.Body)
					fmt.Printf("Réponse de la requête POST (getUserLevel) depuis le port %d: %s\n", port, string(bodyUserLevel))

					// Envoi de la requête POST vers /getUserPoints
					reqUserPoints, err := http.NewRequest("POST", fmt.Sprintf("%s%d/getUserPoints", baseURL, port), bytes.NewBuffer(jsonLevelData))
					if err != nil {
						fmt.Printf("Erreur lors de la création de la requête POST (getUserPoints) pour le port %d: %s\n", port, err)
						return
					}
					reqUserPoints.Header.Set("Content-Type", "application/json")
					respUserPoints, err := client.Do(reqUserPoints)
					if err != nil {
						fmt.Printf("Erreur lors de l'envoi de la requête POST (getUserPoints) sur le port %d: %s\n", port, err)
						return
					}
					defer respUserPoints.Body.Close()
					bodyUserPoints, _ := ioutil.ReadAll(respUserPoints.Body)
					if err != nil {
						fmt.Printf("Erreur lors de la lecture du corps de la réponse pour le port %d: %s\n", port, err)
						continue
					}
					fmt.Printf("Réponse de la requête POST (getUserPoints) depuis le port %d: %s\n", port, string(bodyUserPoints))
				
				
					return
				}
			}
		}()
	}

	// Remplissage du canal avec les numéros de port
	for port := startPort; port <= endPort; port++ {
		ports <- port
	}
	close(ports)

	wg.Wait()
}

