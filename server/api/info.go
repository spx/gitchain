package api

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/inconshreveable/log15"
	"github.com/spx/gitchain/block"
	"github.com/spx/gitchain/server"
	"github.com/spx/gitchain/server/context"
)

type Info struct {
	Mining    server.MiningStatus
	LastBlock *block.Block
	Debug     struct {
		NumGoroutine int
	}
}

func InfoHandler(srv *context.T, log log15.Logger) func(http.ResponseWriter, *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		log := log.New("http")
		lastBlock, err := srv.DB.GetLastBlock()
		if err != nil {
			log.Error("error serving /info", "err", err)
		}
		info := Info{
			Mining:    server.GetMiningStatus(),
			LastBlock: lastBlock,
		}
		info.Debug.NumGoroutine = runtime.NumGoroutine()
		json, err := json.Marshal(info)
		if err != nil {
			log.Error("error serving /info", "err", err)
		}
		resp.Header().Add("Content-Type", "application/json")
		resp.Write(json)
	}

}
