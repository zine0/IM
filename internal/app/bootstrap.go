package app

import (
	v1 "github/zine0/IM/internal/api/v1"
	"github/zine0/IM/internal/config"
	"github/zine0/IM/internal/initialization"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Start() {
	config.SetConfig()

	queries := initialization.InitDB()

	router := gin.New()

	logger := initialization.InitLogger()

	router.Use(ginzap.Ginzap(logger,time.RFC3339,true))
	router.Use(ginzap.RecoveryWithZap(logger,true))
	zap.ReplaceGlobals(logger)

	app := router.Group("/api/v1")
	v1.SetupRoutes(app,queries)
	router.Run(":8080")
}
