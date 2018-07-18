package sqlite

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //only call pq.init()
	"github.com/pkg/errors"
)

var (
	dbconn   map[string]dbReplication
	ErrDBNil = errors.New("database is nil")
)

// Connect creates connection from given config.
// If all connections to database succeed, connect will return nil error and user can access
// database connections using Get function.
func Connect(config interface{}) (err error) {
	if dbconn == nil || len(dbconn) == 0 {
		dbconn = make(map[string]dbReplication)
	}

	// cast interface into config struct
	cfg := config.(map[string]*struct {
		Master string
		Slave  []string
	})

	// loop each config into db connection
	for name, conns := range cfg {
		db := dbReplication{}

		db.Master, err = openConnection(conns.Master)
		if err != nil {
			return errors.Wrap(err, "master")
		}

		var slaveconns []*sqlx.DB
		for _, v := range conns.Slave {
			con, err := openConnection(v)
			if err != nil {
				return errors.Wrap(err, "slave")
			}
			slaveconns = append(slaveconns, con)
		}

		db.Slave = slaveconns
		dbconn[name] = db
	}

	return nil
}

// AddSQLXConn is for add new connection using sqlx.DB obj
func AddSQLXConn(connname string, master *sqlx.DB, slaves ...*sqlx.DB) {
	if dbconn == nil || len(dbconn) == 0 {
		dbconn = make(map[string]dbReplication)
	}

	dbconn[connname] = dbReplication{
		Master: master,
		Slave:  slaves,
	}
}

// AddSQLConn is for add new connection using sql.DB obj
func AddSQLConn(connname string, master *sql.DB, slaves ...*sql.DB) {
	masterx := sqlx.NewDb(master, "postgres")

	slavesx := []*sqlx.DB{}
	for _, v := range slaves {
		slavesx = append(slavesx, sqlx.NewDb(v, "postgres"))
	}

	AddSQLXConn(connname, masterx, slavesx...)
}

// Get database connection based on database name.
// Use replication value "master" to get connection to master database.
func Get(name, replication string) (*sqlx.DB, error) {
	if dbconn == nil || len(dbconn) == 0 {
		return nil, fmt.Errorf("Database connection is not initialized yet")
	}

	if db, ok := dbconn[name]; ok {
		var err error
		if replication == "master" {
			if db.Master == nil {
				err = errors.Wrapf(ErrDBNil, "%s Master", name)
			}
			return db.Master, err
		}

		if len(db.Slave) == 1 {
			if db.Slave[0] == nil {
				err = errors.Wrapf(ErrDBNil, "%s Slave %d", name, 0)
			}
			return db.Slave[0], err
		}

		if len(db.Slave) == 0 {
			return nil, errors.Wrapf(ErrDBNil, "%s Slave", name)
		}

		// if has more than 1 slave, return random slave
		rand.Seed(time.Now().UTC().UnixNano())
		randIndex := rand.Intn(len(db.Slave))
		if db.Slave[randIndex] == nil {
			err = errors.Wrapf(ErrDBNil, "%s Slave %d", name, randIndex)
		}
		return db.Slave[randIndex], err

	}

	return nil, fmt.Errorf("Database connection " + name + " doesn't exist")
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

// Ping all database connections
func Ping() (errMap map[string]error) {
	errMap = make(map[string]error)
	for label, conn := range dbconn {
		if err := conn.Master.Ping(); err != nil {
			errMap[fmt.Sprintf("[Master]%v", label)] = err
		}
		for i, v := range conn.Slave {
			if err := v.Ping(); err != nil {
				errMap[fmt.Sprintf("[Slave-%v]%v", i, label)] = err
			}
		}
	}

	return errMap
}
