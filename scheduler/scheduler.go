package scheduler

import (
	"fmt"
	"mailsys/models"
	"mailsys/services"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/xuri/excelize/v2"
)

func readExcel(filePath string) ([]models.Employee, []string) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		errMsg := fmt.Sprintf("打开 Excel 文件出错: %v", err)
		services.LogError(fmt.Errorf(errMsg))
		return nil, []string{errMsg}
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		errMsg := fmt.Sprintf("读取 Excel 行出错: %v", err)
		services.LogError(fmt.Errorf(errMsg))
		return nil, []string{errMsg}
	}

	var employees []models.Employee
	var errorMsgs []string
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 5 {
			errMsg := fmt.Sprintf("第 %d 行数据不完整", i+1)
			services.LogError(fmt.Errorf(errMsg))
			errorMsgs = append(errorMsgs, errMsg)
			continue
		}

		joinDate, err := time.Parse("2006-01-02", row[2])
		if err != nil {
			errMsg := fmt.Sprintf("第 %d 行入职日期格式错误", i+1)
			services.LogError(fmt.Errorf(errMsg))
			errorMsgs = append(errorMsgs, errMsg)
			continue
		}

		birthDate, err := time.Parse("2006-01-02", row[3])
		if err != nil {
			errMsg := fmt.Sprintf("第 %d 行生日日期格式错误", i+1)
			services.LogError(fmt.Errorf(errMsg))
			errorMsgs = append(errorMsgs, errMsg)
			continue
		}

		if row[4] == "" {
			errMsg := fmt.Sprintf("第 %d 行缺少邮箱", i+1)
			services.LogError(fmt.Errorf(errMsg))
			errorMsgs = append(errorMsgs, errMsg)
			continue
		}

		employee := models.Employee{
			Name:      row[0],
			Category:  row[1],
			JoinDate:  joinDate,
			BirthDate: birthDate,
			Email:     row[4],
		}
		employees = append(employees, employee)
	}

	return employees, errorMsgs
}

func checkAndSendEmails(filePath string) {
	employees, errs := readExcel(filePath)
	if len(errs) > 0 {
		combinedErrMsg := strings.Join(errs, "\n")
		services.SendErrorEmail(fmt.Errorf("读取 Excel 文件过程中发生以下错误:\n%v", combinedErrMsg))
	}

	today := time.Now()
	for _, employee := range employees {
		if employee.BirthDate.Month() == today.Month() && employee.BirthDate.Day() == today.Day() {
			err := services.SendBirthdayEmail(employee)
			if err != nil {
				services.LogError(fmt.Errorf("发送生日邮件至 %s 失败: %v", employee.Name, err))
			}
		}

		if employee.JoinDate.Month() == today.Month() && employee.JoinDate.Day() == today.Day() {
			err := services.SendWorkAnniversaryEmail(employee)
			if err != nil {
				services.LogError(fmt.Errorf("发送入职周年邮件至 %s 失败: %v", employee.Name, err))
			}
		}
	}
	services.LogInfo("每日邮件检查和发送已完成")
}

func ScheduleDailyNotifications(filePath string) {
	c := cron.New()

	// 每天早上8点检查并发送邮件
	c.AddFunc("6 11 * * *", func() {
		checkAndSendEmails(filePath)
	})

	// 每小时重试未发送的邮件
	for hour := 8; hour <= 13; hour++ {
		cronExpr := fmt.Sprintf("0 %d * * *", hour)
		c.AddFunc(cronExpr, func() {
			services.RetryUnsentEmails()
		})
	}

	c.Start()
}
