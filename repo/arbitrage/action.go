package arbitrage

import (
	"fmt"

	"eos-arb/repo/arbitrage/dex_adapter"
	"eos-arb/repo/arbitrage/models"
	"eos-arb/repo/arbitrage/util"

	"github.com/cihub/seelog"
	"github.com/panyanyany/eos-go"
)

func (r *Repo) RunChances(chances [][]*models.PathJob, minProfit eos.Asset) (err error) {
	for _, paths := range chances {
		var actions []*eos.Action
		groups := util.PathGroupByDex(paths)
		for j, group := range groups {
			var pairs []*models.Pair
			for _, path := range group {
				pairs = append(pairs, path.Pair)
			}
			out := group[len(group)-1].Out.Asset
			if j == len(groups)-1 {
				out = groups[0][0].In.Asset.Add(minProfit)
			}
			action, err := group[0].Pair.Dex.(dex_adapter.IDex).GetMultiSwapAction(dex_adapter.PushMultiSwapActionInput{
				Pairs: pairs,
				In:    group[0].In.Asset,
				//Out:   group[len(group)-1].Out.Asset,
				Out: out,
			})
			if err != nil {
				msg := fmt.Sprintf("GetMultiSwapAction: %v", err)
				fmt.Println(msg)
				seelog.Errorf(msg)
				break
			}
			actions = append(actions, action)
		}
		go func(actions []*eos.Action) {
			_, err = r.Api.PushActions(actions)
			if err != nil {
				err = fmt.Errorf("r.Api.PushActions: %s", err)
				fmt.Println(err)
				seelog.Error(err)
			}
		}(actions)
	}
	return
}
