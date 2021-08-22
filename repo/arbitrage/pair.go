package arbitrage

import (
	"fmt"
	"sync"

	"eos-arb/repo/arbitrage/dex_adapter"
	"eos-arb/repo/arbitrage/models"

	"github.com/cihub/seelog"
)

func (r *Repo) MakePairHub(pairs []*models.Pair) (pairHub map[string][]*models.Pair) {
	pairHub = make(map[string][]*models.Pair)
	for _, pair := range pairs {
		_, found0 := pairHub[pair.Token0.GetKey()]
		if !found0 {
			pairHub[pair.Token0.GetKey()] = make([]*models.Pair, 0)
		}
		pairHub[pair.Token0.GetKey()] = append(pairHub[pair.Token0.GetKey()], pair)

		_, found1 := pairHub[pair.Token1.GetKey()]
		if !found1 {
			pairHub[pair.Token1.GetKey()] = make([]*models.Pair, 0)
		}
		pairHub[pair.Token1.GetKey()] = append(pairHub[pair.Token1.GetKey()], pair)
	}
	return
}
func (r *Repo) FilterPairs(pairs []*models.Pair) (newPairs []*models.Pair) {
	contractBlackList := []string{
		"joker.eos", // 老是显示 swap too small
		"token.bank", // 吞币
		"eoscccdotcom", // 转账燃烧
		"yupyupxtoken", // balance error
		"token.yup", // balance error
	}
	for _, pair := range pairs {
		skip := false
		for _, contract := range contractBlackList {
			if pair.Token0.Contract == contract || pair.Token1.Contract == contract {
				skip = true
			}
		}
		if skip {
			continue
		}
		newPairs = append(newPairs, pair)
	}
	return
}
func (r *Repo) GetAllPairs(dexes []dex_adapter.IDex) (pairs []*models.Pair) {
	pairs = []*models.Pair{}
	wg := sync.WaitGroup{}
	cRes := make(chan *GetPairsResult)
	done := make(chan bool)
	for _, dex := range dexes {
		wg.Add(1)

		go func(dex dex_adapter.IDex) {
			defer wg.Done()
			_pairs, err := dex.GetAllPairs()
			cRes <- &GetPairsResult{
				Pairs: _pairs,
				Error: err,
			}
		}(dex)
	}
	go func() {
		for res := range cRes {
			if res.Error != nil {
				err := fmt.Errorf(": %w", res.Error)
				seelog.Error(err)
			} else {
				pairs = append(pairs, res.Pairs...)
			}
		}
		done <- true
	}()
	wg.Wait()
	close(cRes)
	<-done
	return
}

type GetPairsResult struct {
	Pairs []*models.Pair
	Error error
}
