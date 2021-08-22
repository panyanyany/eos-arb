package defisswapcnt

import (
	"eos-arb/repo/eos_api"

	"gorm.io/gorm"
)

const Code string = "defisswapcnt"

type Repo struct {
	Api eos_api.IEosApi
	Db  *gorm.DB
}
