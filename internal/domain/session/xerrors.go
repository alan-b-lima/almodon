package session

import "github.com/alan-b-lima/pkg/problem"

var (
	ErrCreate   = problem.Imp(problem.SemanticalError, "session-create").Message("could not create session")
	ErrUpdate   = problem.Imp(problem.SemanticalError, "session-update").Message("could not update session")
	ErrNotFound = problem.New(problem.NotFound, "session-not-found", "session not found", nil, nil)

	ErrInvalidToken = problem.Imp(problem.Malformed, "invalid-token").Format("token cannot be parsed into %d byte array").Make(TokenLen)
	ErrTooLong      = problem.New(problem.SemanticalError, "session-too-long", "session too long", nil, map[string]any{"hard": HardTimeout, "idle": IdleTimeout})
)
