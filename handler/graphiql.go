package handler

import "net/http"

// GraphiQL is an in-browser IDE for exploring GraphiQL APIs.
// This handler returns GraphiQL when requested.
//
// For more information, see https://github.com/graphql/graphiql.
type GraphiQL struct{}

func (h GraphiQL) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		respond(w, errorJSON("only GET requests are supported"), http.StatusMethodNotAllowed)
		return
	}

	w.Write(graphiql)
}

var graphiql = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<style>body {height: 100vh; margin: 0; width: 100%; overflow: hidden;}</style>
		<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.css" />
		<script src="//cdnjs.cloudflare.com/ajax/libs/fetch/2.0.3/fetch.min.js"></script>
		<script src="//cdnjs.cloudflare.com/ajax/libs/react/16.2.0/umd/react.production.min.js"></script>
		<script src="//cdnjs.cloudflare.com/ajax/libs/react-dom/16.2.0/umd/react-dom.production.min.js"></script>
		<script src="//cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.js"></script>
	</head>
	<body>
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>		
			function fetchGQL(params) {		
				return fetch("/graphql", {		
					method: "post",		
					body: JSON.stringify(params),		
					credentials: "include",		
				}).then(function (resp) {		
					return resp.text();		
				}).then(function (body) {		
					try {		
						return JSON.parse(body);		
					} catch (error) {		
						return body;		
					}		
				});		
			}		
		
			ReactDOM.render(		
				React.createElement(GraphiQL, {fetcher: fetchGQL}),		
				document.getElementById("graphiql")		
			)		
		</script>
	</body>
</html>
`)
