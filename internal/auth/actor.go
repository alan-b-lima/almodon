package auth

import "github.com/alan-b-lima/almodon/pkg/uuid"

// Actor is a system entity that can access services through it's
// role. It is related the user entity, as a Actor is a shrinked
// version of an user.
type Actor struct {
	user uuid.UUID
	role Role
}

// NewLogged creates a new actor. This function does not check
// whether the fact is real.
func NewLogged(user uuid.UUID, role Role) Actor {
	return Actor{
		user: user,
		role: role,
	}
}

// NewUnlogged creates a new unlogged actor. It's also equivalent to
// the zero value of [Actor].
func NewUnlogged() Actor {
	return Actor{role: Unlogged}
}

// User returns the uuid of the user as an actor.
func (act *Actor) User() uuid.UUID {
	return act.user
}

// Role returns the role of the user as an actor.
func (act *Actor) Role() Role {
	return act.role
}
