package services

import (
	"bytes"
	"fmt"
	"html/template"
	"mailsys/internal/models"
	"os"
	"time"
)

func readTemplate(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func SendBirthdayEmail(employee models.Employee) error {
	templateContent, err := readTemplate("D:/go_projects/mailsys/templates/birthday_template.html")
	if err != nil {
		LogError(fmt.Errorf("failed to read birthday template: %v", err))
		LogUnsent(employee.Email, "failed to read birthday template")
		return err
	}

	tmpl, err := template.New("birthday").Parse(templateContent)
	if err != nil {
		LogError(fmt.Errorf("failed to parse birthday template: %v", err))
		LogUnsent(employee.Email, "failed to parse birthday template")
		return err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, employee)
	if err != nil {
		LogError(fmt.Errorf("failed to execute birthday template: %v", err))
		LogUnsent(employee.Email, "failed to execute birthday template")
		return err
	}

	subject := "生日快乐"
	err = SendEmail(employee.Email, subject, body.String())
	if err == nil {
		LogInfo(fmt.Sprintf("Birthday email sent to %s", employee.Name))
	} else {
		LogUnsent(employee.Email, "birthday email send failure")
	}
	return err
}

func SendWorkAnniversaryEmail(employee models.Employee) error {
	templateContent, err := readTemplate("D:/go_projects/mailsys/templates/anniversary_template.html")
	if err != nil {
		LogError(fmt.Errorf("failed to read anniversary template: %v", err))
		LogUnsent(employee.Email, "failed to read anniversary template")
		return err
	}

	tmpl, err := template.New("anniversary").Parse(templateContent)
	if err != nil {
		LogError(fmt.Errorf("failed to parse anniversary template: %v", err))
		LogUnsent(employee.Email, "failed to parse anniversary template")
		return err
	}

	data := struct {
		Name  string
		Years int
	}{
		Name:  employee.Name,
		Years: time.Now().Year() - employee.JoinDate.Year(),
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		LogError(fmt.Errorf("failed to execute anniversary template: %v", err))
		LogUnsent(employee.Email, "failed to execute anniversary template")
		return err
	}

	subject := "入职周年快乐"
	err = SendEmail(employee.Email, subject, body.String())
	if err == nil {
		LogInfo(fmt.Sprintf("Work anniversary email sent to %s", employee.Name))
	} else {
		LogUnsent(employee.Email, "anniversary email send failure")
	}
	return err
}
