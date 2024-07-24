package services

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"gopkg.in/yaml.v3"
	"html/template"
	"log"
	"mailsys/database"
	"mailsys/models"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/gomail.v2"
)

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type AdminConfig struct {
	Email string `yaml:"email"`
}

type Config struct {
	SMTP  SMTPConfig  `yaml:"smtp"`
	Admin AdminConfig `yaml:"admin"`
}

var Cfg Config

func LoadConfig() {
	file, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("无法打开配置文件: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&Cfg); err != nil {
		log.Fatalf("无法解析配置文件: %v", err)
	}
}

func readTemplate(filePath string) (string, error) {
	content, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", Cfg.SMTP.Username, Cfg.SMTP.Name) // 使用配置文件中的值
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(Cfg.SMTP.Host, Cfg.SMTP.Port, Cfg.SMTP.Username, Cfg.SMTP.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("邮件发送失败: %v", err)
	}

	return nil
}

func SendErrorEmail(err error) {
	subject := "错误通知"
	body := fmt.Sprintf("发生错误: %v", err)
	if err := SendEmail(Cfg.Admin.Email, subject, body); err != nil { // 使用配置文件中的管理员邮箱
		logEntry := models.LogEntry{
			Type:    "ERROR",
			Message: fmt.Sprintf("错误邮件发送失败: %v", err),
			Date:    time.Now().Unix(),
		}
		database.InsertLogEntry(logEntry)
	}
}

func sendTemplatedEmail(employee models.Employee, templatePath, subject string, data interface{}) error {
	templateContent, err := readTemplate(templatePath)
	if err != nil {
		LogError(fmt.Errorf("读取模板失败: %v", err))
		LogUnsent(employee.Email, "读取模板失败")
		return err
	}

	tmpl, err := template.New("emailTemplate").Parse(templateContent)
	if err != nil {
		LogError(fmt.Errorf("解析模板失败: %v", err))
		LogUnsent(employee.Email, "解析模板失败")
		return err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		LogError(fmt.Errorf("执行模板失败: %v", err))
		LogUnsent(employee.Email, "执行模板失败")
		return err
	}

	err = SendEmail(employee.Email, subject, body.String())
	if err == nil {
		LogInfo(fmt.Sprintf("邮件已发送给 %s", employee.Name))
	} else {
		LogUnsent(employee.Email, "邮件发送失败")
	}
	return err
}

func SendBirthdayEmail(employee models.Employee) error {
	subject := "生日快乐"
	data := struct {
		Name      string
		BirthDate string
	}{
		Name:      employee.Name,
		BirthDate: employee.BirthDate.Format("2006-01-02"),
	}
	templatePath := "./templates/birthday_template.html"
	err := sendTemplatedEmail(employee, templatePath, subject, data)
	if err == nil {
		LogInfo(fmt.Sprintf("生日邮件已成功发送给 %s (%s)", employee.Name, employee.Email))
	}
	return err
}

func SendWorkAnniversaryEmail(employee models.Employee) error {
	subject := "入职周年快乐"
	data := struct {
		Name  string
		Years int
	}{
		Name:  employee.Name,
		Years: time.Now().Year() - employee.JoinDate.Year(),
	}
	templatePath := "./templates/anniversary_template.html"
	err := sendTemplatedEmail(employee, templatePath, subject, data)
	if err == nil {
		LogInfo(fmt.Sprintf("周年纪念邮件已成功发送给 %s (%s)", employee.Name, employee.Email))
	}
	return err
}
