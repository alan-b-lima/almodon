package userserve

import (
	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/internal/xerrors"
)

type AuthService struct {
	service   user.Service
	hierarchy auth.Hierarchy
}

func New(service user.Service) user.Service {
	return &AuthService{
		service:   service,
		hierarchy: auth.DefaultHierarchy,
	}
}

var (
	permStrictChief = auth.Permit(auth.Chief)
	permChief       = auth.Permit(auth.Promoted)
	permAdmin       = auth.Permit(auth.Admin)
	permLogged      = auth.Permit(auth.User)
	permPermissive  = auth.Permit(auth.Unlogged)
)

func (s *AuthService) List(act auth.Actor, req user.ListRequest) (user.ListResponse, error) {
	if err := service.Authorize(permStrictChief, act); err != nil {
		return user.ListResponse{}, err
	}

	return s.service.List(act, req)
}

func (s *AuthService) Get(act auth.Actor, req user.GetRequest) (user.Response, error) {
	if act.User() == req.UUID {
		goto Do
	}

	if err := service.Authorize(permStrictChief, act); err != nil {
		return user.Response{}, err
	}

Do:
	return s.service.Get(act, req)
}

func (s *AuthService) GetBySIAPE(act auth.Actor, req user.GetBySIAPERequest) (user.Response, error) {
	res, err := s.service.GetBySIAPE(act, req)
	if err != nil {
		return user.Response{}, err
	}

	if act.User() == res.UUID {
		goto Do
	}

	if err := service.Authorize(permStrictChief, act); err != nil {
		return user.Response{}, err
	}

Do:
	return res, nil
}

func (s *AuthService) Create(act auth.Actor, req user.CreateRequest) (user.Response, error) {
	if err := service.Authorize(permStrictChief, act); err != nil {
		return user.Response{}, err
	}

	return s.service.Create(act, req)
}

// TODO: Patch has serious security ploblems, rationalize and break them into smaller use cases
func (s *AuthService) Patch(act auth.Actor, req user.PatchRequest) (user.Response, error) {
	if r, ok := req.Role.Unwrap(); ok {
		role, ok := auth.FromString(r)
		if !ok {
			goto Continue
		}

		if !s.hierarchy(role, act.Role()) {
			return user.Response{}, xerrors.ErrUnpriviledUserPromotion.New(act.Role(), role)
		}
	}

Continue:
	if act.User() == req.UUID {
		goto Do
	}

	if err := service.Authorize(permStrictChief, act); err != nil {
		return user.Response{}, err
	}

Do:
	return s.service.Patch(act, req)
}

func (s *AuthService) Delete(act auth.Actor, req user.DeleteRequest) error {
	if act.User() == req.UUID {
		goto Do
	}

	if err := service.Authorize(permStrictChief, act); err != nil {
		return err
	}

Do:
	return s.service.Delete(act, req)
}

func (s *AuthService) Authenticate(req user.AuthRequest) (user.AuthResponse, error) {
	return s.service.Authenticate(req)
}

func (s *AuthService) Actor(req user.ActorRequest) (auth.Actor, error) {
	return s.service.Actor(req)
}
