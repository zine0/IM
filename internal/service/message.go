package service

import (
	"github/zine0/IM/internal/repository"
	"github/zine0/IM/internal/wserver"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type MessageService struct {
	db *repository.Queries
}

var (
	hub  *wserver.Hub
	once sync.Once
)

func NewMessageService(db *repository.Queries) *MessageService {
	return &MessageService{db}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (m *MessageService) CreateWS(ctx *gin.Context) {
	once.Do(func() {
		hub = wserver.NewHub()
		go hub.Run()
	})

	client, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		zap.L().Error("upgrade to ws", zap.String("error", err.Error()))
		return
	}

	username := ctx.Value("username").(string)

	hub.Register <- wserver.NewClient(username, client, hub)

}
