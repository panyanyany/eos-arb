package arbitrage

import (
	"eos-arb/repo/arbitrage/dex_adapter"
	"eos-arb/repo/arbitrage/printer"
	"eos-arb/repo/eos_api"
)

type Repo struct {
	Dexes   []dex_adapter.IDex
	Api     eos_api.IEosApi
	Printer *printer.PathPrinter
}
