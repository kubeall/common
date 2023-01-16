package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type RsaSecurity struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewRsaSecurityFromRsaKey(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) (result *RsaSecurity) {
	result.publicKey = publicKey
	result.privateKey = privateKey
	return
}
func NewRsaSecurityFromStringKey(publicKey, privateKey string) (result *RsaSecurity, err error) {
	if len(publicKey) > 0 {
		block, _ := pem.Decode([]byte(publicKey))
		if block == nil {
			return nil, errors.New("get public key error")
		}
		// x509 parse public key
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		result.publicKey = pub.(*rsa.PublicKey)
	}
	if len(privateKey) > 0 {
		block, _ := pem.Decode([]byte(privateKey))
		if block == nil {
			return nil, errors.New("get private key error")
		}
		result.privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			pri2, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			result.privateKey = pri2.(*rsa.PrivateKey)
		}
	}
	return
}

func GenerateRASPrivateAndPublicKeys() (privateKey, publicKey []byte, err error) {
	pri, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	derTx := x509.MarshalPKCS1PrivateKey(pri)
	block := pem.Block{Type: "RSA PRIVATE KEY", Bytes: derTx}
	privateKey = pem.EncodeToMemory(&block)
	stream, err := x509.MarshalPKIXPublicKey(&pri.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	block = pem.Block{Type: "RSA PUBLIC KEY", Bytes: stream}
	publicKey = pem.EncodeToMemory(&block)
	return privateKey, publicKey, nil
}

// PublicKeyEncrypt
func (s *RsaSecurity) PublicKeyEncrypt(input []byte) ([]byte, error) {
	if s.publicKey == nil {
		return []byte(""), errors.New(`please set the public key in advance`)
	}
	return pubKeyByte(s.publicKey, input, true)
}

// PublicKeyDecrypt
func (s *RsaSecurity) PublicKeyDecrypt(input []byte) ([]byte, error) {
	if s.publicKey == nil {
		return []byte(""), errors.New(`please set the public key in advance`)
	}
	return pubKeyByte(s.publicKey, input, false)
}

// PrivateKeyEncrypt
func (s *RsaSecurity) PrivateKeyEncrypt(input []byte) ([]byte, error) {
	if s.privateKey == nil {
		return []byte(""), errors.New(`please set the private key in advance`)
	}
	return priKeyByte(s.privateKey, input, true)
}

// PrivateKeyDecrypt
func (s *RsaSecurity) PrivateKeyDecrypt(input []byte) ([]byte, error) {
	if s.privateKey == nil {
		return []byte(""), errors.New(`please set the private key in advance`)
	}
	return priKeyByte(s.privateKey, input, false)
}
