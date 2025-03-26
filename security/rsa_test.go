package security

import (
	"strings"
	"testing"
)

func TestRSA(t *testing.T) {
	private, public, _ := GenerateRASPrivateAndPublicKeys()
	Rsa, err := NewRsaSecurityFromStringKey(string(public), string(private))
	if err != nil {
		t.Fatal(err)
	}
	data := "this is rsa test raw data"
	encryptData, err := Rsa.Encrypt([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	decryptData, err := Rsa.Decrypt(encryptData)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("raw data : %s", data)
	t.Logf("decryptData: %s", decryptData)
	if strings.EqualFold(data, string(decryptData)) {
		t.Log("rsa ok !")
	} else {
		t.Fatal("rsa failed !")
	}
}
