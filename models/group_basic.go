package models

import (
	// "fmt"
	// "ginchat/utils"
	//"time"

	"gorm.io/gorm"
)

type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerId uint
	Icon    uint
	Dsc     string
	Type    string

	
}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}
