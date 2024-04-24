package controller

import (
	"Arc/common"
	"Arc/dto"
	"Arc/model"
	"Arc/response"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// CreateCompetition
// ! 管理员权限
//
//	@Description: 创建比赛
//	@param context
func CreateCompetition(context *gin.Context) {
	// * 检测role是否为管理员
	role, _ := context.Get("role")
	if role != 0 && role != 1 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "权限不足，您不是管理员权限")
		return
	}
	db := common.GetDB()
	// * 获取参数
	var requestCompetition = model.Competition{}
	err := context.ShouldBind(&requestCompetition)
	title := requestCompetition.Title
	description := requestCompetition.Description
	category := requestCompetition.Category
	tags := requestCompetition.Tags
	active := requestCompetition.Active
	startTime := requestCompetition.StartTime
	endTime := requestCompetition.EndTime
	password := requestCompetition.Password
	if err != nil {
		return
	}
	log.Printf("title: %s, description: %s, category: %s, tags: %s, active: %d, startTime: %d, endTime: %d, password: %s", title, description, category, tags, active, startTime, endTime, password)
	// * 参数校验
	if len(title) == 0 || len(description) == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "参数不能为空")
		return
	}
	// * 创建比赛
	newCompetition := model.Competition{
		Title:           title,
		Description:     description,
		Category:        category,
		Tags:            tags,
		Active:          active,
		StartTime:       startTime,
		EndTime:         endTime,
		Password:        password,
		ParticipantsNum: 0,
	}
	db.Create(&newCompetition)
	// * 返回结果
	response.Success(context, nil, "创建比赛成功")
}

func GetAllCompetitions(context *gin.Context) {
	db := common.GetDB()
	var competitions []model.Competition
	db.Find(&competitions)
	response.Success(context, gin.H{"competitions": dto.ToCompetitionsListDto(competitions)}, "获取比赛列表成功")
}

// GetCompetitionInfo4User
// * 用户权限
//
//	@Description: 获取比赛信息
//	@param context
func GetCompetitionInfo4User(context *gin.Context) {
	// * 获取参数
	id := context.Query("id")
	if len(id) == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "参数不能为空")
		return
	}
	db := common.GetDB()
	var competition model.Competition
	db.Where("id = ?", id).First(&competition)
	response.Success(context, gin.H{"competition": dto.ToCompetitionDto(competition)}, "获取比赛信息成功")
}
