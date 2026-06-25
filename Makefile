build:
	go build  -trimpath -o main ./main.go

build-arm64:
	GOARCH=arm64 GOOS=darwin go build -trimpath -o zkproof_darwin_arm64 ./main.go

build-amd64:
	GOARCH=amd64 GOOS=darwin go build -trimpath -o zkproof_darwin_amd64 ./main.go

build-linux:
	GOARCH=amd64 GOOS=linux go build -trimpath -o zkproof_linux_amd64 ./main.go

build-windows:
	GOARCH=amd64 GOOS=windows go build -trimpath -o zkproof_windows_amd64.exe ./main.go