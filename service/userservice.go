package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"math/rand"
	"time"

	//"internal/goversion"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"

	//"golang.org/x/net/websocket"
	"github.com/gorilla/websocket"
)

// GetUserList
// @Summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/list [get]
func GetUserList(c *gin.Context) {
	//data := make([]*models.UserBasic, 10)
	data := models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "find finished 1",
		"data":    data,
	})
}

// GetUserList
// @Summary 所有用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	//name := c.Query("name")
	name := c.Request.FormValue("name")
	//password := c.Query("password")
	password := c.Request.FormValue("password")
	fmt.Println(name, password)
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "user not exist",
			"data":    data,
		})
	}
	fmt.Println(user)
	flag := utils.VaildPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "password incorrect",
			"data":    data,
		})
		return
	}
	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "find finished 2",
		"data":    data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword  query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	// user.Name = c.Query("name")
	// password := c.Query("password")
	// repassword := c.Query("repassword")
	user.Name = c.PostForm("name")
	password := c.PostForm("password")
	repassword := c.PostForm("Identity")
	fmt.Println(user.Name, ">>>>>", password, repassword)

	salt := fmt.Sprintf("%06d", rand.Int31())

	data := models.FindUserByName(user.Name)
	//fmt.Println(user.Name,"······",data.Name)

	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "user is not null",
			"data":    data,
		})
		return
	}
	if data.Name != "" {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "user is register",
			"data":    data,
		})
		return
	}
	if password != repassword {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "2次密码不一致",
			"data":    data,
		})
		return
	}
	//user.PassWord = password
	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt
	fmt.Println(user.PassWord)

	models.CreatUser(user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "add ok",
		"data":    user,
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	fmt.Println("upodate:", user)
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "update no exec",
			"data":    user,
		})
	} else {
		models.UpdateUser(user)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "update ok",
			"data":    user,
		})

	}

}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    0, //0成功 -1失败
		"message": "delete ok",
		"data":    user,
	})
}

// 防跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c)

}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	msg, err := utils.Subscribe(c, utils.PublishKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("发送消息", msg)
	tm := time.Now().Format("2006-01-02 15:45:05")
	m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
	err = ws.WriteMessage(1, []byte(m))
	if err != nil {
		fmt.Println(err)
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

func SearchFriend(c *gin.Context) {
	id,_:=strconv.Atoi(c.Request.FormValue("userId"))
	users:=models.SearchFriend(uint(id))
	// c.JSON(200,gin.H{
	// 	"code":0,
	// 	"message":"查询好友成功",
	// 	"data":users,
	// })
	utils.RespOkList(c.Writer,users,len(users))

}

