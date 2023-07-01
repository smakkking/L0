.PHONY: build
build:
	go build -o main app.go models.go
.PHONY: start_docker
start_docker:
	docker-compose up -d
.PHONY: clear_docker
clear_docker:
	docker-compose down --volumes