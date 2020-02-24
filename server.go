package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Entry into the DB
type Entry struct {
  Key string `json:"key"`
  Value interface{} `json:"value"`
}

var (
  db = map[string]interface{}{}
  dbLock sync.Mutex
)

func sendResponse(entry *Entry, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(entry); err != nil {
		log.Printf("error encoding %+v - %s", entry, err)
	}
}

func dbPostHandler(w http.ResponseWriter, r *http.Request) {
	// Decode request
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	entry := &Entry{}
	if err := dec.Decode(entry); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Do work
	dbLock.Lock()
	defer dbLock.Unlock()
	db[entry.Key] = entry.Value

	// Encode response
	sendResponse(entry, w)
}

func dbGetHandler(w http.ResponseWriter, r *http.Request) {
	// GET request /<key>
	key := r.URL.Path[4:] // Trim leading prefix /db/

	dbLock.Lock() // Acquire lock
	defer dbLock.Unlock() // Deferred call to unlock

	value, ok := db[key]
	if !ok {
		http.Error(w, fmt.Sprintf("Key %q not found", key), http.StatusNotFound)
		return
	}
  // If key is found
	entry := &Entry{
		Key:   key,
		Value: value,
	}
	sendResponse(entry, w)
}

func main() {
	http.HandleFunc("/db", dbPostHandler)
	http.HandleFunc("/db/", dbGetHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
