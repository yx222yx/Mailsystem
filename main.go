package main

import (
	"mailsys/database"
	"mailsys/scheduler"
	"mailsys/services"
)

func main() {
	// 加载配置文件
	services.LoadConfig()

	// 初始化数据库
	database.Init()

	// 使用相对路径
	filePath := "./employees.xlsx"
	scheduler.ScheduleDailyNotifications(filePath)

	select {} // 阻塞主程序以保持调度器运行
}
