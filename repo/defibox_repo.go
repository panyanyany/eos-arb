package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/parnurzeal/gorequest"
)

type DefiboxRepo struct {
	rdb     *redis.Client
	ctx     context.Context
	timeout time.Duration
}

func NewDefiboxRepo(rdb *redis.Client, ctx context.Context, timeout time.Duration) (r *DefiboxRepo) {
	r = new(DefiboxRepo)
	r.ctx = ctx
	r.rdb = rdb
	r.timeout = timeout
	return
}
func (r *DefiboxRepo) GetUserMarket(account string) (body string, err error) {
	key := "user_market:" + account
	body, err = r.rdb.Get(r.ctx, key).Result()
	if err != redis.Nil {
		return
	}
	err = nil
	defer func() {
		if err == nil {
			r.rdb.Set(r.ctx, key, body, r.timeout)
		}
	}()

	req := gorequest.New()
	req.Post("https://defibox.340wan.com/api/swap/account/capital")
	req.Header = map[string][]string{
		"authorization": []string{"Basic Og=="},
		"account":       []string{account},
		"user-agent":    []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"},
		"referer":       []string{"https://defibox.340wan.com/my/"},
	}
	req.Data = map[string]interface{}{
		"owner": account,
	}

	var resp gorequest.Response
	var errs []error
	var bytes []byte
	resp, bytes, errs = req.EndBytes()
	if len(errs) > 0 {
		err = fmt.Errorf("request failed: %v", errs)
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("request failed: status=%v", resp.StatusCode)
		return
	}
	body = string(bytes)
	return
}

func (r *DefiboxRepo) GetBalance(account string) (body string, err error) {
	key := "balance:" + account
	body, err = r.rdb.Get(r.ctx, key).Result()
	if err != redis.Nil {
		return
	}
	err = nil
	defer func() {
		if err == nil {
			r.rdb.Set(r.ctx, key, body, r.timeout)
		}
	}()

	req := gorequest.New()
	req.Get("https://defibox.340wan.com/api/swap/account/getBalances")
	req.Header = map[string][]string{
		"authorization": []string{"Basic Og=="},
		"account":       []string{account},
		"user-agent":    []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"},
		"referer":       []string{"https://defibox.340wan.com/my/"},
	}
	var resp gorequest.Response
	var errs []error
	var bytes []byte
	resp, bytes, errs = req.EndBytes()
	if len(errs) > 0 {
		err = fmt.Errorf("request failed: %v", errs)
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("request failed: status=%v", resp.StatusCode)
		return
	}
	body = string(bytes)
	return
}

func (r *DefiboxRepo) SwapGetMarket(symbol string) (body string, err error) {
	key := "swap_market:" + symbol
	body, err = r.rdb.Get(r.ctx, key).Result()
	if err != redis.Nil {
		return
	}
	err = nil
	defer func() {
		if err == nil {
			r.rdb.Set(r.ctx, key, body, r.timeout)
		}
	}()

	req := gorequest.New()
	req.Post("https://defibox.340wan.com/api/swap/getMarket")
	req.Header = map[string][]string{
		"authorization": []string{"Basic Og=="},
		"user-agent":    []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"},
		//"referer":       []string{"https://defibox.340wan.com/my/"},
	}
	//req.Send("{\"pairId\":\"1134\"}")
	req.Data = map[string]interface{}{
		"limit":       60,
		"type":        2,
		"symbol":      symbol,
		"orderColumn": "eos_balance",
		"isAsc":       false,
	}
	//req.Param("pairId", pairId)
	//seelog.Infof("%#v", req)
	var resp gorequest.Response
	var errs []error
	var bytes []byte
	resp, bytes, errs = req.EndBytes()
	//seelog.Info(string(bytes))
	if len(errs) > 0 {
		err = fmt.Errorf("request failed: %v", errs)
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("request failed: status=%v", resp.StatusCode)
		return
	}
	body = string(bytes)
	return
}
