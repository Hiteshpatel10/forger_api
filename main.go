package main

import (
	"forger/db"
	flutterforge "forger/flutter_forge"
	forgeicons "forger/forge_icons"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	db.Init()
	defer func() {
		if err := db.Database.Close(); err != nil {
			log.Fatalf("Error closing database: %v", err)
		}
	}()
	router := mux.NewRouter()

	router.Use(corsMiddleware)
	router.HandleFunc("/forge", flutterforge.ForgeCategory)
	router.HandleFunc("/components/{slug}", flutterforge.ForgeComponents)

	// Apply CORS middleware
	router.HandleFunc("/icons", forgeicons.GetForgeIcons)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

// CORS middleware function
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If the request is an OPTIONS request, we just return OK and exit
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
