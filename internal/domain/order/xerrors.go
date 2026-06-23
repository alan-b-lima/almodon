package order

import "github.com/alan-b-lima/pkg/problem"

var (
	ErrCreate   = problem.Imp(problem.SemanticalError, "order-create").Message("could not create order")
	ErrNotFound = problem.New(problem.NotFound, "order-not-found", "order not found", nil, nil)

	ErrSignatureFailed = problem.New(problem.Unauthorized, "order-signature-failed", "invalid password", nil, nil)

	ErrGenerationalFailure = problem.New(problem.LostUpdate, "stock-generational-failure", "stock has been modified since last read", nil, nil)
)
