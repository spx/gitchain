package api

import (
	"net/http"

	"github.com/inconshreveable/log15"
	"github.com/spx/gitchain/server/context"
)

type NetService struct {
	srv *context.T
	log log15.Logger
}

type JoinArgs struct {
	Host string
}

type JoinReply struct {
}

func (service *NetService) Join(r *http.Request, args *JoinArgs, reply *JoinReply) error {
	service.srv.Router.Pub("/dht/join", args.Host)
	return nil
}
