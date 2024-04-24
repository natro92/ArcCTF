package middleware

import (
	"Arc/common"
	"Arc/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(context *gin.Context) {
		// * 获取authorization
		tokenString := context.GetHeader("Authorization")
		// * 验证token格式
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Token格式错误，权限不足"})
			context.Abort()
			return
		}
		tokenString = tokenString[7:]

		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Token权限不足"})
			context.Abort()
			return
		}
		// * 验证通过后获取claims中的userId
		userId := claims.UserId
		db := common.GetDB()
		var user model.User
		db.First(&user, userId)

		// * 验证用户是否存在
		if user.ID == 0 {
			context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "ID不存在，权限不足"})
			context.Abort()
			return
		}

		// * 验证用户角色
		if user.Role == 0 {
			context.Set("role", 0)
		} else if user.Role == 1 {
			context.Set("role", 1)
		} else if user.Role == 2 {
			context.Set("role", 2)
		} else {
			context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "无效的用户角色"})
			context.Abort()
			return
		}

		// * 用户存在 将user信息写入上下文
		context.Set("user", user)
		context.Next()
	}
}
