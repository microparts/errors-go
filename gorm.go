package errors

import "github.com/jinzhu/gorm"

// Извлечение ошибок из gorm запросов (для update|delete)
func GetResultErrors(result *gorm.DB) []error {
	err := result.GetErrors()
	affected := result.RowsAffected
	if affected == 0 {
		err = append(err, NoRowsAffected)
	}
	return err
}
