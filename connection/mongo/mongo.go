package mongo

import (
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

type (
	// Pkg is interface to store new object from New
	Pkg interface {
		NewConn(conn string) (Connect, error)
	}
	pkgMongo struct{}

	// Connect is interface to store object connection in private after get connection
	Connect interface {
		SetDB(aliasName, databaseName string) error
	}
	pDB struct {
		Conn *mgo.Session
	}

	// DB is interface to store object db and collection
	DB interface {
		Collection(collectionName string) (*mgo.Collection, error)
		ListCollection() ([]string, error)
	}
	pMongoCol struct {
		Db  *mgo.Database
		Col *mgo.Collection
	}
)

var (
	dbConn map[string]*mgo.Database
)

// New to create new object of this package
func New() Pkg {
	return &pkgMongo{}
}

// NewConn creates a new session to given connection.
func (p *pkgMongo) NewConn(conn string) (Connect, error) {
	mongo, err := mgo.Dial(conn)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("fail to create %v connecction", conn))
	}

	mongo.SetMode(mgo.Monotonic, true)
	err = mongo.Ping()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("can't ping to %v", conn))
	}

	return &pDB{Conn: mongo}, err
}

// SetDB to get interface of database and collection
func (p *pDB) SetDB(alias, databaseName string) error {
	if dbConn == nil || len(dbConn) == 0 {
		dbConn = make(map[string]*mgo.Database)
	}
	db := p.Conn.DB(databaseName)
	if db == nil {
		return errors.New("Database name not found")
	}
	dbConn[alias] = db
	return nil
}

// Get to get interface of database and collection
func Get(aliasName string) (DB, error) {
	if dbConn == nil || len(dbConn) == 0 {
		return nil, errors.New("Db config for mongo is empty")
	}

	if _, exist := dbConn[aliasName]; !exist {
		return nil, errors.New("Database not found")
	}

	db := dbConn[aliasName]
	if db == nil {
		return nil, errors.New("Database name can't be found")
	}
	return &pMongoCol{Db: db}, nil
}

// Collection is func to get collection from db
func (p *pMongoCol) Collection(collectionName string) (*mgo.Collection, error) {
	colName := p.Db.C(collectionName)
	if colName == nil {
		allCollection, _ := p.ListCollection()
		return nil, errors.New(fmt.Sprintf("Collection name can't be found, here is your colection list name %v", allCollection))
	}
	return colName, nil
}

// ListCollection is func to get all collection
func (p *pMongoCol) ListCollection() ([]string, error) {
	return p.Db.CollectionNames()
}
