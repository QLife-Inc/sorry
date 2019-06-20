package lib

import (
	"crypto/tls"
	"net"
)

func createTlsConfig(certPairs []CertificatePair) (*tls.Config, error) {
	config := &tls.Config{}
	config.Certificates = make([]tls.Certificate, len(certPairs))
	for i, pair := range certPairs {
		if cert, err := pair.ToCertificate(); err != nil {
			return nil, err
		} else {
			config.Certificates[i] = *cert
		}
	}
	config.BuildNameToCertificate()
	return config, nil
}

func CreateHttpsListener(port string, certPairs []CertificatePair) (net.Listener, error) {
	if config, err := createTlsConfig(certPairs); err != nil {
		return nil, err
	} else {
		return tls.Listen("tcp", ":" + port, config)
	}
}
