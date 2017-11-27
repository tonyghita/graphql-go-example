package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

func parse(r *http.Request) (Request, error) {
	// We always need to read and close the request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return Request{}, errors.New("unable to read request body")
	}
	_ = r.Body.Close()

	var request Request

	switch r.Method {
	case "POST":
		request = parsePost(body)
	case "GET":
		request = parseGet(r.URL.Query())
	default:
		err = errors.New("only POST and GET requests are supported")
	}

	return request, err
}

func parseGet(v url.Values) Request {
	var (
		queries   = v["query"]
		names     = v["operationName"]
		variables = v["variables"]
		qLen      = len(queries)
		nLen      = len(names)
		vLen      = len(variables)
	)

	if qLen == 0 {
		return Request{}
	}

	var requests = make([]Query, 0, qLen)
	var isBatch bool

	// This loop assumes there will be a corresponding element at each index
	// for query, operation name, and variable fields.
	//
	// NOTE: This could be a bad assumption. Maybe we want to do some validation?
	for i, q := range queries {
		var n string
		if i < nLen {
			n = names[i]
		}

		var m = map[string]interface{}{}
		if i < vLen {
			str := variables[i]
			if err := json.Unmarshal([]byte(str), &m); err != nil {
				m = nil // TODO: Improve error handling here.
			}
		}

		requests = append(requests, Query{Query: q, OpName: n, Variables: m})
	}

	if qLen > 1 {
		isBatch = true
	}

	return Request{queries: requests, isBatch: isBatch}
}

func parsePost(b []byte) Request {
	if len(b) == 0 {
		return Request{}
	}

	var queries []Query
	var isBatch bool

	// Inspect the first character to inform how the body is parsed.
	switch b[0] {
	case '{':
		q := Query{}
		err := json.Unmarshal(b, &q)
		if err == nil {
			queries = append(queries, q)
		}
	case '[':
		isBatch = true
		_ = json.Unmarshal(b, &queries)
	}

	return Request{queries: queries, isBatch: isBatch}
}
