package postgre

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //only call pq.init()
	"github.com/pkg/errors"
)

var (
	dbconn map[string]*sqlx.DB

	// ErrDBNill is var to define error if db nill
	ErrDBNill = errors.New("database is nil")
)

type (
	// Pkg is interface to store new object after calling New func
	Pkg interface {
		NewDB(dbName ...string) error
		AddDB(name, connString string) error
		Ping() (errMap map[string]error)
		KeepAlive(t time.Duration) error
		AddWebhook() error
	}
	pkgDatabase struct {
		Config Config
	}

	// Config is struct for storing db information
	Config struct {
		Connection   map[string]string
		isUseWebhock bool
	}
)

/*
	New will be passing with func
	can add config weebhok
	add new connection db
	this param maps is temporary
*/

// New func is func for creating new object of this package
func New(conf map[string]string) Pkg {
	config := Config{
		Connection: conf,
	}
	return &pkgDatabase{
		Config: config,
	}
}

// NewDB is func to connect to each db that registered in New
func (p *pkgDatabase) NewDB(dbName ...string) error {
	if p.Config.Connection == nil || len(p.Config.Connection) == 0 {
		return errors.New("No config found, please config first")
	}
	if dbconn == nil || len(dbconn) == 0 {
		dbconn = make(map[string]*sqlx.DB)
	}
	if len(dbName) == 0 {
		return errors.New("Please fill database connection")
	}

	for _, name := range dbName {
		if _, exist := p.Config.Connection[name]; !exist {
			return errors.New("Database didn't exist, please add config first")
		}
		p.AddDB(name, p.Config.Connection[name])
	}
	return nil
}

// AddDB is func to add db and connect it and save it to
func (p *pkgDatabase) AddDB(name, connString string) error {
	if dbconn == nil || len(dbconn) == 0 {
		dbconn = make(map[string]*sqlx.DB)
	}
	db, err := openConnection(connString)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Fail to connect to %v with connString %v", name, connString))
	}
	dbconn[name] = db
	return nil
}

// Ping is func to ping all database connections and get map error
func (p *pkgDatabase) Ping() (errMap map[string]error) {
	errMap = make(map[string]error)
	for name, db := range dbconn {
		err := db.Ping()
		if err != nil {
			errMap[name] = err
		}
	}
	return errMap
}

func (p *pkgDatabase) KeepAlive(t time.Duration) error {
	return nil
}

func (p *pkgDatabase) AddWebhook() error {
	return nil
}

// Get Database after creating connection with New func
func Get(name string) (*sqlx.DB, error) {
	if dbconn == nil || len(dbconn) == 0 {
		return nil, ErrDBNill
	}
	if _, exist := dbconn[name]; !exist {
		return nil, ErrDBNill
	}
	return dbconn[name], nil
}

func openConnection(connstring string) (db *sqlx.DB, err error) {
	db, err = sqlx.Connect("postgres", connstring)
	if err != nil {
		return db, errors.Wrapf(err, "connect to %v", connstring)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(30)

	return db, err
}
