package defisswapcnt

import (
	"fmt"

	"github.com/cihub/seelog"
	"github.com/panyanyany/eos-go"
)

func (r *Repo) GetMultiSwapAction(params PushMultiSwapActionInput) (action *eos.Action, err error) {
	inputContract := params.Pairs[0].Contract0
	inputSymbol := params.Pairs[0].Sym0
	if inputSymbol.String() != params.In.Symbol.String() {
		inputContract = params.Pairs[0].Contract1
		inputSymbol = params.Pairs[0].Sym1
	}
	if inputSymbol.String() != params.In.Symbol.String() {
		err = fmt.Errorf("not found symbol: %v", params.In.Symbol.String())
		//errs = append(errs, err)
		return
	}

	//memo := fmt.Sprintf("swap,%d,%d", params.Out.Amount, params.Pair.ID)
	memo := r.GetSwapMemo(params.Pairs, params.Out)

	actionData := eos.NewActionData(TransferAction{
		From:     r.Api.GetActor(),
		To:       eos.AccountName(Code),
		Quantity: params.In,
		Memo:     memo,
	})
	action = &eos.Action{
		Account: eos.AccountName(inputContract),
		Name:    "transfer",
		Authorization: []eos.PermissionLevel{
			{Actor: r.Api.GetActor(), Permission: eos.PermissionName("active")},
		},
		ActionData: actionData,
	}
	return
}
func (r *Repo) PushMultiSwapAction(params PushMultiSwapActionInput) (resp *eos.PushTransactionFullResp, err error) {
	var action *eos.Action
	action, err = r.GetMultiSwapAction(params)
	if err != nil {
		err = fmt.Errorf("r.GetMultiSwapAction: %w", err)
		return
	}

	seelog.Debugf("pushActions: %v", action.ActionData)
	resp, err = r.Api.PushActions([]*eos.Action{action})
	return
}

type PushMultiSwapActionInput struct {
	Pairs []*Pair
	In    eos.Asset
	Out   eos.Asset
}

type TransferAction struct {
	From     eos.AccountName `json:"from"`
	To       eos.AccountName `json:"to"`
	Quantity eos.Asset       `json:"quantity"`
	Memo     string          `json:"memo"`
}
