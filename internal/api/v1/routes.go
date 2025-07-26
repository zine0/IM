package v1

import (
	"github/zine0/IM/internal/repository"
	"github/zine0/IM/internal/service"
	"github/zine0/IM/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup, db *repository.Queries){
	SetupUserRoutes(r,db)
	SetupWSRoutes(r,db)
}

func SetupUserRoutes(r *gin.RouterGroup,db *repository.Queries) {
	s := service.NewUserService(db)
	userRoutes := r.Group("/user")

	userRoutes.POST("/create",s.CreateUser)
	userRoutes.POST("/login",s.Login)
}

func SetupWSRoutes(r *gin.RouterGroup,db *repository.Queries) {
	s := service.NewMessageService(db)
	wsRoutes := r.Group("/ws")
	wsRoutes.Use(middleware.AuthUser())
	wsRoutes.GET("",s.CreateWS)
}
