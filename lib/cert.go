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

func (p CertificatePair) ToCertificate() tls.Certificate {
	if cert, err := tls.LoadX509KeyPair(p.chainCert, p.privateKey); err != nil {
		panic(err)
	} else {
		return cert
	}
}

func GetCertificatePairs() []CertificatePair {
	if dir, err := os.Stat(sslDir); os.IsExist(err) || !dir.IsDir() {
		return []CertificatePair{}
	}
	if children, err := ioutil.ReadDir(sslDir); err != nil {
		panic(err)
	} else {
		return walkDirs(children)
	}
}

func walkDirs(dirs []os.FileInfo) []CertificatePair {
	var pairs []CertificatePair
	for _, child := range dirs {
		if !child.IsDir() {
			continue
		}
		if pair := getCertificatePair(child); pair != nil {
			pairs = append(pairs, *pair)
		}
	}
	return pairs
}

func getCertificatePair(dir os.FileInfo) *CertificatePair {
	path := filepath.Join(sslDir, dir.Name())
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	var cert string = ""
	var key string = ""
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), certExt) {
			cert = filepath.Join(path, file.Name())
			if key != "" {
				return &CertificatePair{
					chainCert: cert,
					privateKey: key,
				}
			}
			continue
		}
		if strings.HasSuffix(file.Name(), keyExt) {
			key = filepath.Join(path, file.Name())
			if cert != "" {
				return &CertificatePair{
					chainCert: cert,
					privateKey: key,
				}
			}
			continue
		}
	}
	return nil
}
