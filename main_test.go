package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"testing"
)

const (
	encFile        = "testdata/Atzc5Drr8l7ngis8GmUTMI6vMQNjOU9zGQ2SYRkjwq44.enc"
	fileHashBase64 = "1CDSBESf8QCJRbX+EWVXRf3n8y7TZNu41K5BW4UbMTo="
	hexZMEDIAKEY   = "0A2069A349914734B9359DA0CD8923E6DFDE06F1E2BCE23222C738C521570BA8242A1220A1F5AEB2E620F73007FA853200559B2669455BB5818F619397C638042D8F7F2A18B984A5F1052000"
)

func TestDecryptFile(t *testing.T) {
	// https://mmg-fna.whatsapp.net/d/f/Atzc5Drr8l7ngis8GmUTMI6vMQNjOU9zGQ2SYRkjwq44.enc

	data, err := decryptMediaFile(encFile, hexZMEDIAKEY, mediaTypeVideo)
	if err != nil {
		t.Fatal(err)
	}

	fileHash, err := base64.StdEncoding.DecodeString(fileHashBase64)
	if err != nil {
		t.Fatal(err)
	}

	mac := sha256.New()
	mac.Write(data)
	testFileHash := mac.Sum(nil)

	if !bytes.Equal(testFileHash, fileHash) {
		t.Fatalf("got %v, want %v", testFileHash, fileHash)
	}
}
