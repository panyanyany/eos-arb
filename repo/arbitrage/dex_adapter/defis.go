package dex_adapter

import (
	"fmt"

	"eos-arb/repo/arbitrage/models"
	"eos-arb/repo/defis/defisswapcnt"

	"github.com/panyanyany/eos-go"
	"gorm.io/gorm"
)

type Defis struct {
	BaseAdapter
	Repo *defisswapcnt.Repo
}

func (r *Defis) Init(db *gorm.DB) (err error) {
	r.Repo = &defisswapcnt.Repo{
		Api: r.Api,
		Db:  db,
	}
	return
}
func (r *Defis) GetAllPairs() (pairs []*models.Pair, err error) {
	var _pairs []*defisswapcnt.Pair
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
				Contract: _pair.Contract0,
				Symbol:   _pair.Sym0,
			},
			Reserve:   _pair.Reserve0,
			PriceLast: _pair.Price0Last.Value,
		}
		token1 := &models.TokenDetail{
			TokenMeta: models.TokenMeta{
				Contract: _pair.Contract1,
				Symbol:   _pair.Sym1,
			},
			Reserve:   _pair.Reserve1,
			PriceLast: _pair.Price1Last.Value,
		}
		pairs = append(pairs,
			&models.Pair{
				Id:      _pair.Mid,
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
func (r *Defis) GetSwapCmd(params GetSwapCmdInput) (cmd string, err error) {
	cmd, err = r.Repo.GetSwapCmd(defisswapcnt.GetSwapCmdInput{
		Pair: params.Pair.RawPair.(*defisswapcnt.Pair),
		From: params.From,
		In:   params.In,
		Out:  params.Out,
	})
	return
}
func (r *Defis) GetMultiSwapCmd(params GetMultiSwapCmdInput) (cmd string, err error) {
	pairs := []*defisswapcnt.Pair{}
	for _, pair := range params.Pairs {
		pairs = append(pairs, pair.RawPair.(*defisswapcnt.Pair))
	}
	cmd, err = r.Repo.GetMultiSwapCmd(defisswapcnt.GetMultiSwapCmdInput{
		Pairs: pairs,
		From:  params.From,
		In:    params.In,
		Out:   params.Out,
	})
	return
}
func (r *Defis) PushMultiSwapAction(params PushMultiSwapActionInput) (resp *eos.PushTransactionFullResp, err error) {
	pairs := []*defisswapcnt.Pair{}
	for _, pair := range params.Pairs {
		pairs = append(pairs, pair.RawPair.(*defisswapcnt.Pair))
	}
	resp, err = r.Repo.PushMultiSwapAction(defisswapcnt.PushMultiSwapActionInput{
		Pairs: pairs,
		In:    params.In,
		Out:   params.Out,
	})
	return
}
func (r *Defis) GetMultiSwapAction(params PushMultiSwapActionInput) (action *eos.Action, err error) {
	pairs := []*defisswapcnt.Pair{}
	for _, pair := range params.Pairs {
		pairs = append(pairs, pair.RawPair.(*defisswapcnt.Pair))
	}
	action, err = r.Repo.GetMultiSwapAction(defisswapcnt.PushMultiSwapActionInput{
		Pairs: pairs,
		In:    params.In,
		Out:   params.Out,
	})
	return
}
func (r *Defis) GetName() string {
	return "Defis"
}
