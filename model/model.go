package model

import "gorm.io/gorm"

type PairIdList struct {
	gorm.Model

	Contract string `gorm:"type:varchar(255);unique"`
	IdList   string `gorm:""`
	Total    int
}
