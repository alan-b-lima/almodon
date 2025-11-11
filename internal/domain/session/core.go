package session

import (
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

const _MaxAge = 10 * time.Minute

func Get(repo Getter, uuid uuid.UUID) (Entity, error) {
	return repo.Get(uuid)
}

// TODO: verify validity of _MaxAge and turn it to an internal error
func Create(repo Creater, uuid uuid.UUID) (Entity, error) {
	return repo.Create(uuid, _MaxAge)
}

func CreateWithMaxAge(repo Creater, uuid uuid.UUID, maxAge time.Duration) (Entity, error) {
	return repo.Create(uuid, maxAge)
}

// TODO: verify validity of _MaxAge and turn it to an internal error
func Update(repo Updater, uuid uuid.UUID) (Entity, error) {
	return repo.Update(uuid, _MaxAge)
}

func UpdateWithMaxAge(repo Updater, uuid uuid.UUID, maxAge time.Duration) (Entity, error) {
	return repo.Update(uuid, maxAge)
}
