FROM golang:buster AS build
WORKDIR /go/src/github.com/tonyghita/graphql-go-example
COPY . .
RUN go version && go env && go build -o /example

FROM gcr.io/distroless/base
COPY --from=build /example /
ENTRYPOINT ["/example"]
EXPOSE 8000
