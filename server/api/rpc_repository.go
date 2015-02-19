package api

import (
	"encoding/hex"
	"net/http"

	"github.com/inconshreveable/log15"
	"github.com/spx/gitchain/repository"
	"github.com/spx/gitchain/server/context"
)

type repo struct {
	Name             string
	Status           string
	NameAllocationTx string
}

type RepositoryService struct {
	srv *context.T
	log log15.Logger
}

type ListRepositoriesArgs struct {
}

type ListRepositoriesReply struct {
	Repositories []repo
}

var status = map[int]string{
	repository.PENDING: "pending",
	repository.ACTIVE:  "active",
}

func (service *RepositoryService) ListRepositories(r *http.Request, args *ListRepositoriesArgs, reply *ListRepositoriesReply) error {
	repos := service.srv.DB.ListRepositories()
	for i := range repos {
		r, err := service.srv.DB.GetRepository(repos[i])
		if err != nil {
			return err
		}
		reply.Repositories = append(reply.Repositories,
			repo{
				Name:             r.Name,
				Status:           status[r.Status],
				NameAllocationTx: hex.EncodeToString(r.NameAllocationTx)})
	}
	return nil
}
