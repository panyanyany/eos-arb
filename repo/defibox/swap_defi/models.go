package swap_defi

import (
	"eos-arb/util/json_util"

	"github.com/panyanyany/eos-go"
)

type Pair struct {
	ID     int `json:"id"`
	Token0 struct {
		Contract string     `json:"contract"`
		Symbol   eos.Symbol `json:"symbol"`
	} `json:"token0"`
	Token1 struct {
		Contract string     `json:"contract"`
		Symbol   eos.Symbol `json:"symbol"`
	} `json:"token1"`
	Reserve0             eos.Asset              `json:"reserve0"`
	Reserve1             eos.Asset              `json:"reserve1"`
	LiquidityToken       json_util.StrOrUint64  `json:"liquidity_token"`
	Price0Last           json_util.StrOrFloat64 `json:"price0_last"`
	Price1Last           json_util.StrOrFloat64 `json:"price1_last"`
	Price0CumulativeLast json_util.StrOrUint64  `json:"price0_cumulative_last"`
	Price1CumulativeLast json_util.StrOrUint64  `json:"price1_cumulative_last"`
	BlockTimeLast        string                 `json:"block_time_last"`
}

type PairPrice struct {
	Contract    string
	Symbol      string
	PriceInUsdt float64
}
