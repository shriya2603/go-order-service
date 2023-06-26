package app

import (
	"fmt"

	"gorm.io/gorm"
)

const (
	INVALID_TABLENAME = "Invalid Tablename"
)

// isEmpty func: checks whether a string is empty or not
func isEmpty(str string) bool {
	if str == "" || len(str) == 0 {
		return true
	}
	return false
}

func validID(id uint, tablename string, db *gorm.DB) bool {
	var qresult string
	query := fmt.Sprintf("SELECT ID FROM %s WHERE ID=%d;", tablename, id)
	db.Raw(query).Scan(&qresult)
	if qresult == "" {
		return false
	}
	return true
}
