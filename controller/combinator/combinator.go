package combinator

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	mongo "github.com/fakihariefnoto/mongomigrate/connection/mongo"
	postgre "github.com/fakihariefnoto/mongomigrate/connection/postgre"
	receiver "github.com/fakihariefnoto/mongomigrate/data/receiver"
	updater "github.com/fakihariefnoto/mongomigrate/data/updater"
	config "github.com/fakihariefnoto/mongomigrate/utils/config"
	formatter "github.com/fakihariefnoto/mongomigrate/utils/formatter"
)

type (
	Pkg interface {
		Runner() error
	}
	pkgCombinator struct {
		cfg config.Config
	}
)

var (
	mongoPkg    mongo.Pkg
	postgrePkg  postgre.Pkg
	receiverPkg receiver.Pkg
	updaterPkg  updater.Pkg

	ErrInitializeNotComplete = errors.New("Please initialize config first")
)

func New() Pkg {
	conf := config.Get()
	return &pkgCombinator{
		cfg: conf,
	}
}

func (p *pkgCombinator) setPostre() {

	postgreProductConn := map[string]string{
		p.cfg.Config.Updater.Connection.Name: p.cfg.Config.Updater.Connection.ConnString,
	}
	postgrePkg = postgre.New(postgreProductConn)
	postgrePkg.NewDB(p.cfg.Config.Updater.Connection.Name)

	updaterPkg = updater.New(p.cfg.Config.Updater.Connection.Name)
}

func (p *pkgCombinator) setMongo() {

	if mongoPkg == nil {
		mongoPkg = mongo.New()
	}

	if receiverPkg == nil {
		receiverPkg = receiver.New()
	}
}

func (p *pkgCombinator) Runner() error {

	p.setPostre()
	p.setMongo()

	if mongoPkg == nil || postgrePkg == nil {
		return ErrInitializeNotComplete
	}

	mongoBanNoteConn := map[string]string{
		"conn": p.cfg.Config.Receiver.Connection.ConnString,
		"db":   p.cfg.Config.Receiver.Connection.Database,
	}

	c, err := mongoPkg.NewConn(mongoBanNoteConn["conn"])
	if err != nil {
		log.Println(err)
	}

	err = c.SetDB(p.cfg.Config.Receiver.Connection.Name, mongoBanNoteConn["db"])
	if err != nil {
		log.Println(err)
	}

	getStartstr := p.cfg.Config.Start
	splitDate := strings.Split(getStartstr, "-")
	yearParse, _ := strconv.ParseInt(splitDate[0], 10, 64)
	monthParse, _ := strconv.ParseInt(splitDate[1], 10, 64)
	dayParse, _ := strconv.ParseInt(splitDate[2], 10, 64)
	start := time.Date(int(yearParse), formatter.IntToMonth(monthParse), int(dayParse), 0, 0, 0, 0, time.UTC)

	getEndstr := p.cfg.Config.End
	splitDate = strings.Split(getEndstr, "-")
	yearParse, _ = strconv.ParseInt(splitDate[0], 10, 64)
	monthParse, _ = strconv.ParseInt(splitDate[1], 10, 64)
	dayParse, _ = strconv.ParseInt(splitDate[2], 10, 64)
	end := time.Date(int(yearParse), formatter.IntToMonth(monthParse), int(dayParse), 0, 0, 0, 0, time.UTC)

	res, err := receiverPkg.Get(p.cfg.Config.Receiver.Connection.Name, p.cfg.Config.Receiver.Connection.Collection, start, end)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, bn := range res {
		err := updaterPkg.Update(p.cfg.Config.Query, updater.BanNote{
			Status:     1,
			Note:       bn.Notes,
			ProductID:  bn.ProductID,
			CreateTime: bn.CreateTime,
			CreateBy:   -1,
		})

		fmt.Println("pID : ", bn.ProductID, " | Notes : ", bn.Notes, " | date :", bn.CreateTime)

		if err != nil {
			log.Fatalln(err)
			return err
		}
	}

	p.cfg.Config.LastEndDate = getEndstr

	return nil
}
