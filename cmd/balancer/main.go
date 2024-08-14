package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"net/url"
	"os"
	backend "simple-balancer/balancer"
)

var (
	config     Config
	serverPool = backend.NewServerPool()
)

func loadBalance(w http.ResponseWriter, r *http.Request) {
	peer := serverPool.NextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}

	http.Error(w, "Service not available", http.StatusServiceUnavailable)

	return
}

type Config struct {
	Backends map[string][]string `yaml:"backends"`
}

func main() {

	var configPath string
	var port int

	// Parse flags
	configPath = *flag.String("config", "./config.yaml", "")
	port = *flag.Int("port", 8080, "")
	flag.Parse()

	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalln("error opening the file")
	}
	fmt.Println("args", os.Args)
	fmt.Println("path", configPath)

	// Read the values into the config struct
	yamlDecoder := yaml.NewDecoder(file)
	if err := yamlDecoder.Decode(&config); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Print the parsed result
	fmt.Printf("Parsed YAML into struct: %+v\n", config)

	for path, backends := range config.Backends {
		for i, be := range backends {
			u, err := url.Parse(be)
			if err != nil {
				log.Fatalf("malformed backend url at index %d: %s\n", i, path)
			}
			serverPool.RegisterBackend(u)
		}
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(loadBalance),
	}

	log.Println("listening on", port)
	err = srv.ListenAndServe()
	if err != nil {
		return
	}
}
