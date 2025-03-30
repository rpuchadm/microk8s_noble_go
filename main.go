package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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
	http.HandleFunc("/fichero/", ficheroHandler)

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

// Handler para el endpoint /fichero/<<nombre_fichero>>
// nombre_fichero es un entero de 0 a 9
// cuando se hace un post se crea el fichero con el contenido que se pase en body
// cuando se hace un get se devuelve el contenido del fichero
// cuando se hace un delete se borra el fichero
// cuando se hace un put se actualiza el contenido del fichero

const base_path = "/var/data/ficheros/"

func ficheroHandler(w http.ResponseWriter, r *http.Request) {

	// Obtiene el ID del fichero
	path := strings.TrimPrefix(r.URL.Path, "/fichero/")
	id := strings.Split(path, "/")[0]

	// Validar que el nombre del fichero es un entero de 0 a 9
	if len(id) != 1 || id < "0" || id > "9" {
		http.Error(w, fmt.Sprintf("Nombre de fichero inválido %s", id), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, base_path+id)

	case http.MethodPost:
		file, err := os.Create(base_path + id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error al crear el fichero %v", err), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		if _, err := file.ReadFrom(r.Body); err != nil {
			http.Error(w, "Error al leer el contenido del body", http.StatusInternalServerError)
			return
		}

	case http.MethodDelete:
		if err := os.Remove(base_path + id); err != nil {
			http.Error(w, "Error al borrar el fichero", http.StatusInternalServerError)
			return
		}

	case http.MethodPut:
		file, err := os.OpenFile(base_path+id, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			http.Error(w, "Error al abrir el fichero", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		if _, err := file.ReadFrom(r.Body); err != nil {
			http.Error(w, "Error al leer el contenido del body", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
