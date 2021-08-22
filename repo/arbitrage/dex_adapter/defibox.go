package dex_adapter

import (
	"fmt"

	"eos-arb/repo/arbitrage/models"
	"eos-arb/repo/defibox/swap_defi"

	"github.com/panyanyany/eos-go"
	"gorm.io/gorm"
)

type Defibox struct {
	BaseAdapter
	Repo *swap_defi.Repo
}

func (r *Defibox) Init(db *gorm.DB) (err error) {
	r.Repo = &swap_defi.Repo{
		Api: r.Api,
		Db:  db,
	}
	return
}

func (r *Defibox) GetAllPairs() (pairs []*models.Pair, err error) {
	var _pairs []*swap_defi.Pair
	//_pairs, err = r.Repo.GetAllPairs()
	_pairs, err = r.Repo.GetAllPairsBatch()
	if err != nil {
		err = fmt.Errorf("r.Repo.GetAllPairs(): %w", err)
		return
	}

	pairs = make([]*models.Pair, 0, len(_pairs))

	for _, _pair := range _pairs {
		token0 := &models.TokenDetail{
			TokenMeta: models.TokenMeta{
				Contract: _pair.Token0.Contract,
				Symbol:   _pair.Token0.Symbol,
			},
			Reserve:   _pair.Reserve0,
			PriceLast: _pair.Price0Last.Value,
		}
		token1 := &models.TokenDetail{
			TokenMeta: models.TokenMeta{
				Contract: _pair.Token1.Contract,
				Symbol:   _pair.Token1.Symbol,
			},
			Reserve:   _pair.Reserve1,
			PriceLast: _pair.Price1Last.Value,
		}
		pairs = append(pairs,
			&models.Pair{
				Id:      _pair.ID,
				Token0:  token0,
				Token1:  token1,
				Dex:     r,
				RawPair: _pair,
				TheOther: map[string]*models.TokenDetail{
					token0.GetKey(): token1,
					token1.GetKey(): token0,
				},
			},
		)
	}

	return
}
func (r *Defibox) GetSwapCmd(params GetSwapCmdInput) (cmd string, err error) {
	cmd, err = r.Repo.GetSwapCmd(swap_defi.GetSwapCmdInput{
		Pair: params.Pair.RawPair.(*swap_defi.Pair),
		From: params.From,
		In:   params.In,
		Out:  params.Out,
	})
	return
}
func (r *Defibox) GetMultiSwapCmd(params GetMultiSwapCmdInput) (cmd string, err error) {
	pairs := []*swap_defi.Pair{}
	for _, pair := range params.Pairs {
		pairs = append(pairs, pair.RawPair.(*swap_defi.Pair))
	}
	cmd, err = r.Repo.GetMultiSwapCmd(swap_defi.GetMultiSwapCmdInput{
		Pairs: pairs,
		From:  params.From,
		In:    params.In,
		Out:   params.Out,
	})
	return
}
func (r *Defibox) PushMultiSwapAction(params PushMultiSwapActionInput) (resp *eos.PushTransactionFullResp, err error) {
	pairs := []*swap_defi.Pair{}
	for _, pair := range params.Pairs {
		pairs = append(pairs, pair.RawPair.(*swap_defi.Pair))
	}
	resp, err = r.Repo.PushMultiSwapAction(swap_defi.PushMultiSwapActionInput{
		Pairs: pairs,
		In:    params.In,
		Out:   params.Out,
	})
	return
}
func (r *Defibox) GetMultiSwapAction(params PushMultiSwapActionInput) (action *eos.Action, err error) {
	pairs := []*swap_defi.Pair{}
	for _, pair := range params.Pairs {
		pairs = append(pairs, pair.RawPair.(*swap_defi.Pair))
	}
	action, err = r.Repo.GetMultiSwapAction(swap_defi.PushMultiSwapActionInput{
		Pairs: pairs,
		In:    params.In,
		Out:   params.Out,
	})
	return
}
func (r *Defibox) GetName() string {
	return "Defibox"
}
