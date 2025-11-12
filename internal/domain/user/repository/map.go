package userrepo

import (
	"cmp"
	"sync"

	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/errors"
	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Map struct {
	uuidIndex  map[uuid.UUID]int
	siapeIndex map[int]int

	repo []user.User
	mu   sync.RWMutex
}

func NewMap() user.Repository {
	repo := Map{
		uuidIndex:  make(map[uuid.UUID]int),
		siapeIndex: make(map[int]int),
	}

	return &repo
}

func (m *Map) List(offset, limit int) (user.Entities, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	lo := clamp(0, offset, len(m.repo))
	hi := clamp(0, offset+limit, len(m.repo))

	if lo >= hi {
		return user.Entities{
			Records:      []user.Entity{},
			TotalRecords: len(m.repo),
		}, nil
	}

	res := make([]user.Entity, hi-lo)
	for i := range m.repo[lo:hi] {
		transform(&res[i], &m.repo[i])
	}

	return user.Entities{
		Offset:       lo,
		Length:       len(res),
		Records:      res,
		TotalRecords: len(m.repo),
	}, nil
}

func (m *Map) Get(uuid uuid.UUID) (user.Entity, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	index, in := m.uuidIndex[uuid]
	if !in {
		return user.Entity{}, xerrors.ErrUserNotFound
	}

	var res user.Entity
	transform(&res, &m.repo[index])
	return res, nil
}

func (m *Map) GetBySIAPE(siape int) (user.Entity, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	index, in := m.siapeIndex[siape]
	if !in {
		return user.Entity{}, xerrors.ErrUserNotFound
	}

	var res user.Entity
	transform(&res, &m.repo[index])
	return res, nil
}

func (m *Map) Create(siape int, name, email, password string, role auth.Role) (user.Entity, error) {
	defer m.mu.Unlock()
	m.mu.Lock()

	u, err := user.New(siape, name, email, password, role)
	if err != nil {
		return user.Entity{}, err
	}

	if _, in := m.siapeIndex[u.SIAPE()]; in {
		return user.Entity{}, xerrors.ErrSiapeTaken
	}

	m.uuidIndex[u.UUID()] = len(m.repo)
	m.siapeIndex[u.SIAPE()] = len(m.repo)
	m.repo = append(m.repo, u)

	var res user.Entity
	transform(&res, &u)
	return res, nil
}

func (m *Map) Patch(uuid uuid.UUID, name, email, password opt.Opt[string], role opt.Opt[auth.Role]) (user.Entity, error) {
	defer m.mu.Unlock()
	m.mu.Lock()

	index, in := m.uuidIndex[uuid]
	if !in {
		return user.Entity{}, xerrors.ErrUserNotFound
	}

	u := m.repo[index]

	err := errors.Join(
		some_then(name, u.SetName),
		some_then(email, u.SetEmail),
		some_then(password, u.SetPassword),
		some_then(role, u.SetRole),
	)
	if err != nil {
		return user.Entity{}, err
	}

	m.repo[index] = u

	var res user.Entity
	transform(&res, &u)
	return res, nil
}

func (m *Map) Delete(uuid uuid.UUID) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	index, in := m.uuidIndex[uuid]
	if !in {
		return nil
	}

	u := &m.repo[index]

	delete(m.uuidIndex, u.UUID())
	delete(m.siapeIndex, u.SIAPE())

	m.repo[index] = m.repo[len(m.repo)-1]
	m.repo = m.repo[:len(m.repo)-1]

	return nil
}

func some_then[F any](src opt.Opt[F], fn func(F) error) error {
	val, ok := src.Unwrap()
	if !ok {
		return nil
	}
	return fn(val)
}

func transform(r *user.Entity, u *user.User) {
	r.UUID = u.UUID()
	r.Name = u.Name()
	r.SIAPE = u.SIAPE()
	r.Email = u.Email()
	r.Password = u.Password()
	r.Role = u.Role()
}

func clamp[T cmp.Ordered](mn, val, mx T) T {
	return min(max(mn, val), mx)
}
