package auth

import "github.com/alan-b-lima/almodon/pkg/uuid"

type Actor struct {
	user  uuid.UUID
	level Level
}

func NewLogged(user uuid.UUID, level Level) Actor {
	return Actor{
		user:  user,
		level: level,
	}
}

func NewUnlogged() Actor {
	return Actor{level: Unlogged}
}

func (act *Actor) User() uuid.UUID {
	return act.user
}

func (act *Actor) Level() Level {
	return act.level
}
