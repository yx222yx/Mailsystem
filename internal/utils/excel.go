package utils

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
	"mailsys/internal/models"
	"mailsys/internal/services"
)

func ReadExcel(filePath string) ([]models.Employee, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		logError := fmt.Errorf("Error opening Excel file: %v", err)
		services.LogError(logError)
		return nil, logError
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		logError := fmt.Errorf("Error reading Excel rows: %v", err)
		services.LogError(logError)
		return nil, logError
	}

	var employees []models.Employee
	errorOccurred := false

	for i, row := range rows {
		if i == 0 {
			continue
		}

		if len(row) < 5 {
			logError := fmt.Errorf("Row %d is invalid, skipping: %v", i+1, row)
			services.LogError(logError)
			errorOccurred = true
			continue
		}

		joinDate, err := time.Parse("2006-01-02", row[2])
		if err != nil {
			logError := fmt.Errorf("Error parsing join date in row %d: %v", i+1, err)
			services.LogError(logError)
			errorOccurred = true
			continue
		}

		birthDate, err := time.Parse("2006-01-02", row[3])
		if err != nil {
			logError := fmt.Errorf("Error parsing birth date in row %d: %v", i+1, err)
			services.LogError(logError)
			errorOccurred = true
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

	if !errorOccurred {
		services.LogInfo("Excel data read successfully without errors")
	}

	return employees, nil
}
