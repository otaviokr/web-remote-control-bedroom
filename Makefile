build:
	mkdir -p functions
	GOOS=linux
	GOARCH=amd64
	GO111MODULE=on
	go build -o functions/blinkt ./src/blinkt.go
