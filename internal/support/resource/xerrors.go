package resource

import "github.com/alan-b-lima/pkg/problem"

var (
	ErrBadUUID        = problem.New(problem.Malformed, "bad-uuid", "invalid uuid", nil, nil)
	ErrBadInteger     = problem.New(problem.Malformed, "bad-integer", "invalid integer", nil, nil)
	ErrBadQueryParams = problem.Imp(problem.Malformed, "bad-query-params").Message("invalid query params")

	ErrNoContentType          = problem.New(problem.UnsupportedContentType, "no-content-type", "content type must be informed", nil, nil)
	ErrUnsupportedContentType = problem.Imp(problem.UnsupportedContentType, "unsupported-content-type").Format("content type must be %s")
	ErrNotAcceptable          = problem.Imp(problem.UnsupportedAcceptable, "unsupported-acceptable").Format("client does not accept %s")

	ErrJSON = problem.Imp(problem.Malformed, "json-error").Message("unexpected error ocurred while processing json")

	ErrNotFound = problem.Imp(problem.NotFound, "resource-not-found").Format("resource %s not found")
)
