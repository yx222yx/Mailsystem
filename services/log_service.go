package services

import (
	"fmt"
	"mailsys/database"
	"mailsys/models"
	"strings"
	"time"
)

func LogError(err error) {
	entry := models.LogEntry{
		Type:    "ERROR",
		Message: err.Error(),
		Date:    time.Now().Unix(),
	}
	database.InsertLogEntry(entry)
	SendErrorEmail(err)
}

func LogInfo(info string) {
	entry := models.LogEntry{
		Type:    "INFO",
		Message: info,
		Date:    time.Now().Unix(),
	}
	database.InsertLogEntry(entry)
}

func LogUnsent(email, reason string) {
	entry := models.UnsentEmail{
		Email:  email,
		Reason: reason,
	}
	database.InsertUnsentEmail(entry)
}

func RetryUnsentEmails() {
	unsentEmails, _ := database.GetUnsentEmails()
	for _, entry := range unsentEmails {
		email := entry.Email
		reason := entry.Reason

		if strings.Contains(reason, "birthday") {
			employee := models.Employee{Email: email}
			err := SendBirthdayEmail(employee)
			if err == nil {
				LogInfo(fmt.Sprintf("成功重发生日邮件至 %s", email))
				database.DeleteUnsentEmail(entry.ID)
			} else {
				LogError(fmt.Errorf("重发生日邮件至 %s 失败: %v", email, err))
			}
		} else if strings.Contains(reason, "anniversary") {
			employee := models.Employee{Email: email}
			err := SendWorkAnniversaryEmail(employee)
			if err == nil {
				LogInfo(fmt.Sprintf("成功重发周年纪念邮件至 %s", email))
				database.DeleteUnsentEmail(entry.ID)
			} else {
				LogError(fmt.Errorf("重发周年纪念邮件至 %s 失败: %v", email, err))
			}
		}
	}
}
