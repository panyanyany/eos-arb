package eos_api

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/cihub/seelog"
	"github.com/go-redis/redis/v8"
	"github.com/panyanyany/eos-go"
)

type EosApi struct {
	Rdb            *redis.Client
	RedisTimeout   time.Duration
	EnableCache    bool
	FindLock       sync.Mutex
	LastRequest    time.Time
	PrivateKey     string
	QueryApiList   []string
	PostApiList    []string
	HistoryApiList []string
	ApiLastRequest map[string]time.Time
	Apis           map[string]*eos.API
	UrlLocks       map[string]*sync.Mutex
	Actor          eos.AccountName
}

type IEosApi interface {
	GetTableRows(params eos.GetTableRowsRequest) (out *eos.GetTableRowsResp, baseUrl string, err error)
	PushActions(actions []*eos.Action) (response *eos.PushTransactionFullResp, err error)
	GetActor() (actor eos.AccountName)
	GetActions(params GetActionsInput) (out *ActionOutput, err error)
}

func NewEosApi(rdb *redis.Client, redisTimeout time.Duration) (r *EosApi) {
	r = new(EosApi)
	r.RedisTimeout = redisTimeout
	r.Rdb = rdb
	r.FindLock = sync.Mutex{}
	r.EnableCache = rdb != nil && redisTimeout.Seconds() != 0
	// https://gist.github.com/akme/89a4e596587cb605b530bd825994a0db
	r.QueryApiList = []string{
		// 测试 http://api1.eosasia.one/v1/chain/get_info
		"https://eos.newdex.one",
		"http://eos.greymass.com",
		"http://api.eossweden.org",
		"https://api.eoslaomao.com",

		"https://api1.eosasia.one",
		"http://api-mainnet.starteos.io",
		"http://api1.eosasia.one",
		"http://fn001.eossv.org",
		"http://mainnet.eosamsterdam.net",
		"http://eosbp-0.atticlab.net",
		"http://eosbp-1.atticlab.net",
		"http://eos.eoscafeblock.com",
		"http://api.eosn.io",
		"http://node2.eosphere.io",
		"http://node1.eosphere.io",
		"http://seed01.eosusa.news",
		"http://seed02.eosusa.news",
		"http://api.eosrio.io",
		"http://api.eostitan.com",

		"http://api.eoseoul.io",    //易错
		"http://bp.cryptolions.io", //易错
		//"http://mainnet.eosio.sg",
		//"https://mainnet.coscannon.io":    "https://mainnet.coscannon.io", 无效
		//"https://eos.rrdy.com":            "https://eos.rrdy.com", 无效
	}
	r.PostApiList = []string{
		"https://eospush.tokenpocket.pro",
	}
	r.HistoryApiList = []string{
		"https://eos.hyperion.eosrio.io",
	}
	r.ApiLastRequest = map[string]time.Time{}
	r.Apis = map[string]*eos.API{}
	r.UrlLocks = map[string]*sync.Mutex{}

	allUrlList := []string{}
	allUrlListExists := make(map[string]bool)
	for _, url := range r.QueryApiList {
		_, found := allUrlListExists[url]
		if found {
			continue
		}
		allUrlListExists[url] = true
		allUrlList = append(allUrlList, url)
	}
	for _, url := range r.PostApiList {
		_, found := allUrlListExists[url]
		if found {
			continue
		}
		allUrlListExists[url] = true
		allUrlList = append(allUrlList, url)
	}
	for _, url := range r.HistoryApiList {
		_, found := allUrlListExists[url]
		if found {
			continue
		}
		allUrlListExists[url] = true
		allUrlList = append(allUrlList, url)
	}

	for _, url := range allUrlList {
		r.ApiLastRequest[url] = time.Time{}
		api := eos.New(url)
		api.HttpClient.Transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 5 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       10 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives:     true, // default behavior, because of `nodeos`'s lack of support for Keep alives.
		}
		r.Apis[url] = api
		r.UrlLocks[url] = &sync.Mutex{}
	}
	return
}

var ctx = context.Background()

func (r *EosApi) GetActor() (actor eos.AccountName) {
	if r.Actor.String() == "" {
		panic("unset Actor")
	}
	return r.Actor
}
func (r *EosApi) GetLock(apiList []string) string {
	r.FindLock.Lock()
	var foundUrl string
	interval := int64(math.MaxInt64)

	for _, url := range apiList {
		if r.ApiLastRequest[url].IsZero() {
			foundUrl = url
			interval = 0
			break
		}
		now := time.Now()
		diff := 15 - int64(now.Sub(r.ApiLastRequest[url]).Seconds())
		if diff < interval {
			interval = diff
			foundUrl = url
		}
		//seelog.Infof("diff=%v, interval=%v, foundUrl=%v, url=%v, last=%v, now=%v",
		//	diff, interval, foundUrl,
		//	url, r.ApiLastRequest[url], now,
		//)
	}
	r.ApiLastRequest[foundUrl] = time.Now()
	r.FindLock.Unlock()

	r.UrlLocks[foundUrl].Lock()
	if interval > 0 {
		seelog.Debugf("sleep: %v", interval)
		time.Sleep(time.Duration(interval) * time.Second)
	}
	return foundUrl
}
func (r *EosApi) ReleaseLock(url string) {
	//r.LastRequest = time.Now()
	//r.ApiLastRequest[url] = time.Now()
	//r.FindLock.Unlock()
	r.UrlLocks[url].Unlock()
}
func (r *EosApi) PushActions(actions []*eos.Action) (response *eos.PushTransactionFullResp, err error) {
	keyBag := &eos.KeyBag{}
	err = keyBag.ImportPrivateKey(ctx, r.PrivateKey)
	if err != nil {
		err = fmt.Errorf("import private key: %w", err)
		return
	}

	//baseUrl := r.GetLock(r.PostApiList)
	baseUrl := r.PostApiList[0]

	api := r.Apis[baseUrl]
	api.SetSigner(keyBag)

	txOpts := &eos.TxOptions{}
	if err = txOpts.FillFromChain(ctx, api); err != nil {
		err = fmt.Errorf("filling tx opts: %w", err)
		return
	}

	tx := eos.NewTransaction(actions, txOpts)
	signedTx, packedTx, err := api.SignTransaction(ctx, tx, txOpts.ChainID, eos.CompressionNone)
	if err != nil {
		err = fmt.Errorf("sign transaction: %w", err)
		return
	}

	content, err := json.MarshalIndent(signedTx, "", "  ")
	if err != nil {
		err = fmt.Errorf("json marshalling transaction: %w", err)
		return
	}

	seelog.Debug(string(content))
	for _, action := range actions {
		seelog.Debugf("actionData: %#v", action.ActionData.Data)
	}

	response, err = api.PushTransaction(ctx, packedTx)
	if err != nil {
		err = fmt.Errorf("push transaction: %w", err)
		return
	}
	//r.ReleaseLock(baseUrl)

	seelog.Debugf("Transaction [%s] submitted to the network succesfully.\n", hex.EncodeToString(response.Processed.ID))
	return
}

func (r *EosApi) GetTableRows(params eos.GetTableRowsRequest) (out *eos.GetTableRowsResp, baseUrl string, err error) {
	var bs []byte
	key := "GetTableRows:"
	bs, err = json.Marshal(params)
	if err != nil {
		seelog.Debugf("params: %#v", params)
		err = fmt.Errorf("json.Marshal params for key: %w", err)
		return
	}

	key += fmt.Sprintf("%x", md5.Sum(bs))

	var str string

	if r.EnableCache {
		str, err = r.Rdb.Get(ctx, key).Result()
		if err != redis.Nil && err != nil {
			err = fmt.Errorf("r.Rdb.Get: %v, key=%v", err, key)
			return
		}
		if err == redis.Nil {
			//seelog.Debugf("params: %#v", params)
		}
		//seelog.Debugf("r.Rdb.Get key=%v, err=%v", key, err)
	}

	if str == "" {
		baseUrl = r.GetLock(r.QueryApiList)
		api := r.Apis[baseUrl]
		seelog.Debugf("GetTableRows, url=%v", api.BaseURL)
		tsStart := time.Now()
		ctx2, _ := context.WithTimeout(ctx, time.Second*10)
		out, err = api.GetTableRows(ctx2, params)
		r.ReleaseLock(baseUrl)
		if err != nil {
			seelog.Debugf("params: %#v", params)
			err = fmt.Errorf("r.Api.GetTableRows: %w, time=%v, url=%v", err, time.Now().Sub(tsStart).Seconds(), baseUrl)
			//seelog.Error(err)
			return
		}
		bs, err = json.Marshal(out)
		if err != nil {
			err = fmt.Errorf("json.Marshal: %w", err)
			return
		}
		if r.EnableCache {
			err = r.Rdb.Set(ctx, key, string(bs), r.RedisTimeout).Err()
			if err != nil {
				err = fmt.Errorf("r.Rdb.Set: %w", err)
				return
			}
		}
	} else {
		err = json.Unmarshal([]byte(str), &out)
		if err != nil {
			seelog.Debugf("str: %v", str)
			err = fmt.Errorf("json.Unmarshal: %w", err)
			return
		}
	}
	return
}
