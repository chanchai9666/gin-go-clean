package repositories

import (
	"fmt"

	"gorm.io/gorm"
)

func WhereIsActive(table ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		condition := "is_active=?"
		if len(table) > 0 && table[0] != "" {
			condition = table[0] + ".is_active=?"
		}
		return db.Where(condition, 1)
	}
}

func WhereUserId(data ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(data) == 0 {
			return db // ไม่มีค่า input คืนค่า db กลับไปโดยไม่ทำอะไร
		}
		value := data[0]
		condition := "user_id = ?"
		if len(data) > 1 {
			tableName := data[0]
			value = data[1]
			condition = fmt.Sprintf("%s.user_id = ?", tableName)
		}

		return db.Where(condition, value)
	}
}
