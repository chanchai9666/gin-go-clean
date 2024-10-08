package models

import (
	"errors"

	"github.com/chanchai9666/aider"
	"gorm.io/gorm"
)

// ข้อมูล User
type Users struct {
	BaseColumn
	UserId      string  `json:"user_id" gorm:"primaryKey"`                                    //ไอดี ของผู้ใช้งาน
	Password    string  `json:"password" gorm:"type:varchar(255);comment:รหัสผ่าน"`           //ชื่อ โปรไฟล์
	Name        string  `json:"name" gorm:"type:varchar(50);comment:ชื่อ"`                    //ชื่อ
	SurName     string  `json:"sur_name" gorm:"type:varchar(50);comment:นามสกุล"`             //นามสกุล
	Birthday    *string `json:"birth_day" gorm:"type:date;comment:วันเกิด;default:null;"`     //วันเกิด
	PhoneNumber *string `json:"phone_number" gorm:"type:varchar(20);comment:หมายเลขโทรศัพท์"` //หมายเลขโทรศัพท์
	IdCard      *string `json:"id_card" gorm:"type:varchar(30);comment:รหัสบัตรประจำตัว"`     //รหัสบัตรประจำตัว
}

func (u *Users) BeforeCreate(tx *gorm.DB) (err error) {
	// เรียก BeforeCreate ของ BaseColumn
	if err := u.BaseColumn.BeforeCreate(tx); err != nil {
		return err
	}

	if u.Password == "" {
		u.Password = u.UserId
	}
	HashPassword, err := aider.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = HashPassword
	return
}

func (u *Users) BeforeUpdate(tx *gorm.DB) (err error) {
	// if tx.Statement.Changed("DeletedAt") {
	// 	// ถ้ากำลังทำ Soft Delete ไม่ต้องเช็ค UserId
	// 	return nil
	// }
	// ตรวจสอบว่ากำลังลบหรือไม่
	if u.UserId == "" {
		return errors.New("UserId is Empty")
	}
	return nil
}

func (u *Users) BeforeDelete(tx *gorm.DB) (err error) {
	// อัปเดตฟิลด์ IsActive และ DeletedAt โดยใช้ฟังก์ชัน Updates ของ GORM
	return nil
}

// func (u *Users) BeforeDelete(tx *gorm.DB) (err error) {
// 	// เรียก BeforeCreate ของ BaseColumn
// 	if err := u.BaseColumn.BeforeDelete(tx); err != nil {
// 		return err
// 	}
// 	return
// }

type UserLevels struct {
	BaseColumn
	ID uint `gorm:"column:id;type:int;primaryKey;autoIncrement"`
}
