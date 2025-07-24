package service

import (
	"context"
	"fmt"
	"github/zine0/IM/internal/repository"
	"github/zine0/IM/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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

	exists, err := s.checkUserExists(req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "error"})
		return
	}
	if exists {
		ctx.JSON(http.StatusOK, gin.H{"msg": "username is exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "error"})
		return
	}

	fmt.Printf("%+v\n",req)
	err = s.db.CreateUser(context.Background(), repository.CreateUserParams{
		Username:  pgtype.Text{String: req.Username,Valid: true},
		Password:  pgtype.Text{String: hashedPassword,Valid: true},
		CreatedAt: pgtype.Timestamp{Time: time.Now(),Valid: true},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "error"})
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{"msg": "success"},
	)
}

func (s *UserService) checkUserExists(username string) (bool, error) {
	exists, err := s.db.UserExists(context.Background(), pgtype.Text{String: username,Valid: true})
	if err != nil {
		return false, nil
	}
	return exists, nil
}
