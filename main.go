package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Rhymen/go-whatsapp/crypto/cbc"
	"github.com/Rhymen/go-whatsapp/crypto/hkdf"
	"github.com/golang/protobuf/proto"
)

type mediaType int

const (
	_              = iota
	mediaTypeImage = mediaType(iota)
	mediaTypeVideo
	mediaTypeAudio
	mediaTypeDocument
)

var (
	appInfo = map[mediaType]string{
		mediaTypeImage:    "WhatsApp Image Keys",
		mediaTypeVideo:    "WhatsApp Video Keys",
		mediaTypeAudio:    "WhatsApp Audio Keys",
		mediaTypeDocument: "WhatsApp Document Keys",
	}
)

func main() {
	var Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage: %s -o FILE -t TYPE ENCFILE HEXMEDIAKEY\n\nOptions:\n",
			os.Args[0])
		flag.PrintDefaults()
	}

	mt := flag.Int("t", 0, "media `TYPE` (1 = image, 2 = video, 3 = audio, 4 = doc)")
	outputFileName := flag.String("o", "", "write decrypted output to `FILE`")
	flag.Parse()

	if *mt == 0 || len(*outputFileName) == 0 || flag.NArg() < 2 {
		Usage()
		os.Exit(1)
	}

	data, err := decryptMediaFile(flag.Args()[0], flag.Args()[1], mediaType(*mt))
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(*outputFileName, data, 0400)
	if err != nil {
		log.Fatal(err)
	}
}

func decryptMediaFile(encFilePath string, hexMediaKey string, mt mediaType) (
	[]byte,
	error,
) {
	encFileData, err := ioutil.ReadFile(encFilePath)
	if err != nil {
		return nil, err
	}

	mediaKeyBlob, err := hex.DecodeString(hexMediaKey)
	if err != nil {
		return nil, err
	}

	var mediaKey []byte
	var fileHash []byte

	if len(mediaKeyBlob) == 32 {
		mediaKey = mediaKeyBlob
	} else {
		// Decode as protobuf
		mk := &MediaKey{}
		err = proto.Unmarshal(mediaKeyBlob, mk)
		if err != nil {
			return nil, err
		}

		mediaKey = []byte(*mk.MediaKey)
		fileHash = []byte(*mk.FileEncSha256)
	}

	if len(fileHash) > 0 {
		// Verify the fileHash of the .enc file data
		h := sha256.New()
		h.Write(encFileData)
		encFileHash := h.Sum(nil)
		if !bytes.Equal(encFileHash, fileHash) {
			return nil, fmt.Errorf(".enc file hash does not match mediaKey")
		}
	}

	data, err := decryptMedia(encFileData, mediaKey, mt)
	if err != nil {
		return nil, err
	}

	return data, err
}

func decryptMedia(encFileData []byte, mediaKey []byte, mt mediaType) (
	[]byte,
	error,
) {
	//
	// Implement reverse engineered media decryption algorithm from:
	// https://github.com/sigalor/whatsapp-web-reveng#decryption
	//

	// mediaKey should be 32 bytes
	if len(mediaKey) != 32 {
		return nil, fmt.Errorf("mediaKey length %d != 32",
			len(mediaKey))
	}

	mediaKeyExpanded, err := hkdf.Expand(mediaKey, 112, appInfo[mt])
	if err != nil {
		return nil, err
	}

	iv := mediaKeyExpanded[0:16]
	cipherKey := mediaKeyExpanded[16:48]
	macKey := mediaKeyExpanded[48:80]
	//refKey := mediaKeyExpanded[80:]

	fileLen := len(encFileData) - 10
	file := encFileData[:fileLen]
	mac := encFileData[fileLen:]

	err = validateMedia(iv, file, macKey, mac)
	if err != nil {
		return nil, err
	}

	data, err := cbc.Decrypt(cipherKey, iv, file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// copied from https://github.com/Rhymen/go-whatsapp/blob/master/media.go
func validateMedia(iv []byte, file []byte, macKey []byte, mac []byte) error {
	h := hmac.New(sha256.New, macKey)
	n, err := h.Write(append(iv, file...))
	if err != nil {
		return err
	}
	if n < 10 {
		return fmt.Errorf("hash to short")
	}
	if !hmac.Equal(h.Sum(nil)[:10], mac) {
		return fmt.Errorf("invalid media hmac")
	}
	return nil
}
