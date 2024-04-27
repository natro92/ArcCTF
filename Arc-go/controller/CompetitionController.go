package controller

import (
	"Arc/common"
	"Arc/dto"
	"Arc/model"
	"Arc/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
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
	//log.Printf("title: %s, description: %s, category: %s, tags: %s, active: %d, startTime: %d, endTime: %d, password: %s", title, description, category, tags, active, startTime, endTime, password)
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

// JoinCompetition
//
//	@Description: 加入比赛
//	@param context
func JoinCompetition(context *gin.Context) {
	db := common.GetDB()
	user, _ := context.Get("user")
	competitionUserInput := context.Query("id")

	// * 检测当前用户是否已经注册过
	var userm model.User
	db.Where("id = ?", int(user.(model.User).ID)).First(&userm)
	if userm.ID == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}
	if userm.TeamId == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "用户未加入队伍")
		return
	}
	// * 获取team参数
	var team model.Team
	db.Where("id = ?", userm.TeamId).First(&team)
	if team.ID == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "队伍不存在")
		return
	}

	// * 检测比赛是否存在
	if len(competitionUserInput) == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "比赛参数为空")
		return
	}
	// * 检测team是否已经加入比赛
	if len(team.CompetitionId) != 0 {
		for _, v := range team.CompetitionId {
			if v == competitionUserInput {
				response.Response(context, http.StatusUnprocessableEntity, 422, nil, "队伍已经加入比赛")
				return
			}
		}
	}
	var competition model.Competition
	db.Where("id = ?", competitionUserInput).First(&competition)
	if competition.ID == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "比赛不存在")
		return
	}
	// * 检测比赛是否已经开始 这个我觉得不太需要吧
	//if competition.StartTime.After(time.Now()) {
	//	response.Response(context, http.StatusUnprocessableEntity, 422, nil, "比赛已经开始")
	//	return
	//}
	// * 检测比赛是否已经结束
	if competition.EndTime.Before(time.Now()) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "比赛已经结束")
		return
	}
	//// TODO 检测队伍是否已经满员 加个配置
	//if team.MemberNum >= 4 {
	//	response.Response(context, http.StatusUnprocessableEntity, 422, nil, "队伍已满员")
	//	return
	//}
	// * 加入比赛
	team.CompetitionId = append(team.CompetitionId, strconv.Itoa(int(competition.ID)))
	db.Save(&team)
	// * 返回结果
	response.Success(context, nil, "加入比赛成功")
}
