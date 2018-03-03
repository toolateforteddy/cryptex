package process

import (
	"testing"

	"github.com/toolateforteddy/cryptex/cryptopasta"
)

func TestEncryptDecryptPassword(t *testing.T) {
	password := "I am a password"
	cryptKey := cryptopasta.NewEncryptionKey()

	cb, err := cryptopasta.Encrypt([]byte(password), cryptKey)
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}
	if string(cb) == password {
		t.Fatalf("encrypt not obfuscate Got %q", string(cb))
	}
	db, err := cryptopasta.Decrypt(cb, cryptKey)
	if err != nil {
		t.Fatalf("failed to decrypt: %v", err)
	}
	if string(db) != password {
		t.Fatalf("roundtrip was not clean. Got %q", string(db))
	}
}

func TestEncryptDecryptLine(t *testing.T) {
	password := `foobar = "Iam a p assword"`
	pKey := Key{cryptopasta.NewEncryptionKey()}
	cb := pKey.encryptLine(password)

	if string(cb) == password {
		t.Fatalf("encrypt not obfuscate Got %q", string(cb))
	}
	db := pKey.decryptLine(cb)

	if string(db) != password {
		t.Fatalf("roundtrip was not clean. Got %q", string(db))
	}
}
