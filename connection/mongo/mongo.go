package mongo

import (
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

type (
	// Pkg is interface to store new object from New
	Pkg interface {
		NewConn(conn string) (MongoDB, error)
	}
	pkgMongo struct{}

	// MongoDB is interface to store object connection in private after get connection
	MongoDB interface {
		GetDB(databaseName string) (MongoCollection, error)
	}
	pDB struct {
		Conn *mgo.Session
	}

	// MongoCollection is interface to store object db and collection
	MongoCollection interface {
		Collection(collectionName string) (*mgo.Collection, error)
		ListCollection() ([]string, error)
	}
	pMongoCol struct {
		Db  *mgo.Database
		Col *mgo.Collection
	}
)

// func New to create new object of this package
func New() Pkg {
	return &pkgMongo{}
}

// NewConn creates a new session to given connection.
func (p *pkgMongo) NewConn(conn string) (MongoDB, error) {
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

// func GetDB to get interface of database and collection
func (p *pDB) GetDB(databaseName string) (MongoCollection, error) {
	db := p.Conn.DB(databaseName)
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
