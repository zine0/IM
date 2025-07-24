package app

import (
	v1 "github/zine0/IM/internal/api/v1"
	"github/zine0/IM/internal/config"
	"github/zine0/IM/internal/initialization"

	"github.com/gin-gonic/gin"
)

func Start() {
	config.SetConfig()

	queries := initialization.InitDB()

	router := gin.Default()
	app := router.Group("/api/v1")
	v1.SetupRoutes(app,queries)
	router.Run(":8080")
}
