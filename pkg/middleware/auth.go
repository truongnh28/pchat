package middleware

import (
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func HTTPAuthentication(ctx *gin.Context) {
	var (
		token string
		err   error
	)

	token, err = ctx.Cookie(common.CookieName)
	if err != nil {
		glog.Errorln(err)
		ctx.AbortWithStatus(401)
		return
	}
	dt, err := service.GetJWTInstance().GetDataFromToken(token)
	if err != nil {
		glog.Errorln("GetDataFromToken failed: ", err)
		ctx.AbortWithStatus(401)
		return
	}
	actor := dt.(*domain.Account)
	glog.Infoln("actor", actor)
	ctx.Set("actor", actor)

	ctx.Next()
}
