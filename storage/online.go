package storage

import (
	"context"

	"github.com/ipfs/go-bitswap"
	bsnet "github.com/ipfs/go-bitswap/network"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/p2p"
)

// Online initializes a p2p host for the underlying blockservice.
func (s *Store) Online(ctx context.Context) error {
	priv, err := s.ReadKey()
	if err != nil {
		return err
	}

	host, router, err := p2p.NewHost(ctx, priv)
	if err != nil {
		return err
	}

	net := bsnet.NewFromIpfsHost(host, router)
	exc := bitswap.New(ctx, net, s.bstore)

	bserv := blockservice.New(s.bstore, exc)
	dag := merkledag.NewDAGService(bserv)

	s.Dag = dag
	s.Host = host
	s.Router = router

	return nil
}