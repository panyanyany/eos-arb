package arbitrage

import (
	"sort"

	"eos-arb/repo/arbitrage/models"

	"github.com/panyanyany/eos-go"
)

func (r *Repo) FilterChances(chances [][]*models.PathJob, minProfit eos.Asset) (filtered [][]*models.PathJob) {
	if len(chances) == 0 {
		return
	}

	sort.Slice(chances, func(i, j int) bool {
		iPathJob := chances[i][len(chances[i])-1]
		jPathJob := chances[j][len(chances[j])-1]
		return jPathJob.Out.Asset.Amount < iPathJob.Out.Asset.Amount
	})

	pairExists := make(map[int]bool)
	for _, paths := range chances {
		skip := false
		profit := paths[len(paths)-1].Out.Asset.Sub(paths[0].In.Asset)
		if profit.Amount < minProfit.Amount {
			continue
		}
		for _, path := range paths {
			//if j == 0 {
			//	continue
			//}
			_, found := pairExists[path.Pair.Id]
			if found {
				skip = true
				break
			}
			pairExists[path.Pair.Id] = true
		}
		if skip {
			continue
		}
		filtered = append(filtered, paths)
	}

	return
}
