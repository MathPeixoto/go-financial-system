package gapi

import (
	"fmt"
	"github.com/MathPeixoto/go-financial-system/worker"

	db "github.com/MathPeixoto/go-financial-system/db/sqlc"
	"github.com/MathPeixoto/go-financial-system/pb"
	"github.com/MathPeixoto/go-financial-system/token"
	"github.com/MathPeixoto/go-financial-system/util"
)

type Server struct {
	pb.UnimplementedBankServer
	config      util.Config
	store       db.Store
	tokenMaker  token.Maker
	distributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, distributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:      config,
		store:       store,
		tokenMaker:  tokenMaker,
		distributor: distributor,
	}

	return server, nil
}
