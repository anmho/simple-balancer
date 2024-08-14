package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	var port int
	if len(os.Args) < 2 {
		log.Fatalln("please provide port to run on")
	}
	x, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalln("invalid port value", os.Args[1])
	}
	port = x

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": fmt.Sprintf("Hello World from %d", port),
		})
	})

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	log.Printf("starting on port %d\n", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln("failed to start server", err)
	}

}
