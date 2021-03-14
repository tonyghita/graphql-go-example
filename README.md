# graphql-go-example

This codebase is an example of how [`github.com/graph-gophers/graphql-go`][0]
might be used in production.

This project launches a server that exposes the StarWars API hosted at
<https://swapi.dev> as a GraphQL API implemented in Go.

## Usage

```sh
make server
```

This launches the HTTP server on port `8000` of your machine.

Visiting http://localhost:8000 will return a GraphiQL client that you can use to make
requests against the API.

[0]: https://github.com/graph-gophers/graphql-go
