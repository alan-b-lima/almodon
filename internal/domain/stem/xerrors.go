package stem

import "github.com/alan-b-lima/pkg/problem"

var ErrNotFound = problem.New(problem.NotFound, "stem-not-found", "stem not found", nil, nil)
