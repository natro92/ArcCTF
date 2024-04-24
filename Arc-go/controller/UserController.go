package controller

import (
	"Arc/common"
	"Arc/dto"
	"Arc/model"
	"Arc/response"
	"Arc/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func Register(context *gin.Context) {
	db := common.GetDB()
	// ! 不能用form 因为穿的使json
	var requestUser = model.User{}
	// * 获取参数
	err := context.ShouldBind(&requestUser)
	if err != nil {
		return
	}
	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password

	// * 数据验证
	if len(telephone) != 11 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	// * 如果名称没有上传 则生成一个10位的随机字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	// * 判断手机号是否存在
	if isTelephoneExist(db, telephone) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "手机号已经注册过，用户已存在")
		return
	}
	// * 新建用户
	// * 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(context, http.StatusInternalServerError, 500, nil, "加密错误")
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashedPassword),
	}
	db.Create(&newUser)
	//	* 发放token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(context, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}
	response.Success(context, gin.H{"token": token}, "注册成功")
}

func Login(context *gin.Context) {
	var requestUser = model.User{}
	db := common.GetDB()
	err := context.ShouldBind(&requestUser)
	if err != nil {
		return
	}
	//	* 获取参数
	telephone := requestUser.Telephone
	password := requestUser.Password
	//	* 数据验证
	if len(telephone) != 11 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	//	* 判断手机号是否存在
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}
	//	* 判断密码是否正确
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		response.Fail(context, nil, "密码错误")
		return
	}
	//	* 发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(context, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}
	response.Success(context, gin.H{"token": token}, "登录成功")
}

func Info(context *gin.Context) {
	user, _ := context.Get("user")
	response.Success(context, gin.H{"user": dto.ToUserDto(user.(model.User))}, "")
	//context.JSON(http.StatusOK, gin.H{
	//	"code": 200,
	//	"data": gin.H{"user": dto.ToUserDto(user.(model.User))},
	//})
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
