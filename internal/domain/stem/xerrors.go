package stem

import "github.com/alan-b-lima/pkg/problem"

var (
	ErrCreate   = problem.Imp(problem.SemanticalError, "stem-create").Message("could not create stem")
	ErrRename   = problem.Imp(problem.SemanticalError, "stem-rename").Message("could not rename stem")
	ErrNotFound = problem.New(problem.NotFound, "stem-not-found", "stem not found", nil, nil)

	ErrTitleTooShort = problem.Imp(problem.SemanticalError, "stem-title-too-short").Format("title must be at least %d characters").Make(NameMinLen)
	ErrTitleTooLong  = problem.Imp(problem.SemanticalError, "stem-title-too-long").Format("title must be at most %d characters").Make(NameMaxLen)
	ErrTitleInvalid  = problem.New(problem.SemanticalError, "stem-title-invalid", "title must consist only of unicode letter and digits, - and /", nil, nil)

	ErrDeleteNonSprout = problem.New(problem.Unauthorized, "stem-delete-non-sprout", "cannot delete a stem that has orders", nil, nil)
)
