M = $(shell printf "\033[34;1m▶\033[0m")

build: dep ; $(info $(M) Building project...)
	go build

clean: ; $(info $(M) [TODO] Removing generated files... )
	$(RM) schema/bindata.go

schema: dep ; $(info $(M) Embedding schema files into binary...)
	go generate ./schema

server: schema ; $(info $(M) Starting development server...)
	go run server.go

image: ; $(info $(M) Building application image...)
	docker build -t graphql-go-example .

container: image ; $(info $(M) Running application container...)
	docker run -p 8000:8000 graphql-go-example:latest

.PHONY: build clean container dep image schema setup server
