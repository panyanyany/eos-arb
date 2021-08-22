package dolphinsswap

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/panyanyany/eos-go"
)

func (r *Repo) GetSwapCmd(params GetSwapCmdInput) (cmd string, err error) {
	inputContract := params.Pair.Tokens[0].Symbol.Contract
	inputSymbol := params.Pair.Tokens[0].Symbol.Symbol
	if inputSymbol.String() != params.In.Symbol.String() {
		inputContract = params.Pair.Tokens[1].Symbol.Contract
		inputSymbol = params.Pair.Tokens[1].Symbol.Symbol
	}
	if inputSymbol.String() != params.In.Symbol.String() {
		err = fmt.Errorf("not found symbol: %v", params.In.Symbol.String())
		return
	}

	//memo := fmt.Sprintf("swap,%d,%d", params.Out.Amount, params.Pair.Mid)
	memo := r.GetSwapMemo([]*Pair{params.Pair}, params.Out)
	data := fmt.Sprintf(`{"from":"%s", "to":"%s", "quantity":"%s", "memo": "%s"}`, params.From, Code, params.In, memo)
	cmd = fmt.Sprintf("cleos -u https://eospush.tokenpocket.pro push action %v transfer '%v' -p %v",
		inputContract,
		data,
		params.From,
	)

	return
}
func (r *Repo) GetMultiSwapCmd(params GetMultiSwapCmdInput) (cmd string, err error) {
	inputContract := params.Pairs[0].Tokens[0].Symbol.Contract
	inputSymbol := params.Pairs[0].Tokens[0].Symbol.Symbol
	if inputSymbol.String() != params.In.Symbol.String() {
		inputContract = params.Pairs[0].Tokens[1].Symbol.Contract
		inputSymbol = params.Pairs[0].Tokens[1].Symbol.Symbol
	}
	if inputSymbol.String() != params.In.Symbol.String() {
		err = fmt.Errorf("not found symbol: %v", params.In.Symbol.String())
		return
	}

	//memo := fmt.Sprintf("swap,%d,%d", params.Out.Amount, params.Pair.Mid)
	memo := r.GetSwapMemo(params.Pairs, params.Out)
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
	memo = fmt.Sprintf("swap:%s;min:%d", strings.Join(sIds, "-"), out.Amount)
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
