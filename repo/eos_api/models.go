package eos_api

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

type GetActionOutput struct {
	QueryTimeMs float64 `json:"query_time_ms"`
	Cached      bool    `json:"cached"`
	Lib         int     `json:"lib"`
	Total       struct {
		Value    int    `json:"value"`
		Relation string `json:"relation"`
	} `json:"total"`
	Actions []struct {
		Timestamp string `json:"timestamp"`
		BlockNum  int    `json:"block_num"`
		TrxID     string `json:"trx_id"`
		Act       struct {
			Account       string `json:"account"`
			Name          string `json:"name"`
			Authorization []struct {
				Actor      string `json:"actor"`
				Permission string `json:"permission"`
			} `json:"authorization"`
			Data struct {
				From     string  `json:"from"`
				To       string  `json:"to"`
				Amount   float64 `json:"amount"`
				Symbol   string  `json:"symbol"`
				Memo     string  `json:"memo"`
				Quantity string  `json:"quantity"`
			} `json:"data"`
		} `json:"act"`
		Notified             []string `json:"notified"`
		CPUUsageUs           int      `json:"cpu_usage_us,omitempty"`
		NetUsageWords        int      `json:"net_usage_words,omitempty"`
		GlobalSequence       int64    `json:"global_sequence"`
		Producer             string   `json:"producer"`
		ActionOrdinal        int      `json:"action_ordinal"`
		CreatorActionOrdinal int      `json:"creator_action_ordinal"`
		AccountRAMDeltas     []struct {
			Account string `json:"account"`
			Delta   int    `json:"delta"`
		} `json:"account_ram_deltas,omitempty"`
	} `json:"actions"`
}

type OfferActionOutputOfferLog struct {
	Assetid  string `json:"assetid"`
	Owner    string `json:"owner"`
	Category string `json:"category"`
	Idata    string `json:"idata"`
	Mdata    string `json:"mdata"`
	Color    string `json:"color"`
	RoleId   string `json:"role_id"`
	Level    string `json:"level"`
	Power0   string `json:"power0"`
	Power    string `json:"power"`
	Status   int    `json:"status"`
	Price    string `json:"price"`
	Offering int    `json:"offering"`
}

type OfferActionOutputBuyOfferLog struct {
	Assetid  string `json:"assetid"`
	PreOwner string `json:"pre_owner"`
	NowOwner string `json:"now_owner"`
	//Category string `json:"category"`
	//Idata    string `json:"idata"`
	//Mdata    string `json:"mdata"`
	Color  string `json:"color"`
	RoleId string `json:"role_id"`
	Level  string `json:"level"`
	Power0 string `json:"power0"`
	Power  string `json:"power"`
	Price  string `json:"price"`
}

type Act struct {
	Account       string `json:"account"`
	Name          string `json:"name"`
	Authorization []struct {
		Actor      string `json:"actor"`
		Permission string `json:"permission"`
	} `json:"authorization"`
	Data json.RawMessage `json:"data"`
}

type ActionOutput struct {
	QueryTimeMs float64 `json:"query_time_ms"`
	Cached      bool    `json:"cached"`
	Lib         int     `json:"lib"`
	Total       struct {
		Value    int    `json:"value"`
		Relation string `json:"relation"`
	} `json:"total"`
	Actions []struct {
		Timestamp            string   `json:"timestamp"`
		BlockNum             int      `json:"block_num"`
		TrxID                string   `json:"trx_id"`
		Act                  Act      `json:"act"`
		Notified             []string `json:"notified"`
		GlobalSequence       int64    `json:"global_sequence"`
		Producer             string   `json:"producer"`
		ActionOrdinal        int      `json:"action_ordinal"`
		CreatorActionOrdinal int      `json:"creator_action_ordinal"`
	} `json:"actions"`
}

type ActionOutputRaw struct {
	QueryTimeMs float64 `json:"query_time_ms"`
	Cached      bool    `json:"cached"`
	Lib         int     `json:"lib"`
	Total       struct {
		Value    int    `json:"value"`
		Relation string `json:"relation"`
	} `json:"total"`
	Actions []json.RawMessage `json:"actions"`
}

type GetActionsInput struct {
	Account string    `json:"account"`
	Filter  string    `json:"filter"`
	Skip    int       `json:"skip"`
	Limit   int       `json:"limit"`
	Sort    string    `json:"sort"`
	After   time.Time `json:"after"`
	Before  time.Time `json:"before"`
}

func (r *GetActionsInput) ToMap() (out map[string]string) {
	out = map[string]string{
		"account": r.Account,
		"filter":  r.Filter,
		"skip":    strconv.Itoa(r.Skip),
		"limit":   strconv.Itoa(r.Limit),
		"sort":    r.Sort,
	}
	if !r.After.IsZero() {
		out["after"] = r.After.UTC().Format(TimeLayoutQuery)
	}
	if !r.Before.IsZero() {
		out["before"] = r.Before.UTC().Format(TimeLayoutQuery)
	}
	return
}
func (r *GetActionsInput) ToValues() (out url.Values) {
	out = url.Values{}
	for key, val := range r.ToMap() {
		out.Set(key, val)
	}
	return
}
