package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

const (
	port = 8443
)

func main() {
	handler := http.NewServeMux()

	handler.HandleFunc("/metadata/endpoints", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		apiVersion := q.Get("api-version")

		metadata, status, err := loadMetadataFromFile(apiVersion)
		if err != nil {
			w.WriteHeader(status)
			fmt.Fprintf(w, `{"error":"%s", "api-version": "%s"}`, err, apiVersion)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(metadata)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: handler,
	}

	log.Printf("Simple Azure Metadata Service listening on 127.0.0.1:%d\n", port)
	if err := server.ListenAndServeTLS(path.Join("tls", "server.crt"), path.Join("tls", "server.key")); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server.ListenAndServeTLS: %v", err)
	}
}

func loadMetadataFromFile(apiVersion string) ([]byte, int, error) {
	filename := ""

	switch apiVersion {
	case "1.0", "2015-01-01":
		filename = "metadata20150101.json"
	case "2018-01-01":
		filename = "metadata20180101.json"
	case "2019-05-01", "2020-06-01":
		filename = "metadata20190501.json"
	case "2022-09-01":
		filename = "metadata20220901.json"
	default:
		return nil, http.StatusBadRequest, fmt.Errorf("unrecognized api-version")
	}

	metadata, err := os.ReadFile(path.Join("metadata", filename))
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("reading metadata from file: %v", err)
	}

	return metadata, http.StatusOK, nil
}
