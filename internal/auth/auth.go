package auth

import (
	"fmt"
	"slices"
)

type Authorizer struct {
	classes []Level
}

func Permission(classes ...Level) Authorizer {
	return Authorizer{
		classes: classes,
	}
}

func (auth *Authorizer) Authorize(level Level) bool {
	if slices.Contains(auth.classes, level) {
		return true
	}

	return false
}

func (auth Authorizer) String() string {
	return fmt.Sprint(auth.classes)
}
