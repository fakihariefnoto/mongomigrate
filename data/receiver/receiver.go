package receiver

import (
	"log"
	"time"

	mongo "github.com/fakihariefnoto/mongomigrate/connection/mongo"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type (
	Pkg interface {
		Get(dbName, ColName string, start, end time.Time) ([]BanNote, error)
	}
	pkgReceiver struct{}

	BanNote struct {
		ID         bson.ObjectId `bson:"_id"`
		ProductID  int64         `bson:"p_id"`
		Notes      string        `bson:"notes"`
		CreateTime time.Time     `bson:"create_time"`
	}
)

var (
	mongoDB mongo.DB
)

func New() Pkg {
	return &pkgReceiver{}
}

func (p *pkgReceiver) Get(dbName, ColName string, start, end time.Time) ([]BanNote, error) {
	var err error
	if mongoDB == nil {
		mongoDB, err = mongo.Get(dbName)
		if err != nil {
			return nil, errors.Wrap(err, "Err when get mongo db")
		}
	}

	var result []BanNote

	col, err := mongoDB.Collection(ColName)
	if err != nil {
		return nil, errors.Wrap(err, "Err when get collection db")
	}

	newBS := bson.M{
		"create_time": bson.M{
			"$lt":  end,
			"$gte": start,
		},
	}

	err = col.Find(newBS).All(&result)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil

}
