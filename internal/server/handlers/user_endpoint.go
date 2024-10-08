package handlers

import (
	"github.com/gin-gonic/gin"

	_ "arczed/internal/entities/models"
	"arczed/internal/entities/schemas"
	"arczed/internal/server/validate"
	"arczed/internal/usecase"
)

type userEndPoint struct {
	service usecase.UsersService
}

func NewUserEndPoint(service usecase.UsersService) *userEndPoint {
	return &userEndPoint{service: service}
}

// @Tags Users
// @Summary ค้นหา User ตามเงื่อนไข
// @Description Show User ตามเงื่อนไข
// @Accept  json
// @Produce  json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Param request query schemas.FindUsersReq false " request body "
// @Success 200 {object} []models.Users
// @Failure 400 {object} schemas.HTTPError
// @Router /api/users [get]
// @Security ApiKeyAuth
func (en *userEndPoint) FindUser(c *gin.Context) {
	validate.RespJson(c, en.service.FindUsers, &schemas.FindUsersReq{})
}

// @Tags Users
// @Summary ค้นหา User ตาม UserId
// @Description Show User ตาม UserId
// @Accept  json
// @Produce  json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Param user_id path string true "User ID"
// @Success 200 {object} models.Users
// @Failure 400 {object} schemas.HTTPError
// @Router /api/users/{user_id} [get]
// @Security ApiKeyAuth
func (en *userEndPoint) FindUsersByUserId(c *gin.Context) {
	validate.RespJson(c, en.service.FindUsersByUserId, &schemas.FindUsersByUserIdReq{})
}

// @Tags Users
// @Summary ค้นหา User ทั้งหมด
// @Description Show User ทั้งหมด
// @Accept  json
// @Produce  json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Success 200 {object} []models.Users
// @Failure 400 {object} schemas.HTTPError
// @Router /api/users/usersAll [get]
// @Security ApiKeyAuth
func (en *userEndPoint) FindUserAll(c *gin.Context) {
	validate.RespJsonNoReq(c, en.service.FindUsersAll)
}

// @Tags Users
// @Summary เพิ่มข้อมูล User
// @Description เพิ่มข้อมูล User
// @Accept  json
// @Produce  json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Param request body schemas.AddUsers false " request body "
// @Success 200 {object} schemas.HTTPError
// @Failure 400 {object} schemas.HTTPError
// @Router /api/users/createUsers [post]
// @Security ApiKeyAuth
func (en *userEndPoint) CreateUsers(c *gin.Context) {
	validate.RespSuccess(c, en.service.CreateUsers, &schemas.AddUsers{})
}

// @Tags Users
// @Summary แก้ไขข้อมูล User
// @Description แก้ไขข้อมูล User
// @Accept  json
// @Produce  json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Param request body schemas.AddUsers false " request body "
// @Success 200 {object} schemas.HTTPError
// @Failure 400 {object} schemas.HTTPError
// @Router /api/users/updateUsers [post]
// @Security ApiKeyAuth
func (en *userEndPoint) UpdateUsers(c *gin.Context) {
	validate.RespSuccess(c, en.service.UpdateUsers, &schemas.AddUsers{})
}

// @Tags Users
// @Summary ลบข้อมูล User
// @Description ลบข้อมูล User
// @Accept  json
// @Produce  json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Param user_id path string true "User ID"
// @Success 200 {object} schemas.HTTPError
// @Failure 400 {object} schemas.HTTPError
// @Router /api/users/deleteUsers/{user_id} [delete]
// @Security ApiKeyAuth
func (en *userEndPoint) DeleteUsers(c *gin.Context) {
	validate.RespSuccess(c, en.service.DeleteUsers, &schemas.AddUsers{})
}
