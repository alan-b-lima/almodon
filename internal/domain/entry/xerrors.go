package entry

import "github.com/alan-b-lima/pkg/problem"

var ErrNotFound = problem.New(problem.NotFound, "order-not-found", "order not found", nil, nil)
