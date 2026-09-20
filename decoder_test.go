package totalcommanderftppasswordrecovery

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestDecryptPassword(t *testing.T) {
	expected, err := hex.DecodeString("fdf3b350e9b8c5fbe82d478d")
	if err != nil {
		t.Fatal(err)
	}
	decoder := New()
	got, err := decoder.DecryptPassword("00112233445566778899aabbccddeeff")
	if err != nil {
		t.Fatalf("DecryptPassword: %v", err)
	}
	if !bytes.Equal(got, expected) {
		t.Fatalf("golden mismatch: got %x want %x", got, expected)
	}
}

func TestDecryptPasswordIgnoresWhitespaceAndCase(t *testing.T) {
	expected, _ := hex.DecodeString("fdf3b350e9b8c5fbe82d478d")
	got, err := New().DecryptPassword("00 11 22 33 44 55 66 77 88 99 AA BB CC DD EE FF")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, expected) {
		t.Fatalf("got %x want %x", got, expected)
	}
}

func TestDecryptPasswordRejectsTooShort(t *testing.T) {
	if _, err := New().DecryptPassword("00112233"); err == nil {
		t.Fatal("expected error for too-short ciphertext")
	}
}

func TestDecryptPasswordRejectsInvalidHex(t *testing.T) {
	if _, err := New().DecryptPassword("00112233445566778899aabbccddeefg"); err == nil {
		t.Fatal("expected error for invalid hex")
	}
}

func TestHexStringToByteArray(t *testing.T) {
	if got := HexStringToByteArray("4142"); !bytes.Equal(got, []byte{0x41, 0x42}) {
		t.Fatalf("got %x", got)
	}
	if HexStringToByteArray("414") != nil {
		t.Fatal("odd length should fail")
	}
	if HexStringToByteArray("41ag") != nil {
		t.Fatal("invalid nibble should fail")
	}
}

func TestRol8(t *testing.T) {
	if Rol8(0xAB, 0) != 0xAB {
		t.Fatal("rol 0")
	}
	if Rol8(0xAB, 1) != 0x57 {
		t.Fatal("rol 1")
	}
	if Rol8(0xAB, 8) != 0xAB {
		t.Fatal("rol 8")
	}
}
