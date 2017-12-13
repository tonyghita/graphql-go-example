package errors

import graphql "github.com/neelance/graphql-go/errors"

type slicer interface {
	Slice() []error
}

type indexedCauser interface {
	Index() int
	Cause() error
}

func Expand(errs []*graphql.QueryError) []*graphql.QueryError {
	expanded := make([]*graphql.QueryError, 0, len(errs))

	for _, err := range errs {
		switch t := err.ResolverError.(type) {
		case slicer:
			for _, e := range t.Slice() {
				qe := &graphql.QueryError{
					Message:   err.Message,
					Locations: err.Locations,
					Path:      err.Path,
				}

				if ic, ok := e.(indexedCauser); ok {
					qe.Path = append(qe.Path, ic.Index())
					qe.Message = ic.Cause().Error()
				}

				expanded = append(expanded, qe)
			}
		default:
			expanded = append(expanded, err)
		}
	}

	return expanded
}
