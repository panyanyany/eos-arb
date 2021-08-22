package defisswapcnt

import (
	"eos-arb/util/json_util"

	"github.com/panyanyany/eos-go"
)

type Pair struct {
	Mid                  int                    `json:"mid"`
	Contract0            string                 `json:"contract0"`
	Contract1            string                 `json:"contract1"`
	Sym0                 eos.Symbol             `json:"sym0"`
	Sym1                 eos.Symbol             `json:"sym1"`
	Reserve0             eos.Asset              `json:"reserve0"`
	Reserve1             eos.Asset              `json:"reserve1"`
	LiquidityToken       json_util.StrOrUint64  `json:"liquidity_token"`
	Price0Last           json_util.StrOrFloat64 `json:"price0_last"`
	Price1Last           json_util.StrOrFloat64 `json:"price1_last"`
	Price0CumulativeLast json_util.StrOrUint64  `json:"price0_cumulative_last"`
	Price1CumulativeLast json_util.StrOrUint64  `json:"price1_cumulative_last"`
	LastUpdate           string                 `json:"last_update"`
}
