package models

import (
	// "fmt"
	// "ginchat/utils"
	//"time"

	"fmt"
	"ginchat/utils"

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	OwnerId  uint   //谁的关系
	TargetId uint   //对应的谁
	Type     string //对应类型 ，1好友 2群 3
	Dsc      string
}

func (table *Contact) TableName() string {
	return "contact"
}

func SearchFriend(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id= ? and type=1 ", userId).First(&contacts)
	for _, v := range contacts {
		fmt.Println(v)
		objIds = append(objIds, uint64(v.TargetId))
	}
	users := make([]UserBasic, 0)
	utils.DB.Where("id=?", objIds).Find(&users)
	return users
}
