package main

import (
	"log"
	"os"

	"github.com/QLife-Inc/sorry/lib"
)

var (
	Version  string
	Revision string
)

func getEnvOrDefault(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getHttpPort() string {
	return getEnvOrDefault("PORT", "80")
}

func getHttpsPort() string {
	return getEnvOrDefault("HTTPS_PORT", "443")
}

func createHttpsListener(contents *lib.ResponseContents) (func(), error) {
	certs, err := lib.GetCertificatePairs()
	if err != nil {
		return nil, err
	}
	if len(certs) == 0 {
		return nil, nil
	}
	var port = getHttpsPort()
	listener, err := lib.CreateHttpsListener(port, certs)
	if err != nil {
		return nil, err
	}
	return func() {
		log.Printf("Start HTTPS Listener on port %s\n", port)
		lib.CreateServer(contents, "https").Serve(listener)
	}, nil
}

func createHttpListener(contents *lib.ResponseContents) func() {
	var server = lib.CreateServer(contents, "http")
	var port = getHttpPort()
	server.Addr = ":" + port
	return func() {
		log.Printf("Start HTTP Listener on port %s\n", port)
		server.ListenAndServe()
	}
}

func createListeners(contents *lib.ResponseContents) (func(), func(), error) {
	httpListener := createHttpListener(contents)
	httpsListener, err := createHttpsListener(contents)
	return httpListener, httpsListener, err
}

func main() {
	contents, err := lib.NewResponseContents()
	http, https, err := createListeners(contents)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	finish := make(chan bool)
	go http()
	go https()
	<-finish
}
