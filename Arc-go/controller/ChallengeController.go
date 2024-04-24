package controller

import (
	"Arc/common"
	"Arc/dto"
	"Arc/model"
	"Arc/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateChallenge
// ! 管理员权限
//
//	@Description: 创建题目
//	@param context
func CreateChallenge(context *gin.Context) {
	// * 检测role是否为管理员
	role, _ := context.Get("role")
	if role != 0 && role != 1 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "权限不足，您不是管理员权限")
		return
	}
	db := common.GetDB()
	// * 获取参数
	var requestChallenge = model.Challenge{}
	err := context.ShouldBind(&requestChallenge)
	competitionId := requestChallenge.CompetitionID
	title := requestChallenge.Title
	description := requestChallenge.Description
	score := requestChallenge.Score
	minScore := requestChallenge.MinScore
	maxScore := requestChallenge.MaxScore
	containerMirror := requestChallenge.ContainerMirror
	flag := requestChallenge.Flag
	attachment := requestChallenge.Attachment
	category := requestChallenge.Category
	tags := requestChallenge.Tags
	hints := requestChallenge.Hints
	visible := requestChallenge.Visible
	if err != nil {
		return
	}
	// * 参数校验
	// 检测是否缺少参数，或者传入参数格式错误
	if len(title) == 0 || len(description) == 0 || len(flag) == 0 || len(category) == 0 || len(tags) == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "参数不能为空")
		return
	}
	// 检测分数是否合法
	if score < 0 || minScore < 0 || maxScore < 0 || minScore > maxScore {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "分数不合法")
		return
	}
	// * 创建题目
	newChallenge := model.Challenge{
		CompetitionID:   competitionId,
		Title:           title,
		Description:     description,
		Score:           score,
		MinScore:        minScore,
		MaxScore:        maxScore,
		ContainerMirror: containerMirror,
		Flag:            flag,
		Attachment:      attachment,
		Category:        category,
		Tags:            tags,
		Hints:           hints,
		Visible:         visible,
	}
	db.Create(&newChallenge)
	// * 返回结果
	response.Success(context, nil, "创建题目成功")
}

// GetChallengeInfo4User
// * User权限
//
// @Description: 获取比赛信息
// @param context
func GetChallengeInfo4User(context *gin.Context) {
	challengeID := context.Query("id")
	db := common.GetDB()
	var challenge model.Challenge
	db.Where("id = ?", challengeID).First(&challenge)
	response.Success(context, gin.H{"challenge": dto.ToChallengeDto(challenge)}, "获取比赛信息成功")
}

// GetAllChallengesByCompetitionId
// * user权限
//
//	@Description: 获取比赛下的所有题目
//	@param context
func GetAllChallengesByCompetitionId(context *gin.Context) {
	competitionID := context.Query("id")
	db := common.GetDB()
	var challenges []model.Challenge
	db.Where("competition_id = ?", competitionID).Find(&challenges)
	response.Success(context, gin.H{"challenges": dto.ToChallengesListDto(challenges)}, "获取比赛下的所有题目成功")
}
