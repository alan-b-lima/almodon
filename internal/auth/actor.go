package auth

import "github.com/alan-b-lima/almodon/pkg/uuid"

type Actor struct {
	user uuid.UUID
	role Role
}

func NewLogged(user uuid.UUID, role Role) Actor {
	return Actor{
		user: user,
		role: role,
	}
}

func NewUnlogged() Actor {
	return Actor{role: Unlogged}
}

func (act *Actor) User() uuid.UUID {
	return act.user
}

func (act *Actor) Role() Role {
	return act.role
}
