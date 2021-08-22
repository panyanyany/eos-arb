package dolphinsswap

import (
	"eos-arb/util/json_util"

	"github.com/panyanyany/eos-go"
)

type Pair struct {
	ID             int                   `json:"id"`
	Code           string                `json:"code"`
	SwapFee        int                   `json:"swap_fee"`
	TotalLptoken   json_util.StrOrUint64 `json:"total_lptoken"`
	CreateTime     int                   `json:"create_time"`
	LastUpdateTime int                   `json:"last_update_time"`
	Tokens         []struct {
		Weight int `json:"weight"`
		Symbol struct {
			Symbol   eos.Symbol `json:"symbol"`
			Contract string     `json:"contract"`
		} `json:"symbol"`
		Reserve eos.Asset `json:"reserve"`
	} `json:"tokens"`
}
