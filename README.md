# WhatsApp Media Decrypt

A recent [high-profile forensic investigation](https://www.vice.com/en_us/article/v74v34/saudi-arabia-hacked-jeff-bezos-phone-technical-report)
reported that “due to end-to-end encryption employed by WhatsApp, it
is virtually impossible to decrypt the contents of the downloader
[.enc file] to determine if it contained any malicious code in
addition to the delivered video.”

This project demonstrates how to decrypt encrypted media files
downloaded from WhatsApp.

## Installation

```
$ go get github.com/ddz/whatsapp-media-decrypt
```

## Usage

```
Usage: ./whatsapp-media-decrypt -o FILE -t TYPE ENCFILE HEXMEDIAKEY

Options:
  -o FILE
    	write decrypted output to FILE
  -t TYPE
    	media TYPE (1 = image, 2 = video, 3 = audio, 4 = doc)
```

## Example

### Extract media key from iOS ChatStorage.sqlite

The media key is stored within a protobuf message that is stored
hex-encoded in the `ZMEDIAKEY` column.

```
$ sqlite ChatStorage.sqlite
SQLite version 3.27.2 2019-03-09 15:45:46
Enter ".help" for usage hints.
sqlite> select ZMEDIAURL,ZVCARDSTRING,hex(ZMEDIAKEY) from ZWAMEDIAITEM where Z_PK = 1795;
https://mmg-fna.whatsapp.net/d/f/Atzc5Drr8l7ngis8GmUTMI6vMQNjOU9zGQ2SYRkjwq44.enc|video/mp4|0A2069A349914734B9359DA0CD8923E6DFDE06F1E2BCE23222C738C521570BA8242A1220A1F5AEB2E620F73007FA853200559B2669455BB5818F619397C638042D8F7F2A18B984A5F1052000
sqlite> .quit
```

### Extract media key from Android msgstore.db

The media key is stored hex-encoded in the `media_key` column.

```
$ sqlite msgstore.db
SQLite version 3.27.2 2019-03-09 15:45:46
Enter ".help" for usage hints.
sqlite> select message_url,mime_type,hex(media_key) from message_media where message_row_id = 1337;
https://mmg-fna.whatsapp.net/d/f/AnUpYQ390rgUBOQRhuwCyNqo_9KGATdmLUq-ghYEx-D9.enc|video/mp4|14F9C1B3BB5E66D9A593999A5E0ED3D03ABFECA84320D17763C2B44205E91C17
sqlite> .quit
```

### Download Encrypted Media File

```
$ curl -O https://mmg-fna.whatsapp.net/d/f/Atzc5Drr8l7ngis8GmUTMI6vMQNjOU9zGQ2SYRkjwq44.enc
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  389k  100  389k    0     0  1966k      0 --:--:-- --:--:-- --:--:-- 1956k
```

### Decrypt Media File

```
$ whatapp-media-decrypt -o Atzc5Drr8l7ngis8GmUTMI6vMQNjOU9zGQ2SYRkjwq44.mp4 -t 2 ./Atzc5Drr8l7ngis8GmUTMI6vMQNjOU9zGQ2SYRkjwq44.enc 0A2069A349914734B9359DA0CD8923E6DFDE06F1E2BCE23222C738C521570BA8242A1220A1F5AEB2E620F73007FA853200559B2669455BB5818F619397C638042D8F7F2A18B984A5F1052000
```

## FAQ

### Does this break WhatsApp encryption?

No. WhatsApp's encryption is end-to-end, which ensures that only the
sender and recipient can read the message and especially not any
servers (or attackers!) in-between them. This uses a cryptographic key
stored on one of the endpoints to decrypt a media attachment in the
same way that the WhatsApp app does to display it on the screen.

### Does this mean my WhatsApp media files are not encrypted at rest?

No. WhatsApp uses [iOS Data
Protection](https://support.apple.com/guide/security/how-data-files-are-created-and-protected-sece8608431d/1/web/1)
to encrypt user data files (including `ChatStorage.sqlite`) using the
device-specific and unrecoverable hardware UID key as well as a key
derived from the user's passcode. It may not be decrypted without
physical access to the specific iOS device that created the file as
well as knowledge of the user's passcode.

### Can you help me decrypt someone's WhatsApp?

No.

## References
Engelke, Lucas. [go-whatsapp](https://github.com/Rhymen/go-whatsapp)

Graham, Robert. [How to decrypt WhatsApp end-to-end media files](https://blog.erratasec.com/2020/01/how-to-decrypt-whatsapp-end-to-end.html)

Marczak, Bill. "[Some Directions for Further Investigation in the Bezos Hack Case](https://medium.com/@billmarczak/bezos-hack-mbs-mohammed-bin-salman-whatsapp-218e1b4e1242)"

Sigalor. [WhatsApp Web Reverse Engineered](https://github.com/sigalor/whatsapp-web-reveng)

WhatsApp. [WhatsApp Encryption Overview](https://www.whatsapp.com/security/WhatsApp-Security-Whitepaper.pdf)
