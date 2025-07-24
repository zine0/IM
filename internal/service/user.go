package service

import (
	"context"
	"github/zine0/IM/internal/repository"
	"github/zine0/IM/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/spf13/viper"
)

type UserService struct {
	db *repository.Queries
}

func NewUserService(db *repository.Queries) *UserService {
	return &UserService{
		db,
	}
}

// Request
type userRegister struct {
	ConfirmPassword string `json:"confirm_password"`
	Password        string `json:"password"`
	Username        string `json:"username"`
}

func (s *UserService) CreateUser(ctx *gin.Context) {
	req := &userRegister{}

	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "bad request"})
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "两次输入的密码不同"})
		return
	}

	_, exists := s.checkUserExists(req.Username)
	if exists {
		ctx.JSON(http.StatusOK, gin.H{"msg": "username is exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "error"})
		return
	}

	id,err := s.db.CreateUser(context.Background(), repository.CreateUserParams{
		Username:  pgtype.Text{String: req.Username,Valid: true},
		Password:  pgtype.Text{String: hashedPassword,Valid: true},
		CreatedAt: pgtype.Timestamp{Time: time.Now(),Valid: true},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "error"})
		return
	}

	token,err := utils.GenerateJWT(req.Username,int(id),viper.GetStringMapString("app")["key"])
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"msg":"error",
		})
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{"msg": "success","token":token},
	)
}

type userLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *UserService) Login(ctx *gin.Context) {
	req := &userLogin{}
	if err := ctx.ShouldBindJSON(req);err!=nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"msg":"bad request"})
		return
	}

	user,exists := s.checkUserExists(req.Username)

	if !exists {
		ctx.JSON(http.StatusBadRequest,gin.H{"msg":"no such user"})
		return
	}

	ok := utils.CheckPassword(user.Password.String,req.Password)
	if !ok {
		ctx.JSON(http.StatusOK,gin.H{"msg":"password error"})
		return
	}

	token,err := utils.GenerateJWT(user.Username.String,int(user.ID),viper.GetStringMapString("app")["key"])
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"msg":"error"})
		return
	}

	ctx.JSON(http.StatusOK,gin.H{
		"msg":"success",
		"token":token,
	})

}

func (s *UserService) checkUserExists(username string) (repository.User, bool) {
	user, err := s.db.UserExists(context.Background(), pgtype.Text{String: username,Valid: true})
	if err != nil {
		return repository.User{}, false
	}
	return user, true
}
