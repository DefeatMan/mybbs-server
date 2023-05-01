package main

import (
	"kome/mybbs-server/dao"
	"kome/mybbs-server/router"
	"kome/mybbs-server/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	utils.InitConfig()
	dao.Init()
	r := gin.Default()
	r = router.InitRoute(r)
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	} else {
		panic(r.Run())
	}
}
