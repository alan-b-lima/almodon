package auth

import (
	"github.com/alan-b-lima/pkg/problem"
)

var (
	ErrUnauthenticated = problem.Imp(problem.Unauthenticated, "unauthenticated").Message("unauthenticated user")

	ErrUnauthorized = problem.Imp(problem.Unauthorized, "unauthorized").Format("actor role %v does not inherit any roles in %v")
	ErrBlocked      = problem.New(problem.Unauthorized, "blocked", "no actor is allowed", nil, nil)
)
