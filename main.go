package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	currentName string
	mu          sync.Mutex // Para evitar race conditions
)

func main() {
	// 1. Leer valor inicial de la variable de entorno
	currentName = os.Getenv("NOMBRE")
	if currentName == "" {
		currentName = "Desconocido" // Valor por defecto
	}

	fmt.Println("Nombre inicial:", currentName)

	// 2. Configurar endpoints
	http.HandleFunc("/nombre", nombreHandler)

	// 3. Iniciar servidor
	fmt.Println("Servidor escuchando en :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func nombreHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mu.Lock()
		defer mu.Unlock()
		json.NewEncoder(w).Encode(map[string]string{"nombre": currentName})

	case http.MethodPut:
		var request struct {
			Nombre string `json:"nombre"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
			return
		}

		mu.Lock()
		currentName = request.Nombre
		mu.Unlock()

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Nombre actualizado a: %s", currentName)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
