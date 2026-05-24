package lotitem

import "github.com/alan-b-lima/pkg/problem"

var (
	ErrCreate   = problem.Imp(problem.SemanticalError, "lotitem-create").Message("could not create lot item")
	ErrUpdate   = problem.Imp(problem.SemanticalError, "lotitem-update").Message("could not update lot item")
	ErrNotFound = problem.New(problem.NotFound, "lotitem-not-found", "lot item not found", nil, nil)
)
