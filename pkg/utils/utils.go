package utils

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
)

func ToHexInt(n *big.Int) string {
	return fmt.Sprintf("%x", n) // or %X or upper case
}

func InsertNth(s string, n int) string {
	if len(s)%2 != 0 {
		s = "0" + s
	}
	var buffer bytes.Buffer
	var n_1 = n - 1
	var l_1 = len(s) - 1
	for i, rune := range s {
		buffer.WriteRune(rune)
		if i%n == n_1 && i != l_1 {
			buffer.WriteRune('-')
		}
	}
	return buffer.String()
}
func ReadCertPool(path string) (*x509.CertPool, error) {
	caCert, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return caCertPool, nil
}
func ReadCert(path string) (*x509.Certificate, error) {
	certContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cpb, _ := pem.Decode(certContent)

	crt, err := x509.ParseCertificate(cpb.Bytes)
	if err != nil {
		return nil, err
	}
	return crt, nil
}

func ReadKey(path string) ([]byte, error) {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return key, nil
}
