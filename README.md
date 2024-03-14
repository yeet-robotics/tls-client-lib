# tls-client-lib

## Compile

```shell
go build -buildmode=c-shared -o tlslib.dll .
```

## Crosscompile for Windows

```shell
sudo apt install mingw-w64
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -buildmode=c-shared -o tlslib.dll .
```

## Fingerprints

- `Firefox89` (desktop)
- `Chrome93` (desktop and mobile)
- `Safari604` (desktop and mobile)
- `Safari605` (desktop and mobile)