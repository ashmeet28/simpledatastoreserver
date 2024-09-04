package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	var dataStore []byte = []byte{0x7b, 0x7d}
	var dataStoreCounter int
	var dataStoreMu sync.Mutex

	http.HandleFunc("GET /datastoreget", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Methods", "*")
		w.Header().Add("Access-Control-Expose-Headers", "*")

		if r.Method == http.MethodGet && r.URL.Path == "/datastoreget" {
			dataStoreMu.Lock()

			w.Header().Add("sdss-data-store-counter", strconv.FormatInt(int64(dataStoreCounter), 10))
			w.WriteHeader(http.StatusOK)
			w.Write(dataStore)

			dataStoreMu.Unlock()
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	http.HandleFunc("OPTIONS /datastoreget", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Methods", "*")
		w.Header().Add("Access-Control-Expose-Headers", "*")

		if r.URL.Path == "/datastoreget" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	http.HandleFunc("POST /datastoreset", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Methods", "*")
		w.Header().Add("Access-Control-Expose-Headers", "*")

		if r.URL.Path == "/datastoreset" {
			c, err1 := strconv.ParseInt(r.Header.Get("sdss-data-store-counter"), 10, 64)
			if err1 != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			newDataStore, err2 := io.ReadAll(r.Body)
			if err2 != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			dataStoreMu.Lock()

			if c == ((int64(dataStoreCounter) + 1) % 1000000000) {
				dataStoreCounter = int(c)
				dataStore = newDataStore
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}

			dataStoreMu.Unlock()
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	http.HandleFunc("OPTIONS /datastoreset", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Methods", "*")
		w.Header().Add("Access-Control-Expose-Headers", "*")

		if r.URL.Path == "/datastoreset" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
