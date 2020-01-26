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
Usage of whatsapp-media-decrypt:
  -o file
    	write decrypted output to file
  -t type
    	media type (1 = image, 2 = video, 3 = audio, 4 = doc)
```

## Example

```
$ whatsapp-media-decrypt -o Atzc5Drr8l7ngis8GmUTMI6vMQNjOU9zGQ2SYRkjwq44.mp4 -t 2 Atzc5Drr8l7ngis8GmUTMI6vMQNjOU9zGQ2SYRkjwq44.enc 0A2069A349914734B9359DA0CD8923E6DFDE06F1E2BCE23222C738C521570BA8242A1220A1F5AEB2E620F73007FA853200559B2669455BB5818F619397C638042D8F7F2A18B984A5F1052000
```

### Acquire ChatStorage.sqlite

(Left as an exercise for the reader)

### Extract mediaKey

```
$ sqlite ChatStorage.sqlite
SQLite version 3.27.2 2019-03-09 15:45:46
Enter ".help" for usage hints.
sqlite> select ZMEDIAURL,ZVCARDSTRING,hex(ZMEDIAKEY) from ZWAMEDIAITEM where Z_PK = 1795;
https://mmg-fna.whatsapp.net/d/f/Atzc5Drr8l7ngis8GmUTMI6vMQNjOU9zGQ2SYRkjwq44.enc|video/mp4|0A2069A349914734B9359DA0CD8923E6DFDE06F1E2BCE23222C738C521570BA8242A1220A1F5AEB2E620F73007FA853200559B2669455BB5818F619397C638042D8F7F2A18B984A5F1052000
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

## References

[WhatsApp Web Reverse Engineered](https://github.com/sigalor/whatsapp-web-reveng)

[go-whatsapp](https://github.com/Rhymen/go-whatsapp)
