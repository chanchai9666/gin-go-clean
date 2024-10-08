package models

import (
	"time"

	"github.com/chanchai9666/aider"
	"gorm.io/gorm"
)

// type DeletedAt sql.NullTime
type BaseColumn struct {
	IsActive    *int            `json:"is_active" gorm:"type:int2;default:1;comment:สถานะใช้งาน;"`                 //สถานะใช้งาน
	CreatedAt   *time.Time      `json:"created_at,omitempty" gorm:"type:timestamp;comment:วันที่สร้าง"`            //วันที่สร้าง
	CreatedUser string          `json:"created_user,omitempty" gorm:"type:varchar(50);comment:ผู้สร้าง"`           //ผู้สร้าง
	UpdatedAt   *time.Time      `json:"updated_at,omitempty" gorm:"type:timestamp;comment:วันเวลาที่อัพเดทล่าสุด"` //วันเวลาที่อัพเดทล่าสุด
	UpdatedUser string          `json:"updated_user,omitempty" gorm:"type:varchar(50);comment:ผู้อัพเดทล่าสุด"`    //ผู้อัพเดทล่าสุด
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index;type:timestamp;comment:วันเวลาที่ลบ"`     //วันเวลาที่ลบ
	DeletedUser string          `json:"deleted_user,omitempty" gorm:"type:varchar(50);comment:ผุ้ลบ"`              //ผุ้ลบ
}

// ใช้ GORM Hooks (BeforeCreate) กรณีที่ไม่ได้กำหนดค่าส่งมา จะ default value ตามที่กำหนดไว้ก่อนนำไป insert
// กรณี struct มี Hooks เป็นของตัวเอง Hooks ที่เป็นของ embedded struct จะไม่ถูกเรียกใช้งาน ถ้าต้องการเรียกใช้งานต้องเรียกเอง
func (u *BaseColumn) BeforeCreate(tx *gorm.DB) (err error) {
	if u.IsActive == nil {
		isActive := 1
		u.IsActive = &isActive
	}

	timeNow := aider.TimeTimeNow()
	if u.CreatedAt == nil {
		u.CreatedAt = &timeNow
	}

	if u.UpdatedAt == nil {
		u.UpdatedAt = &timeNow
	}
	return nil
}

// ใช้ GORM Hooks (BeforeCreate) กรณีที่ไม่ได้กำหนดค่าส่งมา จะ default value ตามที่กำหนดไว้ก่อนนำไป Update
func (u *BaseColumn) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.UpdatedAt == nil {
		timeNow := aider.TimeTimeNow()
		u.UpdatedAt = &timeNow
	}
	u.UpdatedUser = "BeforeUpdate"
	return nil
}

func (u *BaseColumn) BeforeDelete(tx *gorm.DB) (err error) {
	if u.IsActive == nil {
		isActive := 0
		u.IsActive = &isActive
	}

	timeNow := aider.TimeTimeNow() //2024-10-08T15:01:09Z
	deletedAt := gorm.DeletedAt{
		Time:  timeNow,
		Valid: true,
	}
	u.DeletedAt = &deletedAt

	return
}

// // AfterFind Hook ที่จะถูกเรียกหลังจากที่ GORM ทำการค้นหาเสร็จ
// func (bc *BaseColumn) AfterFind(tx *gorm.DB) (err error) {

// 	// ตรวจสอบว่าฟิลด์ CreatedAt, UpdatedAt, DeletedAt ไม่เป็น nil
// 	if bc.CreatedAt != nil {
// 		formattedCreatedAt := bc.CreatedAt.Format("2006-01-02 15:04:05")
// 		tt := aider.DateTime(formattedCreatedAt)
// 		bc.CreatedAt = &tt
// 		fmt.Println("Formatted CreatedAt:", bc.CreatedAt)
// 	}
// 	if bc.UpdatedAt != nil {
// 		formattedUpdatedAt := bc.UpdatedAt.Format("2006-01-02 15:04:05")
// 		fmt.Println("Formatted UpdatedAt:", formattedUpdatedAt)
// 	}
// 	if bc.DeletedAt != nil && bc.DeletedAt.Time != (time.Time{}) {
// 		formattedDeletedAt := bc.DeletedAt.Time.Format("2006-01-02 15:04:05")
// 		fmt.Println("Formatted DeletedAt:", formattedDeletedAt)
// 	}
// 	return nil
// }
