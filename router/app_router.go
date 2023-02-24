package router

import (
	"ginchat/docs"
	"ginchat/service"

	"github.com/gin-gonic/gin"

	//docs "github.com/go-project-name/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func AppRouter() *gin.Engine {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/user/list", service.GetUserList)
	r.POST("/user/createUser", service.CreateUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/deleteUser", service.DeleteUser)
	r.POST("/user/findUserByNameAndPwd", service.FindUserByNameAndPwd)
	//注册
	r.GET("/myregister", service.MyRegister)
	r.GET("/mychat", service.Mychat)
	r.GET("/chat", service.Chat)
	r.POST("/searchfriend", service.SearchFriend)

	//发送消息
	r.GET("/user/SendMsg", service.SendMsg)
	r.GET("/user/SendUserMsg", service.SendUserMsg)

	//静态资源
	r.Static("/asset", "asset/")
	r.LoadHTMLGlob("views/**/*")

	return r
}
