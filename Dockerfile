FROM golang:alpine AS build
WORKDIR /go/src/github.com/tonyghita/graphql-go-example
COPY . .
RUN go version && go env && go build -o /go/bin/example

FROM alpine
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /go/bin/example /app/
ENTRYPOINT ["./example"]
EXPOSE 8000
