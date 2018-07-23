package updater

import (
	"log"
	"time"

	postgre "github.com/fakihariefnoto/mongomigrate/connection/postgre"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type (
	Pkg interface {
		Update(query string, bn BanNote) error
	}
	pkgUpdater struct{}

	BanNote struct {
		Status     int64
		Note       string
		ProductID  int64
		CreateTime time.Time
		CreateBy   int64
	}
)

var (
	postgreProduct *sqlx.DB
)

const (
	upsertBanNoteQuery = ``
)

func New(dbName string) Pkg {
	var err error
	postgreProduct, err = postgre.Get(dbName)
	if err != nil {
		log.Fatal("error when get db -> ", err)
	}
	return &pkgUpdater{}
}

func (p *pkgUpdater) Update(query string, bn BanNote) error {
	var err error
	if postgreProduct == nil {
		log.Fatal("error when get db, package postgre nill")
	}

	if query == "" {
		query = upsertBanNoteQuery
	}

	_, err = postgreProduct.Exec(query, bn.ProductID, bn.Note, bn.CreateBy, bn.CreateTime, bn.Status)
	if err != nil {
		return errors.Wrap(err, "Err when exec db")
	}

	return nil

}
