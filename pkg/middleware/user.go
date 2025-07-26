package middleware

import (
	"github/zine0/IM/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func AuthUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Values("Authorization")[0]
		claims,err:=utils.ValidJWT(token,viper.GetStringMapString("app")["key"])
		if err != nil {
			ctx.Abort()
			return 
		}

		ctx.Set("username",claims.Username)

		ctx.Next()
	}
}
