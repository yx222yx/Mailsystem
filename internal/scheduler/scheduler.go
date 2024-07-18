package scheduler

import (
	"bufio"
	"fmt"
	"log"
	"mailsys/internal/services"
	"mailsys/internal/utils"
	"os"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

func checkAndSendEmails(filePath string) {
	// 读取员工数据
	employees, err := utils.ReadExcel(filePath)
	if err != nil {
		services.LogError(fmt.Errorf("Failed to read Excel file: %v", err))
		log.Fatalf("Failed to read Excel file: %v", err)
	}

	today := time.Now()
	for _, employee := range employees {
		// 检查生日
		if employee.BirthDate.Month() == today.Month() && employee.BirthDate.Day() == today.Day() {
			err := services.SendBirthdayEmail(employee)
			if err != nil {
				services.LogError(fmt.Errorf("Failed to send birthday email to %s: %v", employee.Name, err))
				log.Printf("Failed to send birthday email to %s: %v", employee.Name, err)
			}
		}

		// 检查入职周年
		if employee.JoinDate.Month() == today.Month() && employee.JoinDate.Day() == today.Day() {
			err := services.SendWorkAnniversaryEmail(employee)
			if err != nil {
				services.LogError(fmt.Errorf("Failed to send work anniversary email to %s: %v", employee.Name, err))
				log.Printf("Failed to send work anniversary email to %s: %v", employee.Name, err)
			}
		}
	}
	services.LogInfo("Daily email check and send completed")
}

func checkErrorLogFile(logFilePath string) {
	file, err := os.Open(logFilePath)
	if err != nil {
		services.LogError(fmt.Errorf("Failed to open log file: %v", err))
		log.Printf("Failed to open log file: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	today := time.Now().Format("2006-01-02")
	newErrors := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, today) && strings.Contains(line, "ERROR") {
			newErrors = append(newErrors, line)
		}
	}

	if err := scanner.Err(); err != nil {
		services.LogError(fmt.Errorf("Error reading log file: %v", err))
		log.Printf("Error reading log file: %v", err)
	}

	if len(newErrors) > 0 {
		errorReport := fmt.Sprintf("New errors found in log file on %s:\n%s", today, strings.Join(newErrors, "\n"))
		services.SendErrorEmail(fmt.Errorf(errorReport))
	}
}

func ScheduleDailyNotifications(filePath string) {
	c := cron.New()

	// 每天早上8点发送邮件
	c.AddFunc("43 13 * * *", func() {
		checkAndSendEmails(filePath)
	})

	// 每天8-17点每个整点检查日志文件并重试发送未发送的邮件
	for hour := 8; hour <= 17; hour++ {
		cronExpr := fmt.Sprintf("44 %d * * *", hour)
		c.AddFunc(cronExpr, func() {
			checkErrorLogFile("mailsys_log/mailsys_err.log")
			services.RetryUnsentEmails()
		})
	}

	c.Start()
}
