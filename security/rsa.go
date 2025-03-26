package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"github.com/golang-jwt/jwt/v5"
)

type RsaSecurity struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewRsaSecurityFromRsaKey(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) (result *RsaSecurity) {
	result = &RsaSecurity{}
	result.publicKey = publicKey
	result.privateKey = privateKey
	return
}
func NewRsaSecurityFromStringKey(publicPem, privatePem string) (result *RsaSecurity, err error) {
	result = &RsaSecurity{}
	result.publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(publicPem))
	if err != nil {
		return nil, err
	}
	result.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(privatePem))
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

func (s *RsaSecurity) Encrypt(input []byte) (encryptedBlockBytes []byte, err error) {
	msgLen := len(input)
	h := sha256.New()
	rng := rand.Reader
	label := []byte("efucloud-encrypt")
	step := s.publicKey.Size() - 2*h.Size() - 2
	var encryptedBytes []byte
	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}
		encryptedBlockBytes, err := rsa.EncryptOAEP(h, rng, s.publicKey, input[start:finish], label)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}
	return encryptedBytes, nil
}

func (s *RsaSecurity) Decrypt(input []byte) (decryptedBytes []byte, err error) {
	msgLen := len(input)
	step := s.privateKey.PublicKey.Size()
	h := sha256.New()
	rng := rand.Reader
	label := []byte("efucloud-encrypt")
	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}
		decryptedBlockBytes, err := rsa.DecryptOAEP(h, rng, s.privateKey, input[start:finish], label)
		if err != nil {
			return nil, err
		}
		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}
func DecryptData(private, input []byte) (decryptedBytes []byte, err error) {
	var privateKey *rsa.PrivateKey
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(private)
	if err != nil {
		return nil, err
	}
	msgLen := len(input)
	step := privateKey.PublicKey.Size()
	h := sha256.New()
	rng := rand.Reader
	label := []byte("efucloud-encrypt")
	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}
		decryptedBlockBytes, err := rsa.DecryptOAEP(h, rng, privateKey, input[start:finish], label)
		if err != nil {
			return nil, err
		}
		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}

func DecryptDataByPrivateKey(privateKey *rsa.PrivateKey, input []byte) (decryptedBytes []byte, err error) {
	msgLen := len(input)
	step := privateKey.PublicKey.Size()
	h := sha256.New()
	rng := rand.Reader
	label := []byte("efucloud-encrypt")
	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}
		decryptedBlockBytes, err := rsa.DecryptOAEP(h, rng, privateKey, input[start:finish], label)
		if err != nil {
			return nil, err
		}
		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}
func EncryptData(public, input []byte) (encryptedBytes []byte, err error) {
	var publicKey *rsa.PublicKey
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(public)
	if err != nil {
		return nil, err
	}
	msgLen := len(input)
	h := sha256.New()
	rng := rand.Reader
	label := []byte("efucloud-encrypt")
	step := publicKey.Size() - 2*h.Size() - 2
	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}
		encryptedBlockBytes, err := rsa.EncryptOAEP(h, rng, publicKey, input[start:finish], label)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

func EncryptDataByPublicKey(publicKey *rsa.PublicKey, input []byte) (encryptedBytes []byte, err error) {
	msgLen := len(input)
	h := sha256.New()
	rng := rand.Reader
	label := []byte("efucloud-encrypt")
	step := publicKey.Size() - 2*h.Size() - 2
	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}
		encryptedBlockBytes, err := rsa.EncryptOAEP(h, rng, publicKey, input[start:finish], label)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}
