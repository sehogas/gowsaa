package afip

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"

	"go.mozilla.org/pkcs7"
)

// encodeCMS devuelve el content firmado PKCS#7
func encodeCMS(content []byte, certificate *x509.Certificate, privateKey *rsa.PrivateKey) ([]byte, error) {

	signedData, err := pkcs7.NewSignedData(content)
	if err != nil {
		return nil, fmt.Errorf("encodeCMS: No se pudo inicializar SignedData. %s", err)
	}

	if err := signedData.AddSigner(certificate, privateKey, pkcs7.SignerInfoConfig{}); err != nil {
		return nil, fmt.Errorf("encodeCMS: No se pudo agregar firmante: %s", err)
	}

	detachedSignature, err := signedData.Finish()
	if err != nil {
		return nil, fmt.Errorf("encodeCMS: No se pudo finalizar de firmar mensaje: %s", err)
	}

	return detachedSignature, nil
}

func readPrivateKey(file string) (any, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	p, _ := pem.Decode(buf)
	if p == nil {
		return nil, errors.New("no pem block found")
	}
	return x509.ParsePKCS8PrivateKey(p.Bytes)
}

func readCertificate(file string) (*x509.Certificate, error) {
	var certificate *x509.Certificate

	certPEMBlock, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var blocks [][]byte
	for {
		var certDERBlock *pem.Block
		certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
		if certDERBlock == nil {
			break
		}

		if certDERBlock.Type == "CERTIFICATE" {
			blocks = append(blocks, certDERBlock.Bytes)
		}
	}

	for _, block := range blocks {
		cert, err := x509.ParseCertificate(block)
		if err != nil {
			continue
		}
		certificate = cert
	}
	return certificate, nil
}

func PrintlnAsJSON(obj interface{}) {
	data, err := json.MarshalIndent(obj, "", "    ")
	if err == nil {
		fmt.Println(string(data))
	}
}

func PrintlnAsXML(obj interface{}) {
	data, err := xml.MarshalIndent(obj, " ", "  ")
	if err == nil {
		fmt.Println(string(data))
	}
}
