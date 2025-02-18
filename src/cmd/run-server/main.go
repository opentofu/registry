package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	destinationDir := flag.String("destination", "../generated", "Directory to read the generated responses")
	certFile := flag.String("certificate", "localhost.pem", "Path to certificate file")
	keyFile := flag.String("key", "localhost-key.pem", "Path to key file")

	flag.Parse()

	_, err := os.Stat(*destinationDir)
	if err != nil {
		logger.Error("Could not open folder", slog.Any("Error", err))
		os.Exit(1)
	}

	fs := http.FileServer(http.Dir(*destinationDir))
	http.Handle("/", fs)

	logger.Info("Starting HTTPS server at :8443...")

	err = http.ListenAndServeTLS(":8443", *certFile, *keyFile, nil)
	if err != nil {
		logger.Error("Failed to start server", slog.Any("Error", err))
		os.Exit(1)
	}
}
