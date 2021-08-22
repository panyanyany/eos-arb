package dex_adapter

import (
	"fmt"

	"eos-arb/repo/arbitrage/models"
	"eos-arb/repo/dolphin/dolphinsswap"

	"github.com/panyanyany/eos-go"
	"gorm.io/gorm"
)

type Dolphin struct {
	BaseAdapter
	Repo *dolphinsswap.Repo
}

func (r *Dolphin) Init(db *gorm.DB) (err error) {
	r.Repo = &dolphinsswap.Repo{
		Api: r.Api,
		Db:  db,
	}
	return
}
func (r *Dolphin) GetAllPairs() (pairs []*models.Pair, err error) {
	var _pairs []*dolphinsswap.Pair
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
				Contract: _pair.Tokens[0].Symbol.Contract,
				Symbol:   _pair.Tokens[0].Symbol.Symbol,
			},
			Reserve:   _pair.Tokens[0].Reserve,
			PriceLast: 0,
		}
		token1 := &models.TokenDetail{
			TokenMeta: models.TokenMeta{
				Contract: _pair.Tokens[1].Symbol.Contract,
				Symbol:   _pair.Tokens[1].Symbol.Symbol,
			},
			Reserve:   _pair.Tokens[1].Reserve,
			PriceLast: 0,
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
func (r *Dolphin) GetSwapCmd(params GetSwapCmdInput) (cmd string, err error) {
	cmd, err = r.Repo.GetSwapCmd(dolphinsswap.GetSwapCmdInput{
		Pair: params.Pair.RawPair.(*dolphinsswap.Pair),
		From: params.From,
		In:   params.In,
		Out:  params.Out,
	})
	return
}
func (r *Dolphin) GetMultiSwapCmd(params GetMultiSwapCmdInput) (cmd string, err error) {
	pairs := []*dolphinsswap.Pair{}
	for _, pair := range params.Pairs {
		pairs = append(pairs, pair.RawPair.(*dolphinsswap.Pair))
	}
	cmd, err = r.Repo.GetMultiSwapCmd(dolphinsswap.GetMultiSwapCmdInput{
		Pairs: pairs,
		From:  params.From,
		In:    params.In,
		Out:   params.Out,
	})
	return
}
func (r *Dolphin) PushMultiSwapAction(params PushMultiSwapActionInput) (resp *eos.PushTransactionFullResp, err error) {
	pairs := []*dolphinsswap.Pair{}
	for _, pair := range params.Pairs {
		pairs = append(pairs, pair.RawPair.(*dolphinsswap.Pair))
	}
	resp, err = r.Repo.PushMultiSwapAction(dolphinsswap.PushMultiSwapActionInput{
		Pairs: pairs,
		In:    params.In,
		Out:   params.Out,
	})
	return
}
func (r *Dolphin) GetMultiSwapAction(params PushMultiSwapActionInput) (action *eos.Action, err error) {
	pairs := []*dolphinsswap.Pair{}
	for _, pair := range params.Pairs {
		pairs = append(pairs, pair.RawPair.(*dolphinsswap.Pair))
	}
	action, err = r.Repo.GetMultiSwapAction(dolphinsswap.PushMultiSwapActionInput{
		Pairs: pairs,
		In:    params.In,
		Out:   params.Out,
	})
	return
}
func (r *Dolphin) GetName() string {
	return "Dolphin"
}
