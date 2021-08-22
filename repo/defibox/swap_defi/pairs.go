package swap_defi

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"eos-arb/model"
	"eos-arb/repo/constants/str_const"

	"github.com/cihub/seelog"
	"github.com/panyanyany/eos-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repo) GetPairs(startPairId int) (pairs []*Pair, err error) {
	var resp *eos.GetTableRowsResp
	resp, _, err = r.Api.GetTableRows(eos.GetTableRowsRequest{
		Code:       Code,
		Scope:      Code,
		Table:      "pairs",
		JSON:       true,
		LowerBound: strconv.Itoa(startPairId),
		Limit:      99,
	})
	if err != nil {
		err = fmt.Errorf("get pairs: %w", err)
		return
	}
	pairs = make([]*Pair, 0)
	err = json.Unmarshal(resp.Rows, &pairs)
	if err != nil {
		seelog.Debug(string(resp.Rows))
		err = fmt.Errorf("unmarshal pairs: %w", err)
		return
	}
	return
}
func (r *Repo) GetAllPairs() (allRows []*Pair, err error) {
	var resp *eos.GetTableRowsResp
	var lowerBound string
	exists := make(map[int]bool)
	//var all_rows []*Sasset

	for {
		var rows []*Pair
		params := eos.GetTableRowsRequest{
			JSON:       true,
			Code:       Code,
			Scope:      Code,
			Table:      "pairs",
			LowerBound: lowerBound,
			UpperBound: "",
			Limit:      399,
		}
		resp, _, err = r.Api.GetTableRows(params)
		if err != nil {
			err = fmt.Errorf("GetTableRows of nftoken sassets: %w", err)
			return
		}
		err = json.Unmarshal(resp.Rows, &rows)
		if err != nil {
			seelog.Debugf("rows: %v", string(resp.Rows))
			err = fmt.Errorf("unmarshal: %w", err)
			return
		}
		for _, row := range rows {
			_, found := exists[row.ID]
			if found {
				continue
			}
			exists[row.ID] = true
			allRows = append(allRows, row)
		}
		if !resp.More {
			break
		}
		lowerBound = strconv.Itoa(rows[len(rows)-1].ID)
	}

	return
}
func (r *Repo) GetAllPairsBatch() (allRows []*Pair, err error) {
	mdIdList := model.PairIdList{}
	err = r.Db.Where("contract = ?", Code).First(&mdIdList).Error
	if err == gorm.ErrRecordNotFound || time.Now().Sub(mdIdList.UpdatedAt).Hours() > 7*24 {
		allRows, err = r.GetAllPairs()
		if err != nil {
			err = fmt.Errorf("r.GetAllPairs(): %w", err)
			return
		}
		seelog.Infof("init pairs: %v", len(allRows))
		ids := []int{}
		for _, row := range allRows {
			ids = append(ids, row.ID)
		}
		sort.Ints(ids)
		ids2 := []string{}
		for _, id := range ids {
			ids2 = append(ids2, strconv.Itoa(id))
		}
		idStr := strings.Join(ids2, ",")
		mdIdList.Contract = Code
		mdIdList.IdList = idStr
		mdIdList.Total = len(ids)

		err = r.Db.Clauses(clause.OnConflict{UpdateAll: true}).Save(&mdIdList).Error
		if err != nil {
			err = fmt.Errorf("Save(mdIdList): %w", err)
			return
		}
		return
	}
	idStrList := strings.Split(mdIdList.IdList, ",")
	if len(idStrList) <= 10 {
		err = fmt.Errorf("idStrList to small: %v", len(idStrList))
		return
	}
	exists := make(map[int]bool)

	var wg sync.WaitGroup
	step := 399
	queue := make(chan *GetTableRowsResult)
	resList := []*GetTableRowsResult{}

	done := make(chan bool)
	go func() {
		for res := range queue {
			resList = append(resList, res)
		}
		done <- true
	}()

	for i := 0; i < len(idStrList); i += step {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			lowerBound := idStrList[i]

			var resp *eos.GetTableRowsResp
			params := eos.GetTableRowsRequest{
				JSON:       true,
				Code:       Code,
				Scope:      Code,
				Table:      "pairs",
				LowerBound: lowerBound,
				UpperBound: "",
				Limit:      399,
			}
			baseUrl := ""
			resp, baseUrl, err = r.Api.GetTableRows(params)
			queue <- &GetTableRowsResult{
				Resp:   resp,
				Error:  err,
				Params: params,
				Url:    baseUrl,
			}
		}(i)
	}
	wg.Wait()
	close(queue)
	<-done

	for _, res := range resList {
		//seelog.Infof("url=%v, lowerBound=%v, list=%v", res.Url, res.Params.LowerBound, len(resList))
		resp := res.Resp
		err = res.Error
		if err != nil {
			err = fmt.Errorf("r.Api.GetTableRows: %w", err)
			//seelog.Infof("url=%v, lowerBound=%v, err=%v", res.Url, res.Params.LowerBound, err)
			seelog.Error(err)
			err = nil
			continue
		}
		var rows []*Pair
		err = json.Unmarshal(resp.Rows, &rows)
		if err != nil {
			seelog.Debugf("rows: %v", string(resp.Rows))
			err = fmt.Errorf("unmarshal: %w", err)
			seelog.Error(err)
			err = nil
			continue
		}
		for _, row := range rows {
			_, found := exists[row.ID]
			if found {
				continue
			}
			exists[row.ID] = true
			allRows = append(allRows, row)
		}
		//seelog.Infof("url=%v, lowerBound=%v, len=%v", res.Url, res.Params.LowerBound, len(rows))
	}

	//seelog.Infof("get all pairs: %v", len(allRows))

	return
}
func (r *Repo) GetAllPrices() (prices map[string]*PairPrice, err error) {
	pairs := []*Pair{}
	pairs, err = r.GetAllPairs()
	if err != nil {
		err = fmt.Errorf("r.GetAllPairs: %w", err)
		return
	}
	prices = make(map[string]*PairPrice)

	for _, pair := range pairs {
		baseContract := pair.Token0.Contract
		baseSymbol := pair.Token0.Symbol

		quoteContract := pair.Token1.Contract
		quoteSymbol := pair.Token1.Symbol
		quotePrice := pair.Price1Last

		if baseContract != str_const.CTether || eos.Symbol(baseSymbol).String() != str_const.SUsdt {
			baseContract = pair.Token1.Contract
			baseSymbol = pair.Token1.Symbol

			quoteContract = pair.Token0.Contract
			quoteSymbol = pair.Token0.Symbol
			quotePrice = pair.Price0Last
		}
		// 还是没有
		if baseContract != str_const.CTether || eos.Symbol(baseSymbol).String() != str_const.SUsdt {
			continue
		}

		key := fmt.Sprintf("%v-%v", quoteContract, quoteSymbol)
		prices[key] = &PairPrice{Contract: quoteContract, Symbol: eos.Symbol(quoteSymbol).String(), PriceInUsdt: quotePrice.Value}
	}

	prices[fmt.Sprintf("%v-%v", str_const.CTether, str_const.SUsdt)] = &PairPrice{Contract: str_const.CTether, Symbol: str_const.SUsdt, PriceInUsdt: 1}
	return
}

type GetTableRowsResult struct {
	Resp   *eos.GetTableRowsResp
	Error  error
	Params eos.GetTableRowsRequest
	Url    string
}
