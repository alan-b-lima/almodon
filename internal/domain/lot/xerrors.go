package lot

import "github.com/alan-b-lima/pkg/problem"

var (
	ErrCreate   = problem.Imp(problem.SemanticalError, "lot-create").Message("could not create lot")
	ErrUpdate   = problem.Imp(problem.SemanticalError, "lot-update").Message("could not update lot")
	ErrNotFound = problem.New(problem.NotFound, "lot-not-found", "lot not found", nil, nil)

	ErrSupplierTooLong = problem.New(problem.SemanticalError, "lot-supplier-too-long", "supplier is too long", nil, nil)
	ErrNoteTooLong     = problem.New(problem.SemanticalError, "lot-note-too-long", "note is too long", nil, nil)

	ErrDeleteSigned = problem.New(problem.Unauthorized, "lot-delete-signed", "cannot delete signed lot", nil, nil)
)
