package domain

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"math/bits"
	"os"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	auths "github.com/alan-b-lima/almodon/internal/domain/auth/resource"
	authserve "github.com/alan-b-lima/almodon/internal/domain/auth/service"

	"github.com/alan-b-lima/almodon/internal/domain/item"
	items "github.com/alan-b-lima/almodon/internal/domain/item/resource"
	itemserve "github.com/alan-b-lima/almodon/internal/domain/item/service"
	itemstore "github.com/alan-b-lima/almodon/internal/domain/item/store"

	"github.com/alan-b-lima/almodon/internal/domain/material"
	materials "github.com/alan-b-lima/almodon/internal/domain/material/resource"
	materialserve "github.com/alan-b-lima/almodon/internal/domain/material/service"
	materialstore "github.com/alan-b-lima/almodon/internal/domain/material/store"

	"github.com/alan-b-lima/almodon/internal/domain/promotion"
	promotions "github.com/alan-b-lima/almodon/internal/domain/promotion/resource"
	promotionserve "github.com/alan-b-lima/almodon/internal/domain/promotion/service"
	promotionstore "github.com/alan-b-lima/almodon/internal/domain/promotion/store"

	"github.com/alan-b-lima/almodon/internal/domain/session"
	sessionserve "github.com/alan-b-lima/almodon/internal/domain/session/service"
	sessionstore "github.com/alan-b-lima/almodon/internal/domain/session/store"

	"github.com/alan-b-lima/almodon/internal/domain/user"
	users "github.com/alan-b-lima/almodon/internal/domain/user/resource"
	userserve "github.com/alan-b-lima/almodon/internal/domain/user/service"
	userstore "github.com/alan-b-lima/almodon/internal/domain/user/store"

	"github.com/alan-b-lima/almodon/pkg/closer"
	"github.com/alan-b-lima/pkg/scheduler"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/term"
)

type Domain struct {
	Database  *sql.DB
	Scheduler *scheduler.Scheduler

	Stores    Stores
	Cores     Cores
	Services  Services
	Resources Resources

	Bundle closer.Bundle
}

type (
	Stores struct {
		Items      item.Store
		Materials  material.Store
		Promotions promotion.Store
		Sessions   session.Store
		Users      user.Store
	}

	Cores struct {
		Auths      *authserve.Core
		Items      *itemserve.Core
		Materials  *materialserve.Core
		Promotions *promotionserve.Core
		Sessions   *sessionserve.Core
		Users      *userserve.Core
	}

	Services struct {
		Auths      auth.Service
		Items      item.Service
		Materials  material.Service
		Promotions promotion.Service
		Sessions   session.Service
		Users      user.Service
	}

	Resources struct {
		Auth       *auths.Resource
		Items      *items.Resource
		Materials  *materials.Resource
		Promotions *promotions.Resource
		Users      *users.Resource
	}
)

var (
	ErrInvalidOptComb = errors.New("invalid option combination")
	ErrNoRootUser     = errors.New("no root user found in database")
)

func New(opts ...Option) (*Domain, error) {
	var bundle closer.Bundle
	var err error

	defer func(err *error) {
		if *err != nil {
			bundle.Close()
		}
	}(&err)

	opt, err := Condense(opts...)
	if err != nil {
		return nil, err
	}

	db, err := MountSQLiteDB(opt)
	if err != nil {
		return nil, err
	}
	bundle.Bundle(db)

	scheduler := scheduler.New()
	bundle.BundleFunc(scheduler.Stop)
	scheduler.Start()

	stores := MountStores(db)
	cores := MountCores(scheduler, stores)
	services := MountServices(cores)
	resources := MountResources(services)

	domain := Domain{
		Database:  db,
		Scheduler: scheduler,

		Stores:    stores,
		Cores:     cores,
		Services:  services,
		Resources: resources,

		Bundle: bundle,
	}

	if opt&Structure != 0 {
		if err = PrepareStructure(db); err != nil {
			return nil, fmt.Errorf("prepare structure: %w", err)
		}
	}

	if opt&RootUser != 0 {
		if err = AssertRootUser(cores.Users, opt); err != nil {
			return nil, fmt.Errorf("assert root user: %w", err)
		}
	}

	if opt&Publish != 0 {
		if err = PublishForScheduler(cores.Sessions, cores.Promotions); err != nil {
			return nil, fmt.Errorf("publish for scheduler: %w", err)
		}
	}

	return &domain, nil
}

type Option uint64

const (
	Structure Option = 1 << iota
	RootUser
	Publish

	InMemory
	Interactive

	_Built   = 1 << 63
	_Default = _Built | Structure | RootUser | Publish | Interactive
)

func Condense(opts ...Option) (Option, error) {
	switch len(opts) {
	case 0:
		return _Default, nil

	case 1:
		opt := opts[0]
		if opt&_Built == 0 {
			return 0, ErrInvalidOptComb
		}

		return opt, nil
	}

	final := _Default

	for _, opt := range opts {
		switch bits.OnesCount64(uint64(opt)) {
		case 1:
			final |= opt
		case 63:
			final &= opt
		default:
			return 0, ErrInvalidOptComb
		}
	}

	return final, nil
}

const SQLiteDriver = "sqlite3"

func MountSQLiteDB(opts ...Option) (*sql.DB, error) {
	opt, err := Condense(opts...)
	if err != nil {
		return nil, err
	}

	if opt&InMemory != 0 {
		return OpenSQLiteDBInMemory()
	}

	db, err := OpenSQLiteDB(
		"data/almodon.db",
		"../data/almodon.db",
		"../../data/almodon.db",
	)
	if opt&Interactive != 0 && errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll("data", 0o775); err != nil {
			return nil, err
		}

		db, err = sql.Open(SQLiteDriver, "data/almodon.db")
	}
	if err != nil {
		return nil, err
	}

	return db, nil
}

func MountStores(db *sql.DB) Stores {
	stores := Stores{
		Items:      itemstore.New(db),
		Materials:  materialstore.New(db),
		Promotions: promotionstore.New(db),
		Sessions:   sessionstore.New(db),
		Users:      userstore.New(db),
	}

	return stores
}

func MountCores(scheduler *scheduler.Scheduler, stores Stores) Cores {
	cores := Cores{
		Items:     itemserve.New(stores.Items),
		Materials: materialserve.New(stores.Materials),
		Sessions:  sessionserve.New(stores.Sessions, scheduler, nil),
		Users:     userserve.New(stores.Users),
	}

	cores.Auths = authserve.New(cores.Users, cores.Sessions)
	cores.Promotions = promotionserve.New(stores.Promotions, cores.Users, scheduler)

	return cores
}

func MountServices(cores Cores) Services {
	authed := Services{
		Auths:      cores.Auths,
		Items:      itemserve.NewGate(cores.Items, cores.Auths),
		Materials:  materialserve.NewGate(cores.Materials, cores.Auths),
		Promotions: promotionserve.NewGate(cores.Promotions, cores.Auths),
		Sessions:   cores.Sessions,
		Users:      userserve.NewGate(cores.Users, cores.Auths),
	}

	return authed
}

func MountResources(services Services) Resources {
	resources := Resources{
		Auth:       auths.New(services.Auths),
		Items:      items.New(services.Items),
		Materials:  materials.New(services.Materials),
		Promotions: promotions.New(services.Promotions),
		Users:      users.New(services.Users),
	}

	return resources
}

func OpenSQLiteDB(names ...string) (*sql.DB, error) {
	var name string
	for _, n := range names {
		_, err := os.Stat(n)
		if !errors.Is(err, os.ErrNotExist) {
			name = n
			break
		}
	}
	if name == "" {
		return nil, os.ErrNotExist
	}

	return sql.Open(SQLiteDriver, name)
}

func OpenSQLiteDBInMemory() (*sql.DB, error) {
	return sql.Open(SQLiteDriver, ":memory:")
}

var preconditions = [...]string{
	"PRAGMA foreign_keys = ON",
}

var scripts = [...]string{
	itemstore.Script,
	materialstore.Script,
	promotionstore.Script,
	sessionstore.Script,
	userstore.Script,
}

func PrepareStructure(db *sql.DB) error {
	for _, op := range preconditions {
		if _, err := db.Exec(op); err != nil {
			return err
		}
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, script := range scripts {
		if _, err := tx.Exec(script); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}

			return err
		}
	}

	return tx.Commit()
}

func AssertRootUser(users user.Service, opts ...Option) error {
	opt, err := Condense(opts...)
	if err != nil {
		return err
	}

	if _, err := users.GetBySIAPE(context.TODO(), "0000000"); err != user.ErrNotFound {
		return err
	}

	if opt&Interactive == 0 {
		return ErrNoRootUser
	}

	user := user.Create{
		SIAPE: "0000000",
		Email: "noreply@ufvjm.edu.br",
		Role:  auth.Maintainer,
	}

	// we should be careful with this,
	// our Windows friends might get upset.
	fmt.Print("\033[?1049h\033[1;1H")
	defer fmt.Print("\033[?1049l")

	fmt.Println("Creating root user...")

	fmt.Print("Insert name for root user: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		err := scanner.Err()
		if err == io.EOF {
			return io.ErrUnexpectedEOF
		}

		return err
	}

	user.Name = scanner.Text()

	fmt.Printf("Insert password for %q: ", user.Name)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		if err == io.EOF {
			return io.ErrUnexpectedEOF
		}

		return err
	}
	fmt.Println()

	user.Password = string(password)

	_, err = users.Create(context.TODO(), user)
	return err
}

func PublishForScheduler(sessions *sessionserve.Core, promotions *promotionserve.Core) error {
	ctx := context.TODO()

	var (
		err0 = sessions.Publish(ctx)
		err1 = promotions.Publish(ctx)
	)

	return errors.Join(err0, err1)
}
