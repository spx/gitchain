package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/inconshreveable/log15"
	"github.com/spx/gitchain/block"
	"github.com/spx/gitchain/server/context"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WebsocketHandler(srv *context.T, log log15.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.New("cmp", "websocket")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("error upgrading websocket connetion", "err", err)
			return
		}

		ch := srv.Router.Sub("/block")
		defer srv.Router.Unsub(ch)

	loop:
		select {
		case blki := <-ch:
			if blk, ok := blki.(*block.Block); ok {
				encoded, err := json.Marshal(blk)
				if err != nil {
					log.Error("error encoding block", "err", err)
					return
				}
				if err = conn.WriteMessage(websocket.TextMessage, encoded); err != nil {
					log.Error("error sending data", "err", err)
					return
				}
			}
		}
		goto loop

	}
}
