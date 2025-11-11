package sessionrepo

import (
	"sync"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/session"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/heap"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Map struct {
	uuidIndex   map[uuid.UUID]int
	userIndex   map[uuid.UUID]int
	expiresHeap sleepqueue

	repo []session.Session
	mu   sync.RWMutex
}

func NewMap() session.Repository {
	repo := Map{
		uuidIndex: make(map[uuid.UUID]int),
		userIndex: make(map[uuid.UUID]int),
		expiresHeap: sleepqueue{
			new:    make(chan ess, 32),
			update: make(chan ess, 32),
			cancel: make(chan struct{}, 1),
		},
	}

	go flush(&repo)

	return &repo
}

func (m *Map) Get(uuid uuid.UUID) (session.Entity, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	index, in := m.uuidIndex[uuid]
	if !in {
		return session.Entity{}, xerrors.ErrSessionNotFound
	}

	s := m.repo[index]
	if time.Now().After(s.Expires()) {
		m.delete(s.UUID())
		return session.Entity{}, xerrors.ErrSessionNotFound
	}

	var res session.Entity
	transform(&res, &m.repo[index])
	return res, nil
}

func (m *Map) Create(user uuid.UUID, maxAge time.Duration) (session.Entity, error) {
	m.mu.Lock()

	if index, in := m.userIndex[user]; in {
		s := m.repo[index]
		m.delete(s.UUID())
	}

	s, err := session.New(user, maxAge)
	if err != nil {
		return session.Entity{}, err
	}

	m.uuidIndex[s.UUID()] = len(m.repo)
	m.userIndex[s.User()] = len(m.repo)
	m.repo = append(m.repo, s)

	// unlock before channel send to avoid blocking resources
	m.mu.Unlock()

	m.expiresHeap.new <- ess{s.UUID(), s.Expires()}

	var res session.Entity
	transform(&res, &s)
	return res, nil
}

func (m *Map) Update(uuid uuid.UUID, maxAge time.Duration) (session.Entity, error) {
	m.mu.Lock()

	index, in := m.uuidIndex[uuid]
	if !in {
		return session.Entity{}, xerrors.ErrSessionNotFound
	}

	s := &m.repo[index]

	if err := s.SetMaxAge(maxAge); err != nil {
		return session.Entity{}, err
	}

	// unlock before channel send to avoid blocking resources
	m.mu.Unlock()

	m.expiresHeap.update <- ess{s.UUID(), s.Expires()}

	var res session.Entity
	transform(&res, s)
	return res, nil
}

func (m *Map) delete(uuid uuid.UUID) error {
	index, in := m.uuidIndex[uuid]
	if !in {
		return nil
	}

	s := &m.repo[index]

	delete(m.uuidIndex, s.UUID())
	delete(m.userIndex, s.User())

	m.repo[index] = m.repo[len(m.repo)-1]
	m.repo = m.repo[:len(m.repo)-1]
	return nil
}

func (m *Map) tryDelete(uuid uuid.UUID) bool {
	index, in := m.uuidIndex[uuid]
	if !in {
		return true
	}

	s := &m.repo[index]
	if time.Now().Before(s.Expires()) {
		return false
	}

	delete(m.uuidIndex, s.UUID())
	delete(m.userIndex, s.User())

	m.repo[index] = m.repo[len(m.repo)-1]
	m.repo = m.repo[:len(m.repo)-1]
	return true
}

func transform(r *session.Entity, s *session.Session) {
	r.UUID = s.UUID()
	r.User = s.User()
	r.Expires = s.Expires()
}

func flush(m *Map) {
	h := m.expiresHeap
	var garbage int // count of garbage in the heap

	for {
		var after <-chan time.Time
		if h.heap.Len() > 0 {
			delay := time.Until(h.heap.Peek().expires)
			after = time.After(delay)
		}

		select {
		case <-h.cancel:
			return

		case es := <-h.new:
			h.heap.Push(es)

		case es := <-h.update:
			if garbage < 128 {
				h.heap.Push(es)
				garbage++
				continue
			}

			m.mu.RLock()
			ss := make([]ess, 0, len(m.repo))

			for _, s := range m.repo {
				ss = append(ss, ess{s.UUID(), s.Expires()})
			}
			m.mu.RUnlock()

			nheap := heap.Make[ess](len(m.repo))
			nheap.PushMany(ss...)

		case <-after:
			es := h.heap.Pop()

			m.mu.Lock()
			if m.tryDelete(es.session) {
				garbage--
			}
			m.mu.Unlock()
		}
	}
}

type sleepqueue struct {
	heap   heap.Heap[ess]
	new    chan ess
	update chan ess
	cancel chan struct{}
}

type ess struct {
	session uuid.UUID
	expires time.Time
}

func (o0 ess) Less(o1 ess) bool { return o0.expires.Before(o1.expires) }
