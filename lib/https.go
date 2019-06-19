package lib

import (
	"crypto/tls"
	"net"
)

func createTlsConfig(certPairs []CertificatePair) *tls.Config {
	config := &tls.Config{}
	config.Certificates = make([]tls.Certificate, len(certPairs))
	for i, cert := range certPairs {
		config.Certificates[i] = cert.ToCertificate()
	}
	config.BuildNameToCertificate()
	return config
}

func CreateHttpsListener(port string, certPairs []CertificatePair) net.Listener {
	config := createTlsConfig(certPairs)
	if listener, err := tls.Listen("tcp", ":" + port, config); err != nil {
		panic(err)
	} else {
		return listener
	}
}
