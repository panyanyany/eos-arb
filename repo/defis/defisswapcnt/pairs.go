package defisswapcnt

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"eos-arb/model"

	"github.com/cihub/seelog"
	"github.com/panyanyany/eos-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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
			Table:      "markets",
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
			_, found := exists[row.Mid]
			if found {
				continue
			}
			exists[row.Mid] = true
			allRows = append(allRows, row)
		}
		if !resp.More {
			break
		}
		lowerBound = strconv.Itoa(rows[len(rows)-1].Mid)
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
			ids = append(ids, row.Mid)
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
				Table:      "markets",
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
			_, found := exists[row.Mid]
			if found {
				continue
			}
			exists[row.Mid] = true
			allRows = append(allRows, row)
		}
		//seelog.Infof("url=%v, lowerBound=%v, len=%v", res.Url, res.Params.LowerBound, len(rows))
	}

	//seelog.Infof("get all pairs: %v", len(allRows))

	return
}

type GetTableRowsResult struct {
	Resp   *eos.GetTableRowsResp
	Error  error
	Params eos.GetTableRowsRequest
	Url    string
}
