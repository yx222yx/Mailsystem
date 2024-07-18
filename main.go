package main

import (
	"mailsys/internal/scheduler"
)

func main() {
	filePath := "D:/go_projects/mailsys/dataexcel/employees.xlsx"
	scheduler.ScheduleDailyNotifications(filePath)
	select {} // 阻塞主程序以保持调度器运行
}
