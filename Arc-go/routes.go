package main

import (
	"Arc/controller"
	"Arc/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())

	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.GET("/api/auth/info", middleware.AuthMiddleWare(), controller.Info)

	r.GET("/api/competition/all", controller.GetAllCompetitions)
	r.POST("/api/competition/create", middleware.AuthMiddleWare(), controller.CreateCompetition)
	r.GET("/api/competition/info", middleware.AuthMiddleWare(), controller.GetCompetitionInfo4User)

	r.GET("/api/challenge/all", middleware.AuthMiddleWare(), controller.GetAllChallengesByCompetitionId)
	r.GET("/api/challenge/info", middleware.AuthMiddleWare(), controller.GetChallengeInfo4User)
	r.POST("/api/challenge/create", middleware.AuthMiddleWare(), controller.CreateChallenge)

	r.POST("/api/team/create", middleware.AuthMiddleWare(), controller.CreateNewTeam)
	r.POST("/api/team/join", middleware.AuthMiddleWare(), controller.JoinTeam)

	return r
}
