package services

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"mailsys/internal/models"
	"net/smtp"
	"os"
	"strings"
)

const (
	ADMIN_EMAIL     = "2261576876@qq.com"
	SMTP_HOST       = "smtpdm.aliyun.com"
	SMTP_PORT       = 465
	SMTP_USERNAME   = "noreply@email.jxedc.com"
	SMTP_NAME       = "Clinflash 易迪希"
	SMTP_PASSWORD   = "ABcd123456"
	ERROR_LOG_FILE  = "mailsys_log/mailsys_err.log"
	NORMAL_LOG_FILE = "mailsys_log/mailsys_log.log"
	UNSENT_LOG_FILE = "mailsys_log/mailsys_unsent.log"
)

func init() {
	// 确保日志目录存在
	if _, err := os.Stat("mailsys_log"); os.IsNotExist(err) {
		os.Mkdir("mailsys_log", os.ModePerm)
	}
}

func LogError(err error) {
	// 打开错误日志文件
	f, ferr := os.OpenFile(ERROR_LOG_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if ferr != nil {
		log.Printf("Error opening error log file: %v", ferr)
		return
	}
	defer f.Close()

	// 写入错误日志文件
	logger := log.New(f, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println(err)

	// 发送报错邮件给管理员
	SendErrorEmail(err)
}

func LogInfo(info string) {
	// 打开正常日志文件
	f, ferr := os.OpenFile(NORMAL_LOG_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if ferr != nil {
		log.Printf("Error opening normal log file: %v", ferr)
		return
	}
	defer f.Close()

	// 写入正常日志文件
	logger := log.New(f, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println(info)
}

func LogUnsent(email, reason string) {
	// 打开未发送日志文件
	f, ferr := os.OpenFile(UNSENT_LOG_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if ferr != nil {
		log.Printf("Error opening unsent log file: %v", ferr)
		return
	}
	defer f.Close()

	// 写入未发送日志文件
	logger := log.New(f, "UNSENT: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Printf("Email: %s, Reason: %s\n", email, reason)
}

func SendErrorEmail(err error) {
	subject := "系统报错通知"
	body := fmt.Sprintf("管理员您好，\n\n系统在处理过程中出现以下错误：\n\n%s\n\n请及时处理。", err)
	SendEmail(ADMIN_EMAIL, subject, body)
}

func SendEmail(to, subject, body string) error {
	from := SMTP_USERNAME
	password := SMTP_PASSWORD
	smtpHost := SMTP_HOST
	smtpPort := SMTP_PORT

	// 设置邮件头信息
	header := make(map[string]string)
	header["From"] = fmt.Sprintf("%s <%s>", SMTP_NAME, from)
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"

	// 构建邮件内容
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// 设置SMTP服务器地址和端口
	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)

	// 设置TLS配置
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	// 连接到SMTP服务器
	conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		log.Printf("Error connecting to SMTP server: %v", err)
		return err
	}

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Printf("Error creating SMTP client: %v", err)
		return err
	}

	// 使用AUTH命令进行认证
	auth := smtp.PlainAuth("", from, password, smtpHost)
	if err = client.Auth(auth); err != nil {
		log.Printf("Error authenticating: %v", err)
		return err
	}

	// 设置发件人和收件人
	if err = client.Mail(from); err != nil {
		log.Printf("Error setting sender: %v", err)
		return err
	}

	if err = client.Rcpt(to); err != nil {
		log.Printf("Error setting recipient: %v", err)
		return err
	}

	// 获取邮件写入器
	writer, err := client.Data()
	if err != nil {
		log.Printf("Error getting writer: %v", err)
		return err
	}

	// 写入邮件内容
	_, err = writer.Write([]byte(message))
	if err != nil {
		log.Printf("Error writing message: %v", err)
		return err
	}

	// 关闭写入器
	if err = writer.Close(); err != nil {
		log.Printf("Error closing writer: %v", err)
		return err
	}

	// 发送QUIT命令以关闭连接
	if err = client.Quit(); err != nil {
		log.Printf("Error quitting SMTP client: %v", err)
		return err
	}

	LogInfo(fmt.Sprintf("Email successfully sent to %s", to))
	return nil
}

func RetryUnsentEmails() {
	file, err := os.Open(UNSENT_LOG_FILE)
	if err != nil {
		LogError(fmt.Errorf("Failed to open unsent log file: %v", err))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var unsentEntries []string
	for scanner.Scan() {
		unsentEntries = append(unsentEntries, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		LogError(fmt.Errorf("Error reading unsent log file: %v", err))
		return
	}

	for _, entry := range unsentEntries {
		parts := strings.Split(entry, ",")
		if len(parts) < 2 {
			continue
		}

		email := strings.TrimSpace(parts[0][7:])  // Remove "Email: " prefix
		reason := strings.TrimSpace(parts[1][8:]) // Remove "Reason: " prefix

		// 根据 reason 判断是哪种邮件
		if strings.Contains(reason, "birthday") {
			employee := models.Employee{Email: email}
			err := SendBirthdayEmail(employee)
			if err == nil {
				LogInfo(fmt.Sprintf("Successfully resent birthday email to %s", email))
			} else {
				LogError(fmt.Errorf("Failed to resend birthday email to %s: %v", email, err))
			}
		} else if strings.Contains(reason, "anniversary") {
			employee := models.Employee{Email: email}
			err := SendWorkAnniversaryEmail(employee)
			if err == nil {
				LogInfo(fmt.Sprintf("Successfully resent anniversary email to %s", email))
			} else {
				LogError(fmt.Errorf("Failed to resend anniversary email to %s: %v", email, err))
			}
		}
	}
}
