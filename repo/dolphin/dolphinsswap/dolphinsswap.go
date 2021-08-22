package dolphinsswap

import (
	"eos-arb/repo/eos_api"

	"gorm.io/gorm"
)

const Code string = "dolphinsswap"

type Repo struct {
	Api eos_api.IEosApi
	Db  *gorm.DB
}
