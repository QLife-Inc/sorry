package main

import "os"
import "github.com/QLife-Inc/sorry/lib"

var (
	Version  string
	Revision string
)

func getHttpPort() string {
	var listenPort = os.Getenv("PORT")
	if listenPort == "" {
		return "80"
	}
	return listenPort
}

func getHttpsPort() string {
	var listenPort = os.Getenv("HTTPS_PORT")
	if listenPort == "" {
		return "443"
	}
	return listenPort
}

func startHttpsListener() {
	var certs = lib.GetCertificatePairs()
	if len(certs) == 0 {
		return
	}
	var listener = lib.CreateHttpsListener(getHttpsPort(), certs)
	var server = lib.CreateServer()
	server.Serve(listener)
}

func startHttpListener() {
	var server = lib.CreateServer()
	server.Addr = ":" + getHttpPort()
	server.ListenAndServe()
}

func main() {
	finish := make(chan bool)

	go func() {
		startHttpListener()
	}()

	go func() {
		startHttpsListener()
	}()

	<-finish
}
