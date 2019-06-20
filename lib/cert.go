package lib

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	certExt = ".crt"
	keyExt  = ".key"
	sslDir  = "ssl"
)

type CertificatePair struct {
	chainCert string
	privateKey string
}

func (p CertificatePair) String() string {
	return fmt.Sprintf("cert path = %s, key path = %s", p.chainCert, p.privateKey)
}

func (p CertificatePair) ToCertificate() (*tls.Certificate, error) {
	if cert, err := tls.LoadX509KeyPair(p.chainCert, p.privateKey); err != nil {
		return nil, err
	} else {
		return &cert, nil
	}
}

func GetCertificatePairs() ([]CertificatePair, error) {
	if dir, err := os.Stat(sslDir); os.IsExist(err) || !dir.IsDir() {
		return []CertificatePair{}, nil
	}
	if children, err := ioutil.ReadDir(sslDir); err != nil {
		return nil, err
	} else {
		return walkDirs(children)
	}
}

func walkDirs(dirs []os.FileInfo) ([]CertificatePair, error) {
	var pairs []CertificatePair
	for _, child := range dirs {
		if !child.IsDir() {
			continue
		}
		if pair, err := getCertificatePair(child); err != nil {
			return nil, err
		} else if pair != nil {
			pairs = append(pairs, *pair)
		}
	}
	return pairs, nil
}

func isCertificateFile(file os.FileInfo) bool {
	return strings.HasSuffix(file.Name(), certExt)
}

func isPrivateKeyFile(file os.FileInfo) bool {
	return strings.HasSuffix(file.Name(), keyExt)
}

func buildCertificatePair(cert string, key string) *CertificatePair {
	if cert == "" || key == "" {
		return nil
	}
	return &CertificatePair{
		chainCert: cert,
		privateKey: key,
	}
}

func getCertificatePair(dir os.FileInfo) (*CertificatePair, error) {
	path := filepath.Join(sslDir, dir.Name())
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, err
	}

	var cert string = ""
	var key string = ""

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if isCertificateFile(file) {
			cert = filepath.Join(path, file.Name())
			if pair := buildCertificatePair(cert, key); pair != nil {
				return pair, nil
			}
			continue
		}
		if isPrivateKeyFile(file) {
			key = filepath.Join(path, file.Name())
			if pair := buildCertificatePair(cert, key); pair != nil {
				return pair, nil
			}
			continue
		}
	}

	return nil, nil
}
