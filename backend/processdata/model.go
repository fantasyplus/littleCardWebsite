package processdata

import "gorm.io/gorm"

type PersonInfo struct {
	gorm.Model
    // PersonID int `gorm:"primary_key"`
    CN       string `gorm:"type:varchar(255); not null" json:"cn" binding:"required"`
    QQ       string `gorm:"type:varchar(255); not null" json:"qq" binding:"required"`
}

type CardInfo struct {
	gorm.Model
    CardID         string `gorm:"primary_key"`
    CardName       string `gorm:"type:varchar(255); not null" json:"card_name" binding:"required"`
    CardCharacter  string `gorm:"type:varchar(255); not null" json:"card_character" binding:"required"`
    CardType       string `gorm:"type:varchar(255); not null" json:"card_type" binding:"required"`
    CardCondition  string `gorm:"type:varchar(255); not null" json:"card_condition" binding:"required"`
    Other          string `gorm:"type:varchar(1024); not null" json:"other" binding:"required"`
}

type CardIndex struct {
	gorm.Model
    PersonID int `gorm:"type:int; not null" json:"person_id" binding:"required"`
    CardIDs  string `gorm:"type:varchar(8192); not null" json:"card_ids" binding:"required"`
}

type CardNo struct{
    gorm.Model
    PersonID int `gorm:"type:int; not null" json:"person_id" binding:"required"`
    CardName string `gorm:"type:varchar(255); not null" json:"card_name" binding:"required"`
    CardNum float64 `gorm:"type:float; not null" json:"card_num" binding:"required"`
    Status string `gorm:"type:varchar(255); not null" json:"status" binding:"required"`
}