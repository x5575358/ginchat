package main

import (
	myrouter "ginchat/router"

	"ginchat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()

  	r := myrouter.AppRouter()
  	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}