package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Health struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

var produkList = []Produk{
	{ID: 1, Nama: "Indomie Godog", Harga: 3500, Stok: 10},
	{ID: 2, Nama: "Vit 1000ml", Harga: 3000, Stok: 40},
	{ID: 3, Nama: "kecap", Harga: 12000, Stok: 20},
}

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{
			"status":  "OK",
			"message": "API Running",
		}
		writeJSON(w, http.StatusOK, resp)
	})
	// GET/POST localhost:8080/api/produk
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, http.StatusOK, produkList)
		case http.MethodPost:
			var produkBaru Produk
			if err := json.NewDecoder(r.Body).Decode(&produkBaru); err != nil {
				http.Error(w, "invalid request body", http.StatusBadRequest)
				return
			}
			if len(produkList) > 0 {
				produkBaru.ID = produkList[len(produkList)-1].ID + 1
			} else {
				produkBaru.ID = 1
			}
			produkList = append(produkList, produkBaru)
			writeJSON(w, http.StatusCreated, produkBaru)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	// GET/PUT/DELETE localhost:8080/api/produk/123
	http.HandleFunc("/api/produk/", handleProdukByID)

	fmt.Println("Server running di localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}

func handleProdukByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid produk id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		for _, p := range produkList {
			if p.ID == id {
				writeJSON(w, http.StatusOK, p)
				return
			}
		}
		http.Error(w, "produk belum ada", http.StatusNotFound)
	case http.MethodPut:
		var produkUpdate Produk
		if err := json.NewDecoder(r.Body).Decode(&produkUpdate); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		for i, p := range produkList {
			if p.ID == id {
				produkUpdate.ID = id
				produkList[i] = produkUpdate
				writeJSON(w, http.StatusOK, produkUpdate)
				return
			}
		}
		http.Error(w, "produk belum ada", http.StatusNotFound)
	case http.MethodDelete:
		for i, p := range produkList {
			if p.ID == id {
				// bikin slice baru dengan data sebelum dan sesudah index
				produkList = append(produkList[:i], produkList[i+1:]...)
				writeJSON(w, http.StatusOK, p)
				return
			}
		}
		http.Error(w, "produk belum ada", http.StatusNotFound)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func parseID(path string) (int, error) {
	// Parse ID dari URL path
	// URL: /api/produk/123 -> ID = 123
	idPart := strings.TrimPrefix(path, "/api/produk/")
	if idPart == "" {
		return 0, strconv.ErrSyntax
	}
	return strconv.Atoi(idPart)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
