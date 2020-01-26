# WhatsApp Media Decrypt

## Usage

## Example

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
$ go run main.go -o Atzc5Drr8l7ngis8GmUTMI6vMQNjOU9zGQ2SYRkjwq44.mp4 -t 2 ./Atzc5Drr8l7ngis8GmUTMI6vMQNjOU9zGQ2SYRkjwq44.enc 0A2069A349914734B9359DA0CD8923E6DFDE06F1E2BCE23222C738C521570BA8242A1220A1F5AEB2E620F73007FA853200559B2669455BB5818F619397C638042D8F7F2A18B984A5F1052000
```
