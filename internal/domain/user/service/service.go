package userserve

import (
	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/session"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/pkg/errors"
	"github.com/alan-b-lima/almodon/pkg/opt"
)

type Service struct {
	Repo     user.Repository
	Sessions session.Repository
}

func NewService(users user.Repository, sessions session.Repository) user.Service {
	return &Service{
		Repo:     users,
		Sessions: sessions,
	}
}

func (s *Service) List(act auth.Actor, req user.ListRequest) (user.ListResponse, error) {
	res, err := user.List(s.Repo, req.Offset, req.Limit)
	if err != nil {
		return user.ListResponse{}, err
	}

	lres := user.ListResponse{
		Offset:       res.Offset,
		Length:       res.Length,
		Records:      make([]user.Response, res.Length),
		TotalRecords: res.TotalRecords,
	}
	for i := range res.Records {
		transformP(&lres.Records[i], &res.Records[i])
	}

	return lres, nil
}

func (s *Service) Get(act auth.Actor, req user.GetRequest) (user.Response, error) {
	res, err := user.Get(s.Repo, req.UUID)
	if err != nil {
		return user.Response{}, err
	}

	return transform(&res), nil
}

func (s *Service) GetBySIAPE(act auth.Actor, req user.GetBySIAPERequest) (user.Response, error) {
	res, err := user.GetBySIAPE(s.Repo, req.SIAPE)
	if err != nil {
		return user.Response{}, err
	}

	return transform(&res), err
}

func (s *Service) Create(act auth.Actor, req user.CreateRequest) (user.Response, error) {
	role, ok := auth.FromString(req.Role)
	if !ok {
		return user.Response{}, &errors.Error{}
	}

	res, err := user.Create(s.Repo, req.SIAPE, req.Name, req.Email, req.Password, role)
	if err != nil {
		return user.Response{}, err
	}

	return transform(&res), nil
}

func (s *Service) Patch(act auth.Actor, req user.PatchRequest) (user.Response, error) {
	var role opt.Opt[auth.Role]
	if val, ok := req.Role.Unwrap(); ok {
		if l, ok := auth.FromString(val); ok {
			role = opt.Some(l)
		}
	}

	res, err := user.Patch(s.Repo, req.UUID, req.Name, req.Email, req.Password, role)
	if err != nil {
		return user.Response{}, err
	}

	return transform(&res), nil
}

func (s *Service) Delete(act auth.Actor, req user.DeleteRequest) error {
	return user.Delete(s.Repo, req.UUID)
}

func (s *Service) Authenticate(req user.AuthRequest) (user.AuthResponse, error) {
	res, err := user.Authenticate(s.Repo, s.Sessions, req.SIAPE, req.Password)
	if err != nil {
		return user.AuthResponse{}, err
	}

	return user.AuthResponse(res), nil
}

func (s *Service) Actor(req user.ActorRequest) (auth.Actor, error) {
	return user.Actor(s.Repo, s.Sessions, req.Session)
}

func transform(e *user.Entity) user.Response {
	return user.Response{
		UUID:  e.UUID,
		SIAPE: e.SIAPE,
		Name:  e.Name,
		Email: e.Email,
		Role:  e.Role.String(),
	}
}

func transformP(r *user.Response, e *user.Entity) {
	r.UUID = e.UUID
	r.SIAPE = e.SIAPE
	r.Name = e.Name
	r.Email = e.Email
	r.Role = e.Role.String()
}
