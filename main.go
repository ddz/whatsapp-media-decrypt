package main

import (
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
	mt := flag.Int("t", 0, "media `type` (1 = image, 2 = video, 3 = audio, 4 = doc)")
	outputFileName := flag.String("o", "", "write decrypted output to `file`")
	flag.Parse()

	if *mt == 0 || len(*outputFileName) == 0 || flag.NArg() < 2 {
		flag.Usage()
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
	encFile, err := ioutil.ReadFile(encFilePath)
	if err != nil {
		return nil, err
	}

	mediaKey, err := hex.DecodeString(hexMediaKey)
	if err != nil {
		return nil, err
	}

	if len(mediaKey) != 32 {
		// Assume it's hex(ZMEDIAKEY) from ZWAMEDIAITEM table
		if mediaKey[0] == 0x0A && mediaKey[1] == 0x20 {
			// XXX: Not sure what the encoding and rest are,
			// but this is where the raw key is
			mediaKey = mediaKey[2 : 2+32]
		} else {
			return nil, fmt.Errorf("unknown mediaKey format")
		}
	}

	data, err := decryptMedia(encFile, mediaKey, mt)
	if err != nil {
		return nil, err
	}

	return data, err
}

func decryptMedia(encFile []byte, mediaKey []byte, mt mediaType) (
	[]byte,
	error,
) {
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

	fileLen := len(encFile) - 10
	file := encFile[:fileLen]
	mac := encFile[fileLen:]

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
