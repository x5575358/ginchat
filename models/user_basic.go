package models

import (
	"fmt"
	"ginchat/utils"
	"time"

	//"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string `vaild:"matches(^1[3-9{1}\\d{9}$])"`
	Email         string `vaild:"email"`
	Identity      string
	ClientIp      string
	ClientPort    string
	Salt          string
	LoginTime     string
	HeartbeatTime string
	LoginOutTime  string `gorm:"column:login_out_time" json:"login_out_time"`
	IsLogout      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func FindUserByNameAndPwd(name string,password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name=? and pass_word=?", name,password).First(&user)	
	//token加密
	str:=fmt.Sprintf("%d",time.Now().Unix())
	temp:= utils.Md5Encode(str)
	utils.DB.Model(&user).Where("id=?",user.ID).Update("identity",temp)
	return user

}


func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name=?", name).First(&user)
	return user

}

func FindUserByPhone(phone string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("phone=?", phone).First(&user)
}

func FindUserByEmail(email string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("email=?", email).First(&user)
}

func CreatUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}

func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{Name: user.Name, PassWord: user.PassWord, Phone: user.Phone, Email: user.Email})
}
