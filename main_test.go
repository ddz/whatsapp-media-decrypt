package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"testing"
)

// [0] iOS encrypted mp4
// [1] Android encrypted mp4
var encFile = []string{"testdata/Atzc5Drr8l7ngis8GmUTMI6vMQNjOU9zGQ2SYRkjwq44.enc", "testdata/AnUpYQ390rgUBOQRhuwCyNqo_9KGATdmLUq-ghYEx-D9.enc"}

// From x-wa-metadata header filehash value
// or using bash, zsh, ...: cat <filename> | sha256sum | xxd -r -p | base64
var encFileHashBase64 = []string{"ofWusuYg9zAH+oUyAFWbJmlFW7WBj2GTl8Y4BC2Pfyo=", "NbmAmc7WfNwQFDYjE8iNzjZ0+RS8tR59VHNgVGG/FcM="}
var fileHashBase64 = []string{"1CDSBESf8QCJRbX+EWVXRf3n8y7TZNu41K5BW4UbMTo=", "2z310JnCt9q8ff4K6JIOj2UNrCUFvS1qFy/4JsGK+aE="}

// sqlite> select hex(ZMEDIAKEY) from ZWAMEDIAITEM; -- ChatStorage.sqlite
// sqlite> select hex(media_key) from message_media; -- msgstore.db
var mediaKey = []string{"0A2069A349914734B9359DA0CD8923E6DFDE06F1E2BCE23222C738C521570BA8242A1220A1F5AEB2E620F73007FA853200559B2669455BB5818F619397C638042D8F7F2A18B984A5F1052000",
	"14F9C1B3BB5E66D9A593999A5E0ED3D03ABFECA84320D17763C2B44205E91C17"}

func TestDecryptFile(t *testing.T) {
	// [0] https://mmg-fna.whatsapp.net/d/f/Atzc5Drr8l7ngis8GmUTMI6vMQNjOU9zGQ2SYRkjwq44.enc (iOS)
	// [1] https://mmg-fna.whatsapp.net/d/f/AnUpYQ390rgUBOQRhuwCyNqo_9KGATdmLUq-ghYEx-D9.enc (Android)

	for i := 0; i < len(encFile); i++ {
		encFileHash, err := base64.StdEncoding.DecodeString(encFileHashBase64[i])
		if err != nil {
			t.Fatal(err)
		}

		encFileData, err := ioutil.ReadFile(encFile[i])
		if err != nil {
			t.Fatal(err)
		}

		mac := sha256.New()
		mac.Write(encFileData)
		testEncFileHash := mac.Sum(nil)
		if !bytes.Equal(testEncFileHash, encFileHash) {
			t.Fatalf("got %v, want %v", testEncFileHash, encFileHash)
		}

		data, err := decryptMediaFile(encFile[i], mediaKey[i], mediaTypeVideo)
		if err != nil {
			t.Fatal(err)
		}

		fileHash, err := base64.StdEncoding.DecodeString(fileHashBase64[i])
		if err != nil {
			t.Fatal(err)
		}

		mac = sha256.New()
		mac.Write(data)
		testFileHash := mac.Sum(nil)

		if !bytes.Equal(testFileHash, fileHash) {
			t.Fatalf("got %v, want %v", testFileHash, fileHash)
		}
	}
}
