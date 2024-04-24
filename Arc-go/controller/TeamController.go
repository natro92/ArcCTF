package controller

import (
	"Arc/common"
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
	newTeam := model.Team{
		Name:       name,
		Leader:     int(user.(model.User).ID),
		InviteCode: util.RandomString(20),
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
	response.Success(context, nil, "加入队伍成功")
}

func isTeamExist(db *gorm.DB, name string) bool {
	var team model.Team
	db.Where("name = ?", name).First(&team)
	if team.ID != 0 {
		return true
	}
	return false
}
