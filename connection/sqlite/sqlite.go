package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // only call init gosqlite
	"github.com/pkg/errors"
)

type (

	// Pkg is interface to store new object after calling New func
	Pkg interface {
		NewDB(dbName ...string) error
		AddDB(name, path string) error
	}
	pkgDatabase struct {
		Connection map[string]string
	}
)

var (
	dbconn map[string]*sql.DB

	// ErrDBNill is var to define error if db nill
	ErrDBNill = errors.New("database is nil")
)

// New func is func for creating new object of this package
func New(conf map[string]string) Pkg {
	return &pkgDatabase{
		Connection: conf,
	}
}

// NewDB is func to connect to each db that registered in New
func (p *pkgDatabase) NewDB(dbName ...string) error {
	if p.Connection == nil || len(p.Connection) == 0 {
		return errors.New("No config found, please config first")
	}
	if dbconn == nil || len(dbconn) == 0 {
		dbconn = make(map[string]*sql.DB)
	}
	if len(dbName) == 0 {
		return errors.New("Please fill database connection")
	}

	for _, name := range dbName {
		if _, exist := p.Connection[name]; !exist {
			return errors.New("Database didn't exist, please add config first")
		}
		p.AddDB(name, p.Connection[name])
	}
	return nil
}

// AddDB is func to add db and connect it and save it to
func (p *pkgDatabase) AddDB(name, connString string) error {
	db, err := openConnection(connString)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Fail to connect to %v with connString %v", name, connString))
	}
	dbconn[name] = db
	return nil
}

// Get Database after creating connection with New func
func Get(name string) (*sql.DB, error) {
	if dbconn == nil || len(dbconn) == 0 {
		return nil, ErrDBNill
	}
	if _, exist := dbconn[name]; !exist {
		return nil, ErrDBNill
	}
	return dbconn[name], nil
}

func openConnection(filePath string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", filePath)
	if err != nil {
		return db, errors.Wrapf(err, "can't open sqlite to %v", filePath)
	}

	return db, err
}
