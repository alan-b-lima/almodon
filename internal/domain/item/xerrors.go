package item

import "github.com/alan-b-lima/pkg/problem"

var (
	ErrCreate   = problem.Imp(problem.SemanticalError, "item-create").Message("could not create item")
	ErrUpdate   = problem.Imp(problem.SemanticalError, "item-update").Message("could not update item")
	ErrNotFound = problem.New(problem.NotFound, "item-not-found", "item not found", nil, nil)

	ErrAmountNegative   = problem.New(problem.SemanticalError, "item-amount-negative", "amount cannot be negative", nil, nil)
	ErrUnitCostNegative = problem.New(problem.SemanticalError, "item-unit-cost-negative", "unit cost cannot be negative", nil, nil)
)
