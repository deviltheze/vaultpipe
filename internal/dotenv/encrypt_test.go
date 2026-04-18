package dotenv

import (
	"bytes"
	"testing"
)

func testKey() []byte {
	return bytes.Repeat([]byte("k"), 32)
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	plain := []byte("super-secret-value")
	enc, err := Encrypt(plain, testKey())
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	dec, err := Decrypt(enc, testKey())
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if !bytes.Equal(dec, plain) {
		t.Errorf("got %q, want %q", dec, plain)
	}
}

func TestEncrypt_DifferentCiphertexts(t *testing.T) {
	plain := []byte("value")
	a, _ := Encrypt(plain, testKey())
	b, _ := Encrypt(plain, testKey())
	if a == b {
		t.Error("expected different ciphertexts due to random nonce")
	}
}

func TestEncrypt_BadKeyLength(t *testing.T) {
	_, err := Encrypt([]byte("x"), []byte("short"))
	if err == nil {
		t.Error("expected error for short key")
	}
}

func TestDecrypt_BadKeyLength(t *testing.T) {
	_, err := Decrypt("aGVsbG8=", []byte("short"))
	if err == nil {
		t.Error("expected error for short key")
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	_, err := Decrypt("!!!notbase64!!!", testKey())
	if err == nil {
		t.Error("expected error for invalid base64")
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	enc, _ := Encrypt([]byte("value"), testKey())
	// flip last char
	runes := []rune(enc)
	if runes[len(runes)-1] == 'A' {
		runes[len(runes)-1] = 'B'
	} else {
		runes[len(runes)-1] = 'A'
	}
	_, err := Decrypt(string(runes), testKey())
	if err == nil {
		t.Error("expected error for tampered ciphertext")
	}
}
