package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

func RSAReadKeyFromFile(filename string) []byte {
	f, err := os.Open(filename)
	var b []byte

	if err != nil {
		return b
	}

	defer f.Close()
	fileInfo, _ := f.Stat()
	b = make([]byte, fileInfo.Size())
	f.Read(b)
	return b
}

func RSAEncrypt(data, publicBytes []byte) ([]byte, error) {
	var res []byte
	block, _ := pem.Decode(publicBytes)

	if block == nil {
		return res, fmt.Errorf("Unable to encrypt, public key may be incorrect")
	}

	keyInit, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return res, fmt.Errorf("Unable to encrypt, public key may be incorrect, %v", err)
	}

	pubKey := keyInit.(*rsa.PublicKey)
	res, err = rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)

	if err != nil {
		return res, fmt.Errorf("Unable to encrypt, public key may be incorrect, %v", err)
	}

	return []byte(EncodeStr2Base64(string(res))), nil
}

func RSADecrypt(base64Data, privateBytes []byte) ([]byte, error) {
	var res []byte
	data := []byte(DecodeStrFromBase64(string(base64Data)))
	block, _ := pem.Decode(privateBytes)

	if block == nil {
		return res, fmt.Errorf("Unable to decrypt, private key may be incorrect")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return res, fmt.Errorf("Unable to decrypt, private key may be incorrect, %v", err)
	}
	res, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)

	if err != nil {
		return res, fmt.Errorf("Unable to decrypt, private key may be incorrect, %v", err)
	}

	return res, nil
}

func EncodeStr2Base64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func DecodeStrFromBase64(str string) string {
	decodeBytes, _ := base64.StdEncoding.DecodeString(str)
	return string(decodeBytes)
}
