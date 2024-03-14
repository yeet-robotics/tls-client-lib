module tls-client-lib

go 1.18

require tls-client-go v0.0.0-20210428134155-4694784420cc

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/elliotchance/pie/v2 v2.0.1 // indirect
	github.com/refraction-networking/utls v1.1.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	golang.org/x/crypto v0.0.0-20220518034528-6f7dac969898 // indirect
	golang.org/x/exp v0.0.0-20220321173239-a90fa8a75705 // indirect
	golang.org/x/net v0.0.0-20220520000938-2e3eb7b945c2 // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/text v0.4.0 // indirect
)

replace golang.org/x/net => ./tls-client-go/net

replace github.com/refraction-networking/utls => ./tls-client-go/utls

replace tls-client-go => ./tls-client-go
