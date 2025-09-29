package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

var users = map[int]string{
	1: "Anelya",
	2: "Sanzhar",
	3: "Mansur",
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getUser(w, r)
	} else if r.Method == http.MethodPost {
		createUser(w, r)
	} else {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	if _, ok := users[id]; ok {
		json.NewEncoder(w).Encode(map[string]int{"user_id": id})
	} else {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil || data.Name == "" {
		http.Error(w, `{"error":"invalid name"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"created": data.Name})
}
