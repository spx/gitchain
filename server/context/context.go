package context

import (
	"os"
	"path/filepath"

	"github.com/inconshreveable/log15"
	"github.com/spx/gitchain/db"
	"github.com/spx/gitchain/server/config"
	"github.com/tuxychandru/pubsub"
)

type T struct {
	Config *config.T
	DB     *db.T
	Log    log15.Logger
	Router *pubsub.PubSub
}

func (srv *T) Init() error {
	err := os.MkdirAll(srv.Config.General.DataPath, os.ModeDir|0700)
	if err != nil {
		return err
	}
	database, err := db.NewDB(filepath.Join(srv.Config.General.DataPath, "db"))
	if err != nil {
		return err
	}
	srv.DB = database
	srv.Log = log15.New()
	srv.Router = pubsub.New(100)
	return nil
}
