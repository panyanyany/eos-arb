package printer

import (
	"fmt"

	"eos-arb/repo/arbitrage/dex_adapter"
	"eos-arb/repo/arbitrage/models"
	"eos-arb/repo/arbitrage/util"

	"github.com/cihub/seelog"
	"github.com/panyanyany/eos-go"
)

type PathPrinter struct {
	ShowDex      bool
	ShowContract bool
	ShowCmd      bool
	From         eos.Name
}

func (r *PathPrinter) PrintPaths(paths []*models.PathJob) {
	msg := ""
	for j := 0; j < len(paths); j++ {
		contract := ""
		if r.ShowContract {
			contract = paths[j].In.GetKey() + ":"
		}
		dex := ""
		if r.ShowDex {
			dex = paths[j].Pair.Dex.(dex_adapter.IDex).GetName() + ":"
		}
		msg += fmt.Sprintf("[%v%v%v] -> ",
			dex,
			contract,
			paths[j].In.Asset.String(),
		)
	}

	contract := ""
	if r.ShowContract {
		contract = paths[len(paths)-1].In.GetKey() + ":"
	}
	dex := ""
	if r.ShowDex {
		dex = paths[len(paths)-1].Pair.Dex.(dex_adapter.IDex).GetName() + ":"
	}
	msg += fmt.Sprintf("[%v%v%v] = %v",
		dex,
		contract,
		paths[len(paths)-1].Out.Asset.String(),
		paths[len(paths)-1].Out.Asset.Sub(paths[0].In.Asset),
	)
	seelog.Info(msg)

	//fmt.Println()
	if r.ShowCmd {
		for _, group := range util.PathGroupByDex(paths) {
			r.PrintCmd(group)
		}
	}
}
func (r *PathPrinter) PrintCmd(paths []*models.PathJob) {
	var cmd string
	var err error
	pairs := []*models.Pair{}
	for _, _path := range paths {
		pairs = append(pairs, _path.Pair)
	}
	cmd, err = paths[0].Pair.Dex.(dex_adapter.IDex).GetMultiSwapCmd(dex_adapter.GetMultiSwapCmdInput{
		Pairs: pairs,
		From:  r.From,
		In:    paths[0].In.Asset,
		Out:   paths[len(paths)-1].Out.Asset,
	})
	if err != nil {
		err = fmt.Errorf("GetMultiSwapCmd: %w", err)
		seelog.Error(err)
		return
	}
	seelog.Infof("\t%v", cmd)

}
