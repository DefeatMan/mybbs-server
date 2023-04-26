package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

var ShowNum = 20

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./conf/")
	err := viper.ReadInConfig()
	if err != nil {
        panic(fmt.Errorf("[Config Error]: failed to load config, %w", err))
	}
    InitJwt()
    vshow_num := viper.Get("shownum")
    if vshow_num != nil {
        ShowNum = int(vshow_num.(float64))
    }
    fmt.Println("ShowNum: ",ShowNum)
}
