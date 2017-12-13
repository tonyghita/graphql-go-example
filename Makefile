M = $(shell printf "\033[34;1m▶\033[0m")

build: ; $(info $(M) Building project...)
	go build

clean: ; $(info $(M) [TODO] Removing generated files... )
	$(RM) schema/_bindata.go

dep: ; $(info $(M) Ensuring vendored dependencies are up-to-date...)
	dep ensure

schema: ; $(info $(M) Embedding schema files into binary...)
	go generate ./schema

setup: ; $(info $(M) Fetching github.com/golang/dep...)
	go get -u github.com/golang/dep/cmd/dep

server: ; $(info $(M) Starting development server...)
	go run server.go

.PHONY: build clean dep schema setup server
