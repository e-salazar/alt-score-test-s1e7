package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
)

var (
	systems = map[string]string{ // [key: system_name] = value: system_code
		"navigation":       "NAV-01",
		"communications":   "COM-02",
		"life_support":     "LIFE-03",
		"engines":          "ENG-04",
		"deflector_shield": "SHLD-05",
	}
	lastDamagedSystem string
	mutex             sync.RWMutex
)

type StatusResponse struct {
	DamagedSystem string `json:"damaged_system"`
}

func main() {
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/repair-bay", repairBayHandler)
	http.HandleFunc("/teapot", teapotHandler)

	fmt.Println("API ejecutándose en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Endpoint /status
func statusHandler(httpResponse http.ResponseWriter, httpRequest *http.Request) {
	if httpRequest.Method != http.MethodGet {
		http.Error(httpResponse, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	damagedSystem := getRandomSystem()
	setLastDamagedSystem(damagedSystem)

	response := StatusResponse{
		DamagedSystem: damagedSystem,
	}
	httpResponse.Header().Set("Content-Type", "application/json")
	json.NewEncoder(httpResponse).Encode(response)
}

// Endpoint /repair-bay
func repairBayHandler(httpResponse http.ResponseWriter, httpRequest *http.Request) {
	if httpRequest.Method != http.MethodGet {
		http.Error(httpResponse, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	damagedSystem := getLastDamagedSystem()

	if damagedSystem == "" {
		http.Error(httpResponse, "No hay sistema averiado registrado. Primero llama al endpoint /status.", http.StatusNotFound)
		return
	}

	systemCode := systems[damagedSystem]

	htmlContent := fmt.Sprintf(`
        <!DOCTYPE html>
        <html>
        <body>
            <div class="anchor-point">%s</div>
        </body>
        </html>
    `, systemCode)

	httpResponse.Header().Set("Content-Type", "text/html")
	httpResponse.Write([]byte(htmlContent))
}

// Endpoint /teapot
func teapotHandler(httpResponse http.ResponseWriter, httpRequest *http.Request) {
	if httpRequest.Method != http.MethodPost {
		http.Error(httpResponse, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	httpResponse.WriteHeader(http.StatusTeapot)
	httpResponse.Write([]byte("I'm a teapot"))
}

func getRandomSystem() string {
	keys := make([]string, 0, len(systems))
	for key := range systems {
		keys = append(keys, key)
	}
	return keys[rand.Intn(len(keys))]
}

func setLastDamagedSystem(system string) {
	mutex.Lock()
	lastDamagedSystem = system
	mutex.Unlock()
}
func getLastDamagedSystem() string {
	mutex.RLock()
	defer mutex.RUnlock()
	return lastDamagedSystem
}
