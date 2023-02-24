package chat

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"sync"

	"github.com/LSDXXX/libs/api"
	"github.com/LSDXXX/libs/constant"
	"github.com/LSDXXX/libs/pkg/container"
	"github.com/LSDXXX/libs/pkg/log"
	"github.com/LSDXXX/libs/pkg/util"
	"github.com/LSDXXX/libs/pkg/wsmanager"
	"github.com/LSDXXX/libs/service"
	"github.com/LSDXXX/servers/chatgpt/api/handlers/auth"
	"github.com/LSDXXX/servers/chatgpt/bot"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatHandler struct {
	userMapper   sync.Map
	manager      *wsmanager.WSManager  `container:"type"`
	conversation *service.Conversation `container:"type"`
	bot          *bot.Chatbot          `container:"type"`
	mids         []gin.HandlerFunc
}

func Register(mids ...gin.HandlerFunc) {
	handler := NewChatHandler()
	api.RegisterHttpRouter(handler)
}

func NewChatHandler(mids ...gin.HandlerFunc) *ChatHandler {
	out := ChatHandler{
		mids: mids,
	}
	util.PanicWhenError(container.Fill(&out))
	return &out
}

func (ws *ChatHandler) Use(e *gin.Engine) {
	e.GET("/api/chat", ws.manager.BuildHTTPHandler(ws))
}

func (ws *ChatHandler) OnUpgrade(c *gin.Context, manager *wsmanager.WSManager) (groupID, clientID string, err error) {
	for _, mid := range ws.mids {
		mid(c)
		if c.IsAborted() {
			return "", "", errors.New("error")
		}
	}
	v, _ := c.Get(constant.JWTIdentityKey)
	info := v.(auth.IdentityInfo)
	return strconv.Itoa(info.Id), info.WSKey, nil
}

func (ws *ChatHandler) OnClientRegister(c *wsmanager.WSClient) {
	log.WithContext(context.Background()).Debugf("register ws client: %s, ", c.Id)
}

func (ws *ChatHandler) OnClientDeregister(c *wsmanager.WSClient) {
	log.WithContext(context.Background()).Debugf("deRegister ws client: %s, ", c.Id)
}

func (ws *ChatHandler) OnClientMessage(c *wsmanager.WSClient, message []byte) error {
	var conversationId string
	v, ok := ws.userMapper.Load(c.Group)
	if ok {
		conversationId = v.(string)
	}

	parent, ok := ws.conversation.GetParentId(conversationId)
	if !ok {
		parent = uuid.NewString()
	}
	go func() {
		ch, err := ws.bot.AskStream(string(message), conversationId, parent)
		if err != nil {
			log.WithContext(context.Background()).Errorf("ask stream error: %s", err.Error())
		}
		for msg := range ch {
			data, _ := json.Marshal(msg)
			ws.manager.Send(c.Id, c.Group, data)
		}
	}()
	return nil
}
