package util

import (
	"eos-arb/repo/arbitrage/models"
)

func PathGroupByDex(paths []*models.PathJob) (groups []PathGroup) {
	lastPaths := []*models.PathJob{}
	for _, path := range paths {
		if len(lastPaths) == 0 {
			lastPaths = []*models.PathJob{path}
		} else {
			if lastPaths[len(lastPaths)-1].Pair.Dex == path.Pair.Dex {
				lastPaths = append(lastPaths, path)
			} else {
				groups = append(groups, lastPaths)
				lastPaths = []*models.PathJob{path}
			}
		}
	}
	groups = append(groups, lastPaths)
	return
}

type PathGroup []*models.PathJob
