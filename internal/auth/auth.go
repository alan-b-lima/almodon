package auth

import "fmt"

// Permission represents an authorization requirement.
type Permission struct {
	classes   []Role
	hierarchy func(Role, Role) bool
}

// Permit creates a Permission that authorizes any of the given roles
// according to the default hierarchy.
func Permit(classes ...Role) Permission {
	return Permission{
		classes:   classes,
		hierarchy: DefaultHierarchy,
	}
}

// PermitHierarchy creates a Permission that authorizes any of the
// given roles according to the provided hierarchy.
func PermitHierarchy(hierarchy Hierarchy, classes ...Role) Permission {
	return Permission{
		classes:   classes,
		hierarchy: hierarchy,
	}
}

// Authorize returns whether the given role is authorized by the
// Permission.
func (auth *Permission) Authorize(role Role) bool {
	for _, class := range auth.classes {
		if auth.hierarchy(class, role) {
			return true
		}
	}

	return false
}

// String returns the string representation of the Permission.
func (auth Permission) String() string {
	return fmt.Sprint(auth.classes)
}
