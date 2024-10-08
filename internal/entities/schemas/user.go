package schemas

type AddUsers struct {
	UserId      string `json:"user_id"`      //ผู้ใช้งาน
	Password    string `json:"password"`     //รหัสผ่าน
	Name        string `json:"name"`         //ชื่อ
	SurName     string `json:"sur_name"`     //นามสกุล
	Email       string `json:"email"`        //อีเมล
	BirthDay    string `json:"birth_day"`    //วันเกิด
	PhoneNumber string `json:"phone_number"` //เบอร์โทร
	IdCard      string `json:"id_card"`      //เลขบัตรประจำตัว
}

type FindUsersReq struct {
	UserId  string `json:"user_id" form:"user_id"`   //ผู้ใช้งาน
	Name    string `json:"name" form:"name"`         //ชื่อ
	SurName string `json:"sur_name" form:"sur_name"` //นามสกุล
	Email   string `json:"email" form:"email"`       //อีเมล
}

type FindUsersByUserIdReq struct {
	UserId string `json:"user_id" form:"user_id"` //ผู้ใช้งาน
}

type ValueReq struct {
	Value string `form:"value"` //ค่า string ใดๆ
}

type LoginReq struct {
	UserId   string `json:"user_id"`  //ผู้ใช้งาน
	Password string `json:"password"` //รหัสผ่าน
}
