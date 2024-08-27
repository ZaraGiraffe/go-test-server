package server

import (
	"net/http"
	"time"
	"unicode/utf8"
	"server/src/database"
	"io"
)


func index[T comparable](lst []T, elem T) (int, bool) {
	for i, val := range lst {
		if val == elem {
			return i, true
		}
	}
	return 0, false
}


func setHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	content := make([]byte, 256)
	n_bytes, err := r.Body.Read(content)
	if err != io.EOF {
		if err == nil {
			http.Error(w, "Too long input", http.StatusBadRequest)
			return
		} else {
			http.Error(w, "Bad body", http.StatusBadRequest)
			return
		}
	}
	
	if n_bytes >= 200 {
		http.Error(w, "The body is too big", http.StatusBadRequest)
		return 
	}
	content = content[:n_bytes]

	var seperator byte = 0x2C
	index, ok := index(content, seperator)
	if !ok {
		http.Error(w, "There is no seperator", http.StatusBadRequest)
		return 
	}
	if !utf8.Valid(content[:index]) || !utf8.Valid(content[index+1:]) {
		http.Error(w, "Input not in utf-8 format", http.StatusBadRequest)
		return 
	}

	key := string(content[:index])
	value := string(content[index+1:])
	database.SetKey(key, value)

	w.WriteHeader(http.StatusOK)
}


func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	content := make([]byte, 256)
	n_bytes, err := r.Body.Read(content)
	if err != io.EOF {
		if err == nil {
			http.Error(w, "Too long input", http.StatusBadRequest)
			return
		} else {
			http.Error(w, "Bad body", http.StatusBadRequest)
			return
		}
	}

	if n_bytes >= 200 {
		http.Error(w, "The body is too big", http.StatusBadRequest)
		return 
	}
	content = content[:n_bytes]

	if !utf8.Valid(content) {
		http.Error(w, "Input not in utf-8 format", http.StatusBadRequest)
		return
	}
	key := string(content)

	val, ok := database.GetKey(key)
	if ok != nil {
		http.Error(w, "No such key", http.StatusNotFound)
		return
	}

	response_content := []byte(val)
	_, err = w.Write(response_content)
	if err != nil {
		http.Error(w, "Can not write to response body", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}


func NewServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/get", getHandler)
	mux.HandleFunc("/set", setHandler)

	server := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  15 * time.Second,
    }

	return server
}