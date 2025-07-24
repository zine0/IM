package v1

import (
	"github/zine0/IM/internal/repository"
	"github/zine0/IM/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup, db *repository.Queries){
	SetupUserRoutes(r,db)
}

func SetupUserRoutes(r *gin.RouterGroup,db *repository.Queries) {
	s := service.NewUserService(db)
	userRoutes := r.Group("/user")

	userRoutes.POST("/create",s.CreateUser)
	userRoutes.POST("/login",s.Login)
}
