package validate

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

var validate = validator.New()

type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse represents the structure of a JSON error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func NewSuccessMessage() *Message {
	return &Message{
		Code:    200,
		Message: "Success",
	}
}

func RespJson(c *gin.Context, fn interface{}, input interface{}) {
	// ตรวจสอบและดึงข้อมูลจาก request ตาม content type
	if err := parseInputData(c, input); err != nil {
		RenderJSON(c, err, nil)
		return
	}

	// ตรวจสอบความถูกต้องของข้อมูล
	if err := validateInput(input); err != nil {
		RenderJSON(c, err, nil)
		return
	}

	fnType := reflect.TypeOf(fn)
	var out []reflect.Value

	// ตรวจสอบจำนวนพารามิเตอร์ที่ฟังก์ชัน fn ต้องการ
	if fnType.NumIn() == 0 {
		// ถ้าฟังก์ชันไม่ต้องการพารามิเตอร์
		out = reflect.ValueOf(fn).Call(nil)
	} else if fnType.NumIn() == 1 {
		// ถ้าฟังก์ชันต้องการพารามิเตอร์ 1 ตัว
		out = reflect.ValueOf(fn).Call([]reflect.Value{
			reflect.ValueOf(input),
		})
	} else if fnType.NumIn() == 2 {
		// ถ้าฟังก์ชันต้องการพารามิเตอร์ 2 ตัว (เช่น context และ input)
		out = reflect.ValueOf(fn).Call([]reflect.Value{
			reflect.ValueOf(c),
			reflect.ValueOf(input),
		})
	} else {
		// กรณีจำนวนพารามิเตอร์ไม่ตรง
		RenderJSON(c, fmt.Errorf("invalid function signature"), nil)
		return
	}

	// ตรวจสอบผลลัพธ์
	errObj := out[1].Interface()
	if errObj != nil {
		logrus.Errorf("call service error: %s", errObj)
		RenderJSON(c, errObj.(error), nil)
		return
	}

	ResponseResult := ResponseResult{
		Code:    200,
		Message: "Success",
		Result:  out[0].Interface(),
	}

	RenderJSON(c, nil, ResponseResult)
}

func RespJsonNoReq(c *gin.Context, fn interface{}) {

	out := reflect.ValueOf(fn).Call(nil)

	// ตรวจสอบผลลัพธ์
	errObj := out[1].Interface()
	if errObj != nil {
		logrus.Errorf("call service error: %s", errObj)
		RenderJSON(c, errObj.(error), nil)
		return
	}

	ResponseResult := ResponseResult{
		Code:    200,
		Message: "Success",
		Result:  out[0].Interface(),
	}

	RenderJSON(c, nil, ResponseResult)
}

func RespSuccess(c *gin.Context, fn interface{}, input interface{}) {
	// ตรวจสอบและดึงข้อมูลจาก request ตาม content type
	if err := parseInputData(c, input); err != nil {
		RenderJSON(c, err, nil)
		return
	}

	// ตรวจสอบความถูกต้องของข้อมูล
	if err := validateInput(input); err != nil {
		RenderJSON(c, err, nil)
		return
	}

	fnType := reflect.TypeOf(fn)
	var out []reflect.Value

	// ตรวจสอบจำนวนพารามิเตอร์ที่ฟังก์ชัน fn ต้องการ
	if fnType.NumIn() == 0 {
		// ถ้าฟังก์ชันไม่ต้องการพารามิเตอร์
		out = reflect.ValueOf(fn).Call(nil)
	} else if fnType.NumIn() == 1 {
		// ถ้าฟังก์ชันต้องการพารามิเตอร์ 1 ตัว
		out = reflect.ValueOf(fn).Call([]reflect.Value{
			reflect.ValueOf(input),
		})
	} else if fnType.NumIn() == 2 {
		// ถ้าฟังก์ชันต้องการพารามิเตอร์ 2 ตัว (เช่น context และ input)
		out = reflect.ValueOf(fn).Call([]reflect.Value{
			reflect.ValueOf(c),
			reflect.ValueOf(input),
		})
	} else {
		// กรณีจำนวนพารามิเตอร์ไม่ตรง
		RenderJSON(c, fmt.Errorf("invalid function signature"), nil)
		return
	}
	// ตรวจสอบผลลัพธ์
	errObj := out[0].Interface()
	if errObj != nil {
		logrus.Errorf("call service error: %s", errObj)
		RenderJSON(c, errObj.(error), nil)
		return
	}

	RenderJSON(c, nil, NewSuccessMessage())
}

// ฟังก์ชันนี้ใช้ในการเช็คและดึงข้อมูลจาก request ตาม content type
func parseInputData(c *gin.Context, input interface{}) error {
	Method := c.Request.Method
	switch {
	case strings.HasPrefix(Method, "POST"), strings.HasPrefix(Method, "DELETE"):
		// เรียกใช้ฟังก์ชันแปลง
		if err := parseRequestBody(c, &input); err != nil {
			return err
		}
	case strings.HasPrefix(Method, "GET"):
		// ถ้าต้องการจัดการกับ query parameters ด้วย
		if err := c.ShouldBindQuery(input); err != nil {
			return err
		}

	case c.ContentType() == "multipart/form-data":
		if err := mapFormValues(c, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported content type: %s", Method)
	}

	// แปลงค่า path parameters ให้เข้าไปยัง struct `input`
	val := reflect.ValueOf(input).Elem()
	for _, param := range c.Params {
		field := val.FieldByName(toPascalCase(param.Key)) // หาชื่อฟิลด์ที่ตรงกับ path parameter
		if field.IsValid() && field.CanSet() {
			field.SetString(param.Value) // ตั้งค่าฟิลด์จาก path parameter
		}
	}

	return nil
}

func mapFormValues(c *gin.Context, input interface{}) error {
	val := reflect.ValueOf(input).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name

		value := c.PostForm(fieldName)
		if value != "" {
			// If the field is a slice, assign the value as a single element
			if field.Type().Kind() == reflect.Slice {
				slice := reflect.MakeSlice(field.Type(), 1, 1)
				slice.Index(0).SetString(value)
				field.Set(slice)
			} else {
				field.SetString(value)
			}
		}
	}

	return nil
}

func parseRequestBody(c *gin.Context, inputStruct interface{}) error {
	// ใช้คำสั่ง ShouldBind ในที่นี้
	if err := c.ShouldBind(inputStruct); err != nil {
		return fmt.Errorf("failed to parse request body: %v", err)
	}

	return nil
}

// ฟังก์ชันที่จัดการ error response และ success response
func RenderJSON(c *gin.Context, err error, successResponse interface{}) {
	if err != nil {
		// Create an ErrorResponse instance
		errorResponse := ErrorResponse{
			Code:    400, // or you can use a specific error code
			Message: err.Error(),
		}

		// Set the HTTP status code
		c.JSON(400, errorResponse)
		return
	}

	// Return the success response as JSON
	c.JSON(200, successResponse)
}

// CustomValidator คือ struct ที่ใช้เพื่อสร้าง custom validator
type CustomValidator struct {
	validator *validator.Validate
}

// ValidateDate คือ custom validation function สำหรับการตรวจสอบวันที่
func ValidateDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	// ตรวจสอบว่าค่าวันที่ไม่เป็นค่าว่าง
	if dateStr == "" {
		return true
	}
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// NewValidator คือฟังก์ชั่นที่ใช้สร้าง custom validator
func NewValidator() *CustomValidator {
	v := validator.New()
	v.RegisterValidation("date", ValidateDate)

	return &CustomValidator{validator: v}
}

// ฟังก์ชันนี้ใช้ในการตรวจสอบความถูกต้องของข้อมูล
func validateInput(input interface{}) error {
	// สร้าง custom validator
	validate := NewValidator()

	// ใช้ custom validator เพื่อ validate ข้อมูล
	if err := validate.validator.Struct(input); err != nil {
		// กรณีมี error ในการ validate แสดงข้อความเพิ่มเติม
		errs := err.(validator.ValidationErrors)
		errorMsg := "Invalid request data:"
		for _, e := range errs {
			errorMsg += fmt.Sprintf("\n- Field: %s, Type: %T, Error: %s", e.Field(), e.Value(), e.Tag())
		}
		// ใช้ gin.Error เพื่อสร้าง error ที่เข้ากับ Gin framework
		return fmt.Errorf("%s", errorMsg)
	}
	return nil
}

// ฟังก์ชันนี้จะแปลง snake_case เป็น PascalCase
func toPascalCase(snake string) string {
	// แยก string ตามขีดล่าง _
	parts := strings.Split(snake, "_")
	for i, part := range parts {
		// แปลงตัวอักษรแรกของแต่ละส่วนให้เป็นตัวพิมพ์ใหญ่
		parts[i] = strings.Title(part)
	}
	// รวมกลับมาเป็น string เดียวในรูปแบบ PascalCase
	return strings.Join(parts, "")
}

// ฟังก์ชันนี้จะแปลง snake_case เป็น camelCase
func toCamelCase(snake string) string {
	pascal := toPascalCase(snake)
	// แปลงตัวอักษรแรกให้เป็นตัวพิมพ์เล็กสำหรับ camelCase
	return string(unicode.ToLower(rune(pascal[0]))) + pascal[1:]
}
