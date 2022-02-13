package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"strconv"
)

func main() {
	db := database{"shoes": 50, "socks": 5}
	mux := http.NewServeMux()
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)
	mux.HandleFunc("/update", db.update)
	mux.HandleFunc("/create", db.create)
	mux.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

// Responds with the items in the map and their prices
func (db database) list(w http.ResponseWriter, req *http.Request) {
	var mu sync.RWMutex
	mu.RLock()
	defer mu.RUnlock()

	// Check to see if the request is a GET
	if req.Method == "GET" {
		// send a list of the items in the database
		for item, price := range db {
			fmt.Fprintf(w, "%s: %s\n", item, price)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest) // If the request is not a GET then respond with a bad request
		fmt.Fprintf(w, "Error: Bad Request\n")
	}
}

// Responds with the price of the requested item
func (db database) price(w http.ResponseWriter, req *http.Request) {
	var mu sync.RWMutex
	mu.RLock()
	defer mu.RUnlock()

	// Check to see if the request is a GET
	if req.Method == "GET" {
		item := req.URL.Query().Get("item")

		// Send the price of the requested item if it exists. If not then send an error not found
		if price, ok := db[item]; ok {
			fmt.Fprintf(w, "%s\n", price)
		} else {
			w.WriteHeader(http.StatusNotFound) // 404
			fmt.Fprintf(w, "(Error) No such item: %q\n", item)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest) // If the request is not a GET then respond with a bad request
		fmt.Fprintf(w, "Error: Bad Request\n")
	}
}

// Creates an element in the db map
func (db database) create(w http.ResponseWriter, req *http.Request) {
	var mu sync.RWMutex
	mu.RLock()
	defer mu.RUnlock()

	// Check to see if the request is a POST and not a GET
	if req.Method == "POST" {
			item := req.URL.Query().Get("item")
			price,_ := strconv.ParseFloat(req.URL.Query().Get("price"), 32)

			// Adding the item to the database
			db[item] = dollars(price)
	} else {
		w.WriteHeader(http.StatusBadRequest) // If the request is not a POST then respond with a bad request
		fmt.Fprintf(w, "Error: Bad Request\n")
	}
}

// Updates an item in the db map
func (db database) update(w http.ResponseWriter, req *http.Request) {
	var mu sync.RWMutex
	mu.RLock()
	defer mu.RUnlock()

	// Check to see if the request is a POST and not a GET
	if req.Method == "POST" {
		item := req.URL.Query().Get("item")
		price,_ := strconv.ParseFloat(req.URL.Query().Get("price"), 32)

		// If the item exists then update the price else respond with a not found
		if _, ok := db[item]; ok {
			db[item] = dollars(price)
		} else {
			w.WriteHeader(http.StatusNotFound) // The item does not exist
			fmt.Fprintf(w, "(Error) No such item: %q\n", item)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest) // If the request is not a POST then respond with a bad request
		fmt.Fprintf(w, "Error: Bad Request\n")
	}
}

// Deletes an item in the db map
func (db database) delete(w http.ResponseWriter, req *http.Request) {
	fmt.Println("here")
	var mu sync.RWMutex
	mu.RLock()
	defer mu.RUnlock()

	if req.Method == "POST" {
		item := req.URL.Query().Get("item")
		// If the item exists then delete it
		if _, ok := db[item]; ok {
			delete(db,item)
		} else {
			w.WriteHeader(http.StatusNotFound) // The item does not exist
			fmt.Fprintf(w, "(Error) No such item: %q\n", item)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest) // If the request is not a POST then respond with a bad request
		fmt.Fprintf(w, "Error: Bad Request\n")
	}
}
