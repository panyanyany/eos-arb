package dex_adapter

import (
	"eos-arb/repo/arbitrage/models"
	"eos-arb/repo/eos_api"

	"github.com/panyanyany/eos-go"
	"gorm.io/gorm"
)

type IDex interface {
	GetAllPairs() ([]*models.Pair, error)
	Init(db *gorm.DB) error
	GetSwapCmd(GetSwapCmdInput) (string, error)
	GetMultiSwapCmd(GetMultiSwapCmdInput) (string, error)
	PushMultiSwapAction(params PushMultiSwapActionInput) (resp *eos.PushTransactionFullResp, err error)
	GetMultiSwapAction(params PushMultiSwapActionInput) (action *eos.Action, err error)
	GetName() string
}

type BaseAdapter struct {
	Api eos_api.IEosApi
}

type GetSwapCmdInput struct {
	Pair *models.Pair
	From eos.Name
	In   eos.Asset
	Out  eos.Asset
}
type GetMultiSwapCmdInput struct {
	Pairs []*models.Pair
	From  eos.Name
	In    eos.Asset
	Out   eos.Asset
}
type PushMultiSwapActionInput struct {
	Pairs []*models.Pair
	In    eos.Asset
	Out   eos.Asset
}
