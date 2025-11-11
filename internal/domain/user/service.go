package user

import (
	"github.com/alan-b-lima/almodon/internal/auth"
)

type Service interface {
	List(act auth.Actor, req ListRequest) (ListResponse, error)
	Get(act auth.Actor, req GetRequest) (Response, error)
	GetBySIAPE(act auth.Actor, req GetBySIAPERequest) (Response, error)
	Create(act auth.Actor, req CreateRequest) (Response, error)
	Patch(act auth.Actor, req PatchRequest) (Response, error)
	Delete(act auth.Actor, req DeleteRequest) error
	Authenticate(req AuthRequest) (AuthResponse, error)
	Actor(req ActorRequest) (auth.Actor, error)
}
