package main

import (
	"maboshijikaxitong/model"
	"maboshijikaxitong/utils"
	"maboshijikaxitong/view"
)

var (
	//confFilePath = "./conf/config.toml"
	confFilePath = "../conf/config.toml"
)

func main() {
	//println(utils.DifferenceOfDay("2023-01-13", time.Now().Format(time.DateOnly)))
	utils.ConfInit(confFilePath)
	utils.SetupZapLogger()

	//	初始化数据表
	model.DBInit()

	view.View()
}
