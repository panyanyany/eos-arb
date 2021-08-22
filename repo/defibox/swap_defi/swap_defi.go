package swap_defi

import (
	"eos-arb/repo/eos_api"

	"gorm.io/gorm"
)

const Code string = "swap.defi"

type Repo struct {
	Api eos_api.IEosApi
	Db  *gorm.DB
}
