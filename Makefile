BINARY_NAME=app
DOCKER_IMAGE=monitoring
DOCKERFILE=build/Dockerfile

build-monitoring:
	cd src && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../build/$(BINARY_NAME) *.go

docker-build:
	docker build -f $(DOCKERFILE) -t $(DOCKER_IMAGE) .

run:
	docker run -p 8080:8080 -p 9090:9090 $(DOCKER_IMAGE)

run-test:
	cd src && go run *.go

clean:
	rm -f build/$(BINARY_NAME)

build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/app-linux src/*.go

build-windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o build/app-windows.exe src/*.go

build-macos:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o build/app-macos src/*.go

build-all: 
	build-linux build-windows build-macos

auto:
	@echo "🧹 Nettoyage..."
	@make clean
	@echo "🔧 Compilation locale..."
	@make build-monitoring
	@echo "🐳 Construction Docker..."
	@make docker-build
	@echo "🚀 Lancement du conteneur..."
	@make run

git-push:
	@if [ -z "$(m)" ]; then \
		echo "❌ Merci de fournir un message de commit avec m=\"ton message\""; \
		exit 1; \
	fi
	@git add .
	@git commit -m "$(m)"
	@git push
