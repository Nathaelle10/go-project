BINARY_NAME=app
DOCKER_IMAGE=myapp
DOCKERFILE=build/Dockerfile

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o src/$(BINARY_NAME) src/main.go

docker-build:
	docker build -f $(DOCKERFILE) -t $(DOCKER_IMAGE) .

run:
	docker run -p 8080:8080 -p 9090:9090 $(DOCKER_IMAGE)

clean:
	rm -f src/$(BINARY_NAME)

build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/app-linux src/main.go

build-windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o build/app-windows.exe src/main.go

build-macos:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o build/app-macos src/main.go

build-all: build-linux build-windows build-macos

