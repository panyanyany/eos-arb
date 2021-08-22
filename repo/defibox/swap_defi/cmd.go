package swap_defi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/panyanyany/eos-go"
)

func (r *Repo) GetSwapCmd(params GetSwapCmdInput) (cmd string, err error) {
	inputContract := params.Pair.Token0.Contract
	inputSymbol := params.Pair.Token0.Symbol
	if inputSymbol.String() != params.In.Symbol.String() {
		inputContract = params.Pair.Token1.Contract
		inputSymbol = params.Pair.Token1.Symbol
	}
	if inputSymbol.String() != params.In.Symbol.String() {
		err = fmt.Errorf("not found symbol: %v", params.In.Symbol.String())
		return
	}

	//memo := fmt.Sprintf("swap,%d,%d", params.Out.Amount, params.Pair.ID)
	memo := r.GetSwapMemo([]*Pair{params.Pair}, params.Out)
	data := fmt.Sprintf(`{"from":"%s", "to":"%s", "quantity":"%s", "memo": "%s"}`, params.From, Code, params.In, memo)
	cmd = fmt.Sprintf("cleos -u https://eospush.tokenpocket.pro push action %v transfer '%v' -p %v",
		inputContract,
		data,
		params.From,
	)

	return
}

func (r *Repo) GetSwapMemo(pairs []*Pair, out eos.Asset) (memo string) {
	sIds := []string{}
	for _, pair := range pairs {
		sIds = append(sIds, strconv.Itoa(pair.ID))
	}
	memo = fmt.Sprintf("swap,%d,%s", out.Amount, strings.Join(sIds, "-"))
	return
}

func (r *Repo) GetMultiSwapCmd(params GetMultiSwapCmdInput) (cmd string, err error) {
	inputContract := params.Pairs[0].Token0.Contract
	inputSymbol := params.Pairs[0].Token0.Symbol
	if inputSymbol.String() != params.In.Symbol.String() {
		inputContract = params.Pairs[0].Token1.Contract
		inputSymbol = params.Pairs[0].Token1.Symbol
	}
	if inputSymbol.String() != params.In.Symbol.String() {
		err = fmt.Errorf("not found symbol: %v", params.In.Symbol.String())
		return
	}

	//memo := fmt.Sprintf("swap,%d,%d", params.Out.Amount, params.Pair.ID)
	memo := r.GetSwapMemo(params.Pairs, params.Out)
	data := fmt.Sprintf(`{"from":"%s", "to":"%s", "quantity":"%s", "memo": "%s"}`, params.From, Code, params.In, memo)
	cmd = fmt.Sprintf("cleos -u https://eospush.tokenpocket.pro push action %v transfer '%v' -p %v",
		inputContract,
		data,
		params.From,
	)

	return
}

type GetSwapCmdInput struct {
	Pair *Pair
	From eos.Name
	In   eos.Asset
	Out  eos.Asset
}

type GetMultiSwapCmdInput struct {
	Pairs []*Pair
	From  eos.Name
	In    eos.Asset
	Out   eos.Asset
}
