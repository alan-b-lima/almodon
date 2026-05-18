package stem

import "github.com/alan-b-lima/pkg/problem"

var (
	ErrCreate   = problem.Imp(problem.SemanticalError, "stem-create").Message("could not create stem")
	ErrRename   = problem.Imp(problem.SemanticalError, "stem-rename").Message("could not rename stem")
	ErrNotFound = problem.New(problem.NotFound, "stem-not-found", "stem not found", nil, nil)

	ErrTitleInvalid = problem.New(problem.SemanticalError, "stem-title-invalid", "title must consist only of unicode letter and digits, - and /", nil, nil)

)
