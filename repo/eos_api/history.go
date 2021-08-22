package eos_api

import (
	"crypto/md5"
	"encoding/json"
	"fmt"

	"github.com/cihub/seelog"
	"github.com/go-redis/redis/v8"
	"github.com/parnurzeal/gorequest"
)

func (r *EosApi) GetActions(params GetActionsInput) (out *ActionOutput, err error) {
	var bs []byte
	key := "GetActions:"
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
		baseUrl := r.GetLock(r.HistoryApiList)
		fullUrl := baseUrl + "/v2/history/get_actions?" + params.ToValues().Encode()
		seelog.Debugf("GetActions: %v", fullUrl)
		_, bs, errs := gorequest.New().Get(fullUrl).EndBytes()
		r.ReleaseLock(baseUrl)
		if len(errs) != 0 {
			err = fmt.Errorf("request failed: %v", errs)
			return
		}
		if r.EnableCache {
			err = r.Rdb.Set(ctx, key, string(bs), r.RedisTimeout).Err()
			if err != nil {
				err = fmt.Errorf("r.Rdb.Set: %w", err)
				return
			}
		}
		err = json.Unmarshal(bs, out)
		if err != nil {
			err = fmt.Errorf("json.Unmarshal response body: %w", err)
			return
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
