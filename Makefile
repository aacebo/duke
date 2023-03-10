clean:
	rm -rf ./bin

build:
	go build -o bin/main app/index.go

run:
	go run example/main.go

compile:
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 app/index.go
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 app/index.go
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 app/index.go
