M = $(shell printf "\033[34;1m▶\033[0m")

server: ; $(info $(M) Starting development server...)
	@ go run server.go
.PHONY: server

image: ; $(info $(M) Building application image...)
	@ docker build -t graphql-go-example .
.PHONY: image

container: image ; $(info $(M) Running application container...)
	@ docker run -p 8000:8000 graphql-go-example:latest
.PHONY: container
