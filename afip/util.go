package afip

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"go.mozilla.org/pkcs7"
	"golang.org/x/crypto/pkcs12"
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

	/* solo para debugear
	pemCMS := pem.EncodeToMemory(&pem.Block{Type: "PKCS7", Bytes: detachedSignature})
	log.Printf("\n%s\n", string(pemCMS[:]))
	*/

	return detachedSignature, nil
}

// decodePkcs12 extrae del archivo .p12 el certificado y la clave privada RSA
func decodePkcs12(p12 string, password string) (*x509.Certificate, *rsa.PrivateKey, error) {

	pkcs, err := ioutil.ReadFile(p12)
	if err != nil {
		return nil, nil, fmt.Errorf("decodePkcs12: Error leyendo archivo %s. %s", p12, err.Error())
	}

	privateKey, certificate, err := pkcs12.Decode(pkcs, password)
	if err != nil {
		return nil, nil, fmt.Errorf("decodePkcs12: Error decodificando PKCS#12. %s", err.Error())
	}

	rsaPrivateKey, isRsaKey := privateKey.(*rsa.PrivateKey)
	if !isRsaKey {
		return nil, nil, fmt.Errorf("decodePkcs12: El certificado PKCS#12 debe contener una clave privada RSA")
	}

	/* solo para debugear
	keyBytes := x509.MarshalPKCS1PrivateKey(rsaPrivateKey)
	pemPrivateKey := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes})
	pemCertificate := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificate.Raw})
	log.Printf("\n%s\n", string(pemCertificate[:]))
	log.Printf("\n%s\n", string(pemPrivateKey[:]))
	*/

	return certificate, rsaPrivateKey, nil
}
