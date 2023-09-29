# partiel-GO

                    // Convert the map to JSON
					jsonBody, err := json.Marshal(data)
					if err != nil {
						fmt.Println("Error encoding JSON:", err)
						return
					}

					// Envoi de la requête POST vers /iNeedAHint
					reqHint, err := http.NewRequest("POST", fmt.Sprintf("%s%d/iNeedAHint", baseURL, port), bytes.NewBuffer(jsonBody))
					if err != nil {
						fmt.Printf("Erreur lors de la création de la requête POST (iNeedAHint) pour le port %d: %s\n", port, err)
						return 
					}
					reqHint.Header.Set("Content-Type", "application/json")
					respHint, err := client.Do(reqHint)
					if err != nil {
						fmt.Printf("Erreur lors de l'envoi de la requête POST (iNeedAHint) sur le port %d: %s\n", port, err)
						return 
					}
					defer respHint.Body.Close()
					bodyHint, _ := ioutil.ReadAll(respHint.Body)
					fmt.Printf("Réponse de la requête POST (iNeedAHint) depuis le port %d: %s\n", port, string(bodyHint))
				