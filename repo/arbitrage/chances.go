package arbitrage

import (
	"eos-arb/repo/arbitrage/models"

	"github.com/panyanyany/eos-go"
)

func (r *Repo) GetChances(params GetChancesInput) (chances [][]*models.PathJob, err error) {
	pairHub := r.MakePairHub(params.Pairs)

	pathJob := &models.PathJob{
		Level: 0,
		Out:   &params.BaseAsset,
		Pair:  &models.Pair{},
		//SubPathJobs: []*PathJob{},
	}
	levels := [][]*models.PathJob{
		[]*models.PathJob{pathJob},
	}
	outputLevels := [][]*models.PathJob{}
	for i := 0; i < params.PathDepth; i++ {
		// 遍历第i层所有 pathJob
		// 准备好下一层的空位
		levels = append(levels, []*models.PathJob{})
		outputLevels = append(outputLevels, []*models.PathJob{})
		for _, pathJob := range levels[i] {
			key := pathJob.Out.GetKey()
			//subPathJobs := pathJob.SubPathJobs

			// 当前 pathJob 所有可能的 pair
			for _, pair := range pairHub[key] {
				if pair.HasZero() {
					continue
				}
				// 与上级相同，不处理
				if pair == pathJob.Pair {
					continue
				}
				//// 不到最后一级，不要处理基准货币
				//if i+1 != params.PathDepth && pair.TheOther[key].GetKey() == params.BaseAsset.GetKey() {
				//	continue
				//}
				// 到最后一级了，输出必须是基准货币
				if i+1 == params.PathDepth && pair.TheOther[key].GetKey() != params.BaseAsset.GetKey() {
					continue
				}
				//if pair.TheOther[key].GetKey() != "issue.newdex-8,BTC" {
				//	continue
				//}
				out := pair.GetAmountOut(pathJob.Out.Asset)
				pathJob := models.PathJob{
					Level: i + 1,
					In:    pathJob.Out,
					Out:   &out,
					Pair:  pair,
					//SubPathJobs: []*PathJob{},
					Parent: pathJob,
				}
				//subPathJobs = append(subPathJobs, &pathJob)
				levels[i+1] = append(levels[i+1], &pathJob)
				// 当前输出为基准货币
				if pair.TheOther[key].GetKey() == params.BaseAsset.GetKey() {
					outputLevels[i] = append(outputLevels[i], &pathJob)
				}
			}
			//pathJob.SubPathJobs = subPathJobs
		}
	}

	for _, pathJobs := range outputLevels {
		_chances := r.GetBestChances(GetBestChancesInput{
			PathJobs:  pathJobs,
			BaseAsset: params.BaseAsset,
			MinProfit: params.MinProfit,
		})
		chances = append(chances, _chances...)
	}

	return
}
func (r *Repo) GetBestChances(params GetBestChancesInput) (chances [][]*models.PathJob) {
	for _, pathJob := range params.PathJobs {
		if pathJob.Out.Asset.Amount-params.BaseAsset.Asset.Amount < params.MinProfit.Asset.Amount {
			continue
		}
		paths := []*models.PathJob{}
		path := pathJob
		for path.Parent != nil {
			paths = append(paths, path)
			path = path.Parent
		}
		if len(paths) == 1 {
			panic("wrong paths")
		}
		maxProfit := paths[0].Out.Asset.Sub(paths[len(paths)-1].In.Asset)
		oldProfit := maxProfit
		_ = oldProfit
		// 扩大收益，不停地增加输入资产，找出最大收益
		for {
			// 不要 _paths := paths[:]，相当于直接引用了 paths
			_paths := make([]*models.PathJob, len(paths))
			for j, _path := range paths {
				_paths[j] = &models.PathJob{
					Level: _path.Level,
					In: &eos.ExtendedAsset{
						Asset:    _path.In.Asset,
						Contract: _path.In.Contract,
					},
					Out: &eos.ExtendedAsset{
						Asset:    _path.Out.Asset,
						Contract: _path.Out.Contract,
					},
					Pair: _path.Pair,
					//SubPathJobs: _path.SubPathJobs,
					Parent: _path.Parent,
				}
			}
			_paths[len(_paths)-1].In.Asset = _paths[len(_paths)-1].In.Asset.Add(params.MinProfit.Asset)
			for j := len(_paths) - 1; j >= 0; j-- {
				out := _paths[j].Pair.GetAmountOut(_paths[j].In.Asset)
				_paths[j].Out = &out
				if j > 0 {
					_paths[j-1].In = &out
				}
			}
			profit := _paths[0].Out.Asset.Sub(_paths[len(_paths)-1].In.Asset)
			//PrintPaths(_paths)
			// 同等收益下，尽量多花钱，因为有交易挖矿奖励
			if profit.Amount < maxProfit.Amount {
				break
			}
			//fmt.Printf("%v ~ %v\n", profit, maxProfit)
			maxProfit = profit
			paths = _paths
		}
		_paths := []*models.PathJob{}
		for i := len(paths) - 1; i >= 0; i-- {
			_paths = append(_paths, paths[i])
		}
		paths = _paths
		//PrintPaths(paths)
		chances = append(chances, paths)
	}
	return
}

type GetChancesInput struct {
	Pairs     []*models.Pair
	PathDepth int
	BaseAsset eos.ExtendedAsset
	MinProfit eos.ExtendedAsset
}
type GetBestChancesInput struct {
	PathJobs  []*models.PathJob
	BaseAsset eos.ExtendedAsset
	MinProfit eos.ExtendedAsset
}
