package main

import (
	"microgo/handlers"
	"log"
	"net/http"
)

func main() {
	h := &handlers.ChartHandler{}

	// Маршруты
	http.HandleFunc("/api/chart/sync", h.HandleSyncChart)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Оно живое"))
	})

	log.Println("Стартуем на порте :8080")
	
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}