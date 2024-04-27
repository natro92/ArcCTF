package controller

import (
	"Arc/common"
	"Arc/dto"
	"Arc/model"
	"Arc/response"
	"Arc/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

// CreateNewTeam
//
//	@Description: 创建一个新的队伍
//	@param context
//	@return *model.Team
func CreateNewTeam(context *gin.Context) {
	db := common.GetDB()
	var requestTeam = model.Team{}
	err := context.ShouldBind(&requestTeam)
	if err != nil {
		return
	}
	name := requestTeam.Name
	//  * 数据验证
	if len(name) == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "内容参数不能为空")
		return
	}
	//  * 判断队伍是否存在
	if isTeamExist(db, name) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "队伍已经存在")
		return
	}
	user, _ := context.Get("user")
	//  * 检测当前用户是否已经注册过
	var userm model.User
	db.Where("id = ?", int(user.(model.User).ID)).First(&userm)
	if userm.ID == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}
	if userm.TeamId != 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "用户已经加入队伍")
		return
	}
	//  * 设置创建人为队长
	//  * 新建队伍
	init := make(util.StringSlice, 0)
	init = append(init, "0")
	newTeam := model.Team{
		Name:          name,
		Leader:        int(user.(model.User).ID),
		InviteCode:    util.RandomString(20),
		MemberNum:     1,
		CompetitionId: init,
	}
	result := db.Create(&newTeam)
	if result.Error != nil {
		response.Response(context, http.StatusInternalServerError, 500, nil, "创建队伍失败："+result.Error.Error())
		return
	}
	//  * 获取新建队伍的id，并将队员的teamId设置为新建队伍的id
	db.Model(&model.User{}).Where("id = ?", user.(model.User).ID).Update("team_id", newTeam.ID)
	log.Printf("新建队伍的id为：%d", newTeam.ID)
	response.Success(context, nil, "队伍创建成功")
}

func JoinTeam(context *gin.Context) {
	db := common.GetDB()
	var requestTeam = model.Team{}
	err := context.ShouldBind(&requestTeam)
	if err != nil {
		return
	}
	inviteCode := requestTeam.InviteCode
	//  * 数据验证
	if len(inviteCode) == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "内容参数不能为空")
		return
	}
	user, _ := context.Get("user")
	//  * 检测当前用户是否已经注册过
	var userm model.User
	db.Where("id = ?", int(user.(model.User).ID)).First(&userm)
	if userm.ID == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}
	if userm.TeamId != 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "用户已经加入队伍")
		return
	}
	//  * 判断队伍是否存在
	var team model.Team
	db.Where("invite_code = ?", inviteCode).First(&team)
	if team.ID == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "邀请码错误，队伍不存在")
		return
	}
	//  * 将用户加入队伍
	db.Model(&model.User{}).Where("id = ?", user.(model.User).ID).Update("team_id", team.ID)
	// * 队伍人数+1
	db.Model(&model.Team{}).Where("id = ?", team.ID).Update("member_num", team.MemberNum+1)
	response.Success(context, nil, "加入队伍成功")
}

// isTeamExist
//
//	@Description: 判断队伍是否存在
//	@param db
//	@param name
//	@return bool
func isTeamExist(db *gorm.DB, name string) bool {
	var team model.Team
	db.Where("name = ?", name).First(&team)
	if team.ID != 0 {
		return true
	}
	return false
}

// QuitTeam
//
//	@Description: 退出队伍
//	@param context
func QuitTeam(context *gin.Context) {
	db := common.GetDB()
	user, _ := context.Get("user")
	//  * 检测当前用户是否已经注册过
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
	//  * 将用户退出队伍
	db.Model(&model.User{}).Where("id = ?", user.(model.User).ID).Update("team_id", 0)
	//  * 队伍人数-1
	var team model.Team
	db.Where("id = ?", userm.TeamId).First(&team)
	db.Model(&model.Team{}).Where("id = ?", userm.TeamId).Update("member_num", team.MemberNum-1)
	response.Success(context, nil, "退出队伍成功")
}

// GetTeamInfo
//
//	@Description: 获取队伍信息
//	@param context
func GetTeamInfo(context *gin.Context) {
	//  * 通过get id获取队伍信息
	id := context.Query("id")
	if len(id) == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "内容参数不能为空")
		return
	}
	var team model.Team
	common.GetDB().Where("id = ?", id).First(&team)
	response.Success(context, gin.H{"user": dto.ToTeamDto(team)}, "查询成功")
}

// DismissTeam
//
//	@Description: 解散队伍
//	@param context
func DismissTeam(context *gin.Context) {
	db := common.GetDB()
	user, _ := context.Get("user")
	//  * 检测当前用户是否已经注册过
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
	//  * 判断是否为队长
	var team model.Team
	db.Where("id = ?", userm.TeamId).First(&team)
	if team.Leader != int(user.(model.User).ID) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "只有队长可以解散队伍")
		return
	}
	//  * 解散队伍
	db.Where("id = ?", userm.TeamId).Delete(&model.Team{})
	db.Model(&model.User{}).Where("team_id = ?", userm.TeamId).Update("team_id", 0)
	response.Success(context, nil, "解散队伍成功")
}
