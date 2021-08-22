package arbitrage

import (
	"fmt"
	"time"

	"eos-arb/repo/arbitrage/models"

	"github.com/cihub/seelog"
	"github.com/panyanyany/eos-go"
)

func (r *Repo) RunTask(params RunTaskInput) {
	arb := r
	ts_start := time.Now()
	baseSym := params.BaseSymbol
	baseContract := params.BaseContract
	minProfit := eos.Asset{Amount: eos.Int64(params.MinProfitAmount), Symbol: baseSym}

	chances, err := arb.GetChances(GetChancesInput{
		Pairs:     params.Pairs,
		PathDepth: 4,
		BaseAsset: eos.ExtendedAsset{
			Asset:    eos.Asset{Amount: eos.Int64(params.BaseAmount), Symbol: baseSym},
			Contract: baseContract,
		},
		MinProfit: eos.ExtendedAsset{
			Asset:    eos.Asset{Amount: 1, Symbol: baseSym},
			Contract: baseContract,
		},
	})
	if err != nil {
		err = fmt.Errorf("arb.GetChances: %w", err)
		seelog.Error(err)
	}
	chances = arb.FilterChances(chances, minProfit)
	ts_diff := time.Now().Sub(ts_start)
	seelog.Debugf("calc time cost: %.3fs", ts_diff.Seconds())
	for _, paths := range chances {
		r.Printer.PrintPaths(paths)
	}
	arb.RunChances(chances, minProfit)
}

type RunTaskInput struct {
	BaseSymbol      eos.Symbol
	BaseContract    eos.AccountName
	BaseAmount      int64
	MinProfitAmount int64
	Pairs           []*models.Pair
}
