run: build
	./bin/ngrok-file-server
build: clean
	go build -o bin/ngrok-file-server -v api/main.go
clean:
	rm -rf ./bin
test:
	go test -v ./...