package models

import (
	"fmt"
	"math/big"

	"github.com/panyanyany/eos-go"
)

type Pair struct {
	Id       int
	Token0   *TokenDetail
	Token1   *TokenDetail
	Dex      interface{}
	RawPair  interface{}
	TheOther map[string]*TokenDetail
}

type TokenMeta struct {
	Contract string
	Symbol   eos.Symbol
	Key      string
}

type TokenDetail struct {
	TokenMeta
	Reserve   eos.Asset
	PriceLast float64
}

func (r *TokenMeta) GetKey() string {
	if r.Key == "" {
		r.Key = fmt.Sprintf("%v-%v", r.Contract, r.Symbol.String())
	}
	return r.Key
}

func (r *Pair) HasZero() bool {
	return r.Token1.Reserve.Amount == 0 || r.Token0.Reserve.Amount == 0 || r.Token0.PriceLast == 0 || r.Token1.PriceLast == 0
}
func (r *Pair) GetAmountOut(qty eos.Asset) (out eos.ExtendedAsset) {
	var err error
	_997 := new(big.Float).SetUint64(997)
	_1000 := new(big.Float).SetUint64(1000)

	inputAmount := qty.ToFloat()
	inputToken := r.Token0

	if inputToken.Symbol.String() != qty.Symbol.String() {
		inputToken = r.Token1
	}
	inputReserve := inputToken.Reserve.ToFloat()

	if inputToken.Symbol.String() != qty.Symbol.String() {
		panic(fmt.Sprintf("input symbol(%v) do not match any token.", qty.Symbol.String()))
	}

	outputToken := r.TheOther[inputToken.GetKey()]
	outputReserve := outputToken.Reserve.ToFloat()

	inputAmountWithFee := new(big.Float).Mul(inputAmount, _997)
	//fmt.Printf("inputAmountWithFee: %v\n", inputAmountWithFee)
	//fmt.Printf("outputReserve: %v\n", outputReserve)
	numerator := new(big.Float).Mul(inputAmountWithFee, outputReserve)
	//fmt.Printf("numerator: %v\n", numerator)
	denominator := new(big.Float).Add(new(big.Float).Mul(inputReserve, _1000), inputAmountWithFee)
	//fmt.Printf("inputReserve: %v\n", inputReserve)
	//fmt.Printf("denominator: %v\n", denominator)
	outputAmount := new(big.Float).Quo(numerator, denominator)
	//fmt.Printf("outputAmount: %v\n", outputAmount)

	var ass eos.Asset
	ass, err = eos.NewAssetFromFloat(outputAmount, outputToken.Symbol)
	if err != nil {
		err = fmt.Errorf("eos.NewAssetFromFloat: %w", err)
		panic(err)
	}
	out = eos.ExtendedAsset{
		Asset:    ass,
		Contract: eos.AccountName(outputToken.Contract),
	}
	return
}

type PathLayer struct {
	Level    int
	InAsset  eos.ExtendedAsset
	PathJobs []*PathJob
}

type PathJob struct {
	Level int
	In    *eos.ExtendedAsset
	Out   *eos.ExtendedAsset
	Pair  *Pair
	//SubPathJobs []*PathJob
	Parent *PathJob
}
